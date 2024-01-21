import asyncpg
import httpx
import logging
import os
import uuid
from typing import Annotated, List
from contextlib import asynccontextmanager
from fastapi import Depends, FastAPI, HTTPException, Request, status
from fastapi.responses import JSONResponse, RedirectResponse
from fastapi.security import OAuth2PasswordBearer

import authsdk
import db
import models
from conf import *  # noqa
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
    title='Order Service',
    summary='Allows to create a new order and to view an order history',
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


# Order Service


@app.get('/orders')
async def get_orders(current_user: Annotated[authsdk.models.User, Depends(auth_user)]) -> List[models.Order]:
    query = 'SELECT id, username, item_id, creation_ts, successful FROM "order" WHERE username = $1;'
    records = await app.db_pool.fetch(query, current_user.username)
    return [models.Order.from_record(record) for record in records]


@app.post('/orders')
async def create_order(request: Request, order: models.CreateOrderRequest, current_user: Annotated[authsdk.models.User, Depends(auth_user)]) -> JSONResponse:
    # TODO: support receiving idempotency key
    try:
        item = await get_item(app.db_pool, order.item_id)
        if item is None:
            raise HTTPException(status_code=404, detail='Item not found')
        await withdraw(current_user.username, item.cost, {'Authorization': request.headers['Authorization']})
        try:
            tx_id = uuid.uuid4()
            async with httpx.AsyncClient() as client:
                headers = {'Authorization': request.headers['Authorization']}
                wh_req = client.post('http://{warehouse_service_addr}/items/{item_id}/reservation?transaction_id={tx_id}'.format(warehouse_service_addr=WAREHOUSE_SERVICE_ADDR, item_id=order.item_id, tx_id=tx_id), headers=headers)
                delivery_req = client.post('http://{delivery_service_addr}/couriers/reservation/create?transaction_id={tx_id}'.format(delivery_service_addr=DELIVERY_SERVICE_ADDR, tx_id=tx_id), headers=headers)
                wh_resp = await wh_req
                delivery_resp = await delivery_req
                if wh_resp.status_code != 200 and delivery_resp.status_code != 200:
                    raise HTTPException(status_code=403, detail=f'Warehouse and Delivery services failed on prepare stage:\nwarehouse: {wh_resp.content}\ndelivery: {delivery_resp.content}')
                if wh_resp.status_code != 200:
                    await client.post('http://{delivery_service_addr}/couriers/reservation/rollback?transaction_id={tx_id}'.format(delivery_service_addr=DELIVERY_SERVICE_ADDR, tx_id=tx_id), headers=headers)
                    raise HTTPException(status_code=403, detail=f'Warehouse service failed on prepare stage: {wh_resp.content}')
                if delivery_resp.status_code != 200:
                    await client.post('http://{warehouse_service_addr}/items/{item_id}/reservation/rollback?transaction_id={tx_id}'.format(warehouse_service_addr=WAREHOUSE_SERVICE_ADDR, item_id=order.item_id, tx_id=tx_id), headers=headers)
                    raise HTTPException(status_code=403, detail=f'Delivery service failed on prepare stage: {delivery_resp.content}')
                courier_id = delivery_resp.json()['courier_id']
                wh_req = client.post('http://{warehouse_service_addr}/items/{item_id}/reservation/commit?transaction_id={tx_id}'.format(warehouse_service_addr=WAREHOUSE_SERVICE_ADDR, item_id=order.item_id, tx_id=tx_id), headers=headers)
                delivery_req = client.post('http://{delivery_service_addr}/couriers/reservation/commit?transaction_id={tx_id}'.format(delivery_service_addr=DELIVERY_SERVICE_ADDR, tx_id=tx_id), headers=headers)
                assert (await wh_req).status_code == 200
                assert (await delivery_req).status_code == 200
        except:
            await replenish(current_user.username, item.cost, {'Authorization': request.headers['Authorization']})
            raise
    except:
        query = 'INSERT INTO "order" (username, item_id, successful) VALUES ($1, $2, $3) RETURNING id, creation_ts'
        row = await app.db_pool.fetchrow(query, current_user.username, order.item_id, False)
        logging.debug(f'Failed to create a new order, id = {row["id"]}, creation_ts = {row["creation_ts"]}')
        raise

    query = 'INSERT INTO "order" (username, item_id, courier_id, successful) VALUES ($1, $2, $3, $4) RETURNING id, creation_ts'
    row = await app.db_pool.fetchrow(query, current_user.username, order.item_id, courier_id, True)
    msg = f'Successfully created a new order, id = {row["id"]}, courier_id = {courier_id}, creation_ts = {row["creation_ts"]}'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg, 'courier_id': courier_id})


async def withdraw(username: str, amount: int, headers):
    url = 'http://{payment_service_addr}/balance/{username}/withdraw'.format(payment_service_addr=PAYMENT_SERVICE_ADDR, username=username)
    async with httpx.AsyncClient() as client:
        resp = await client.post(url, headers=headers, json={'amount': amount})
        if resp.status_code != 200:
            raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail=f'Failed to withdraw money from the balance (payment service returned status_code={resp.status_code})')


async def replenish(username: str, amount: int, headers):
    url = 'http://{payment_service_addr}/balance/{username}/replenish'.format(payment_service_addr=PAYMENT_SERVICE_ADDR, username=username)
    async with httpx.AsyncClient() as client:
        resp = await client.post(url, headers=headers, json={'amount': amount})
        if resp.status_code != 200:
            raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail=f'Failed to replenish money from the balance (payment service returned status_code={resp.status_code})')


async def get_item(db_pool: asyncpg.Pool, item_id: str) -> models.Item:
    query = 'SELECT id, cost FROM "item" WHERE id = $1;'
    row = await db_pool.fetchrow(query, item_id)
    if row is None:
        return None
    return models.Item(id=row['id'], cost=row['cost'])


async def item_exists(db_pool: asyncpg.Pool, item_id: str) -> bool:
    query = 'SELECT EXISTS (SELECT 1 FROM "item" WHERE id = $1);'
    res =  await db_pool.fetchval(query, item_id)
    return res
