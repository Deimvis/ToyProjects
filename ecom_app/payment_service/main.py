import asyncpg
import logging
import os
from typing import Annotated
from contextlib import asynccontextmanager
from fastapi import Depends, FastAPI, HTTPException, status
from fastapi.responses import JSONResponse, RedirectResponse
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm

import authsdk
import db
import models
from shortcuts import HTTPBadCredentials, auth_user, auth_username


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
    title='Payment Service',
    summary='Allows to view, withdraw and replenish balance',
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


# Payment Service


@app.get('/balance')
async def get_balance(current_user: Annotated[authsdk.models.User, Depends(auth_user)]) -> models.UserBalance:
    query = 'SELECT username, balance FROM "user_balance" WHERE username = $1;'
    record = await app.db_pool.fetchrow(query, current_user.username)
    if record is None:
        await create_balance_if_not_exists(app.db_pool, current_user.username)
        record = await app.db_pool.fetchrow(query, current_user.username)
    return models.UserBalance(**record)


@app.post('/balance/{username}/replenish')
async def replenish(username: str, replenish: models.ReplenishRequest, current_user: Annotated[authsdk.models.User, Depends(auth_user)]) -> JSONResponse:
    # TODO: support receiving idempotency key
    # NOTE: username in path may be redundant
    if username != current_user.username:
        return JSONResponse(status_code=403, content={'status': 'error', 'error': f'User `{current_user.username}` is not authorised to replenish `{username}` balance'})
    balance = await create_balance_if_not_exists(app.db_pool, username)
    query = 'UPDATE "user_balance" SET balance = balance + $2 WHERE username = $1;'
    await app.db_pool.execute(query, username, replenish.amount)
    msg = f'Successfully replenished the {username} balance, current balance: {balance + replenish.amount} (+{replenish.amount})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


@app.post('/balance/{username}/withdraw')
async def withdraw(username: str, withdraw: models.WithdrawRequest, current_user: Annotated[authsdk.models.User, Depends(auth_user)]) -> JSONResponse:
    # TODO: support receiving idempotency key
    # NOTE: username in path may be redundant
    if username != current_user.username:
        return JSONResponse(status_code=403, content={'status': 'error', 'error': f'User `{current_user.username}` is not authorised to replenish `{username}` balance'})
    async with app.db_pool.acquire() as con:
        async with con.transaction():
            query = 'SELECT balance FROM "user_balance" WHERE username = $1 AND balance >= $2 FOR UPDATE;'
            balance = await con.fetchval(query, username, withdraw.amount)
            if balance is None:
                raise HTTPException(status.HTTP_403_FORBIDDEN, 'Insufficient balance')
            query = 'UPDATE "user_balance" SET balance = balance - $2 WHERE username = $1;'
            await con.execute(query, username, withdraw.amount)
    msg = f'Successfully withdrew the {username} balance, current balance: {balance - withdraw.amount} (-{withdraw.amount})'
    logging.debug(msg)
    return JSONResponse(status_code=200, content={'status': msg})


async def create_balance_if_not_exists(db_pool: asyncpg.Pool, username: str) -> int:
    query = 'INSERT INTO "user_balance" (username, balance) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET balance = "user_balance".balance RETURNING balance;'
    return int(await db_pool.fetchval(query, username, 0))
