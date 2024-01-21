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
    title='Delivery Service',
    summary='Allows to update and manage couriers info',
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


# Delivery Service

@app.get('/couriers/{courier_id}')
async def get_courier(courier_id: str) -> models.Courier:
    courier = await get_courier(app.db_pool, courier_id)
    if courier is None:
        raise HTTPException(status_code=404, detail='Courier not found')
    return courier


@app.post('/couriers/{courier_id}')
async def update_courier(courier_id: str, courier: models.Courier):
    if courier_id != courier.id:
        raise HTTPException(status_code=403, detail='Item id in the path is different from item id in the body')
    query = 'INSERT INTO "courier" (id, status) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET status = $2'
    await app.db_pool.execute(query, courier.id, courier.status)
    return JSONResponse(status_code=200, content={'status': f'Courier {courier.id} has been successfully updated'})


# NOTE: For some reason path "/couriers/reservation" throws a Pydantic exception, solution is not found, another path is used instead
@app.post('/couriers/reservation/create')
async def reserve_courier(transaction_id: str):
    tx_info = await get_tx_info(app.db_pool, transaction_id)
    if tx_info is not None:
        match tx_info.status:
            case models.TransactionInfo.Status.RUNNING:
                return JSONResponse(status_code=200, content={'status': f'Courier(id={tx_info.assigned_courier_id}) is already reserved for current transaction'})
            case models.TransactionInfo.Status.COMMITED:
                raise HTTPException(status_code=403, detail='Transaction is already committed')
            case models.TransactionInfo.Status.ROLLBACKED:
                raise HTTPException(status_code=403, detail='Transaction is already rollbacked')

    async with app.db_pool.acquire() as con:
        con: asyncpg.Connection  # noqa
        async with con.transaction():
            query = 'SELECT id FROM "courier" WHERE status = $1 LIMIT 1 FOR UPDATE;'
            courier_id = await con.fetchval(query, models.Courier.Status.AVAILABLE)
            if courier_id is None:
                raise HTTPException(status.HTTP_403_FORBIDDEN, 'No available courier found')

            query = 'UPDATE "courier" SET status = $2 WHERE id = $1'
            await con.execute(query, courier_id, models.Courier.Status.RESERVED)
            query = 'INSERT INTO "transaction_delivery_info" (transaction_id, assigned_courier_id, status) VALUES ($1, $2, $3);'
            await con.execute(query, transaction_id, courier_id, models.TransactionInfo.Status.RUNNING)

    msg = f'Successfully reserved the courier(id={courier_id}) for transaction(id={transaction_id})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg, 'courier_id': courier_id})


@app.post('/couriers/reservation/commit')
async def commit_reservation(transaction_id: str):
    async with app.db_pool.acquire() as con:
        con: asyncpg.Connection  # noqa
        async with con.transaction():
            query = 'SELECT transaction_id, status, assigned_courier_id FROM "transaction_delivery_info" WHERE transaction_id = $1 FOR UPDATE;'
            row = await con.fetchrow(query, transaction_id)
            if row is None:
                raise HTTPException(status_code=404, detail='Transaction is not found')
            tx_info = models.TransactionInfo(**row)
            if tx_info.status != models.TransactionInfo.Status.RUNNING:
                raise HTTPException(status_code=403, detail=f'Transaction has a wrong status: "{tx_info.status.value}" ("{models.TransactionInfo.Status.RUNNING.value}" is expected)')

            query = 'UPDATE "courier" SET status = $2 WHERE id = $1'
            await con.execute(query, tx_info.assigned_courier_id, models.Courier.Status.BUSY)
            query = 'UPDATE "transaction_delivery_info" SET status = $2 WHERE transaction_id = $1'
            await con.execute(query, transaction_id, models.TransactionInfo.Status.COMMITED)

    msg = f'Successfully completed the reservation({transaction_id=})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


@app.post('/couriers/reservation/rollback')
async def rollback_reservation(transaction_id: str):
    async with app.db_pool.acquire() as con:
        con: asyncpg.Connection  # noqa
        async with con.transaction():
            query = 'SELECT transaction_id, status, assigned_courier_id FROM "transaction_delivery_info" WHERE transaction_id = $1 FOR UPDATE;'
            row = await con.fetchrow(query, transaction_id)
            if row is None:
                raise HTTPException(status_code=404, detail='Transaction is not found')
            tx_info = models.TransactionInfo(**row)
            if tx_info.status != models.TransactionInfo.Status.RUNNING:
                raise HTTPException(status_code=403, detail=f'Transaction has a wrong status: "{tx_info.status.value}" ("{models.TransactionInfo.Status.RUNNING.value}" is expected)')

            query = 'UPDATE "courier" SET status = $2 WHERE id = $1'
            await con.execute(query, tx_info.assigned_courier_id, models.Courier.Status.AVAILABLE)
            query = 'UPDATE "transaction_delivery_info" SET status = $2 WHERE transaction_id = $1'
            await con.execute(query, transaction_id, models.TransactionInfo.Status.ROLLBACKED)

    msg = f'Successfully rollbacked the reservation({transaction_id=})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


async def get_courier(db_pool: asyncpg.Pool, courier_id: str) -> models.Courier | None:
    query = 'SELECT id, status FROM "courier" WHERE id = $1'
    row = await app.db_pool.fetchrow(query, courier_id)
    if row is None:
        return None
    return models.Courier(**row)


async def get_tx_info(db_pool: asyncpg.Pool, transaction_id: str) -> models.TransactionInfo | None:
    query = 'SELECT transaction_id, status, assigned_courier_id FROM "transaction_delivery_info" WHERE transaction_id = $1;'
    row = await db_pool.fetchrow(query, transaction_id)
    if row is None:
        return None
    return models.TransactionInfo(**row)
