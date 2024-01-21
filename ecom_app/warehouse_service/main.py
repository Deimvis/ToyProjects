import asyncpg
import httpx
import logging
import os
import time
from typing import Annotated, List
from contextlib import asynccontextmanager
from fastapi import Depends, FastAPI, HTTPException, Request, status
from fastapi.responses import JSONResponse, RedirectResponse
from fastapi.security import OAuth2PasswordBearer

import db
import models
from shortcuts import auth_user


if os.getenv('DEBUG'):
    logging.basicConfig(
        level=logging.DEBUG,
        format='%(levelname)s [%(name)s.%(funcName)s:%(lineno)d] %(message)s',
        datefmt="%d/%b/%Y %H:%M:%S",
    )
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('yt').setLevel(logging.WARNING)
else:
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s] %(levelname)s [%(name)s] %(message)s',
        datefmt="%d/%b/%Y %H:%M:%S",
    )


@asynccontextmanager
async def lifespan(app: FastAPI):
    app.db_pool = await db.create_pool()
    await db.init_db(app.db_pool)
    yield
    await app.db_pool.close()


app = FastAPI(
    title='Warehouse Service',
    summary='Allows to update warehouse related info for items and manage item reservations',
    version='0.0.1',
    lifespan=lifespan,
)
ouath2_scheme = OAuth2PasswordBearer(tokenUrl='token')

@app.get('/', response_class=RedirectResponse, include_in_schema=False)
async def index():
    return RedirectResponse(url='/docs')


@app.get('/health')
async def health():
    return {'status': 'ok'}


# Warehouse Service

@app.get('/items/{item_id}')
async def get_item(item_id: str) -> models.Item:
    item = await get_item(app.db_pool, item_id)
    if item is None:
        raise HTTPException(status_code=404, detail='Item not found')
    return item


@app.post('/items/{item_id}')
async def update_item(item_id: str, item: models.Item):
    if item_id != item.item_id:
        raise HTTPException(status_code=403, detail='Item id in the path is different from item id in the body')
    query = 'INSERT INTO "item_wh_info" (item_id, in_stock) VALUES ($1, $2) ON CONFLICT (item_id) DO UPDATE SET in_stock = $2'
    await app.db_pool.execute(query, item.item_id, item.in_stock)
    return JSONResponse(status_code=200, content={'status': f'Item {item.item_id} has been successfully updated'})


@app.post('/items/{item_id}/reservation')
async def reserve_item(item_id: str, transaction_id: str):
    tx_info = await get_tx_info(app.db_pool, transaction_id)
    if tx_info is not None:
        match tx_info.status:
            case models.TransactionInfo.Status.RUNNING:
                reserved_item = await get_reserved_item(app.db_pool, transaction_id)
                assert reserved_item is not None, 'Inconsistency happened - some item should be reserved'
                if reserved_item.item_id != item_id:
                    raise HTTPException(status_code=403, detail='Another item is reserved for current transaction')
                return JSONResponse(status_code=200, content={'status': 'Item is already reserved for current transaction'})
            case models.TransactionInfo.Status.COMMITED:
                raise HTTPException(status_code=403, detail='Transaction is already committed')
            case models.TransactionInfo.Status.ROLLBACKED:
                raise HTTPException(status_code=403, detail='Transaction is already rollbacked')
    item = await get_item(app.db_pool, item_id)
    if item is None:
        raise HTTPException(status_code=404, detail='Item not found')

    async with app.db_pool.acquire() as con:
        async with con.transaction():
            query = 'SELECT in_stock FROM "item_wh_info" WHERE item_id = $1 AND in_stock >= 1 FOR UPDATE;'
            in_stock = await con.fetchval(query, item_id)
            if in_stock is None:
                raise HTTPException(status.HTTP_403_FORBIDDEN, 'Not enough items in stock')

            query = 'UPDATE "item_wh_info" SET in_stock = in_stock - 1 WHERE item_id = $1;'
            await con.execute(query, item_id)
            creation_ts = int(time.time())
            expiration_ts = creation_ts + 24 * 60 * 60
            query = 'INSERT INTO "reserved_item" (transaction_id, item_id, creation_ts, expiration_ts) VALUES ($1, $2, $3, $4);'
            await con.execute(query, transaction_id, item_id, creation_ts, expiration_ts)
            query = 'INSERT INTO "transaction_wh_info" (transaction_id, status) VALUES ($1, $2);'
            await con.execute(query, transaction_id, models.TransactionInfo.Status.RUNNING)

    msg = f'Successfully reserved the item(id={item_id}) for transaction(id={transaction_id})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


@app.post('/items/{item_id}/reservation/commit')
async def commit_reservation(item_id: str, transaction_id: str):
    tx_info = await get_tx_info(app.db_pool, transaction_id)
    if tx_info is None:
        raise HTTPException(status_code=403, detail='Transaction is not found')
    if tx_info.status != models.TransactionInfo.Status.RUNNING:
        raise HTTPException(status_code=403, detail=f'Transaction has a wrong status: "{tx_info.status}" ("{models.TransactionInfo.Status.RUNNING.value}" is expected)')
    reserved_item = await get_reserved_item(app.db_pool, transaction_id)
    if reserved_item is None:
        return JSONResponse(status_code=404, content={'status': 'Reservations is not found for current transaction'})
    if reserved_item.item_id != item_id:
        raise HTTPException(status_code=403, detail='Another item is reserved for current transaction')

    async with app.db_pool.acquire() as con:
        con: asyncpg.Connection  # noqa
        async with con.transaction():
            query = 'DELETE FROM "reserved_item" WHERE transaction_id = $1'
            await con.execute(query, transaction_id)
            query = 'UPDATE "transaction_wh_info" SET status = $2 WHERE transaction_id = $1'
            await con.execute(query, transaction_id, models.TransactionInfo.Status.COMMITED)

    msg = f'Successfully completed the reservation({transaction_id=})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


@app.post('/items/{item_id}/reservation/rollback')
async def rollback_reservation(item_id: str, transaction_id: str):
    tx_info = await get_tx_info(app.db_pool, transaction_id)
    if tx_info is None:
        raise HTTPException(status_code=403, detail='Transaction is not found')
    if tx_info.status != models.TransactionInfo.Status.RUNNING:
        raise HTTPException(status_code=403, detail=f'Transaction has a wrong status: "{tx_info.status}" ("{models.TransactionInfo.Status.RUNNING.value}" is expected)')
    reserved_item = await get_reserved_item(app.db_pool, transaction_id)
    if reserved_item is None:
        return JSONResponse(status_code=404, content={'status': 'Reservations is not found for current transaction'})
    if reserved_item.item_id != item_id:
        raise HTTPException(status_code=403, detail='Another item is reserved for current transaction')

    async with app.db_pool.acquire() as con:
        async with con.transaction():
            query = 'DELETE FROM "reserved_item" WHERE transaction_id = $1'
            await con.execute(query, transaction_id)
            query = 'UPDATE "item_wh_info" SET in_stock = in_stock + 1 WHERE item_id = $1'
            await con.execute(query, item_id)
            query = 'UPDATE "transaction_wh_info" SET status = $2 WHERE transaction_id = $1'
            await con.execute(query, transaction_id, models.TransactionInfo.Status.ROLLBACKED)

    msg = f'Successfully rollbacked the reservation({transaction_id=})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


async def get_item(db_pool: asyncpg.Pool, item_id: str) -> models.Item | None:
    query = 'SELECT item_id, in_stock FROM "item_wh_info" WHERE item_id = $1;'
    row = await db_pool.fetchrow(query, item_id)
    if row is None:
        return None
    return models.Item(item_id=row['item_id'], in_stock=row['in_stock'])


async def get_reserved_item(db_pool: asyncpg.Pool, transaction_id: str) -> models.ReservedItem | None:
    query = 'SELECT transaction_id, item_id, creation_ts, expiration_ts FROM "reserved_item" WHERE transaction_id = $1;'
    row = await db_pool.fetchrow(query, transaction_id)
    if row is None:
        return None
    return models.ReservedItem(**row)


async def get_tx_info(db_pool: asyncpg.Pool, transaction_id: str) -> models.TransactionInfo | None:
    query = 'SELECT transaction_id, status FROM "transaction_wh_info" WHERE transaction_id = $1;'
    row = await db_pool.fetchrow(query, transaction_id)
    if row is None:
        return None
    return models.TransactionInfo(**row)


async def item_exists(db_pool: asyncpg.Pool, item_id: str) -> bool:
    query = 'SELECT EXISTS (SELECT 1 FROM "item" WHERE id = $1);'
    return await db_pool.fetchval(query, item_id)
