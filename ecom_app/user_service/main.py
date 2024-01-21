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
    title='User Service',
    summary='Provides CRUD API to user data and is responsible for auth',
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


# User CRUD API


@app.post('/user')
async def create_user(user: models.User) -> JSONResponse:
    if await is_user_exists(user.username):
        return JSONResponse(status_code=403, content={'status': 'error', 'error': 'User already exists'})
    password_hash = authsdk.internal.get_password_hash(user.password)
    query = 'INSERT INTO "user" (username, password_hash, first_name, last_name, email, phone) VALUES ($1, $2, $3, $4, $5, $6);'
    await app.db_pool.execute(query, user.username, password_hash, user.first_name, user.last_name, user.email, user.phone)
    return JSONResponse(status_code=200, content={'status': f'User `{user.username}` was created'})


@app.get('/user/{username}')
async def get_user(username: str) -> models.User:
    query = 'SELECT username, first_name, last_name, email, phone FROM "user" WHERE username = $1;'
    record = await app.db_pool.fetchrow(query, username)
    if record is None:
        return JSONResponse(status_code=404, content={'status': 'error', 'error': 'User not found'})
    return models.User(**record)


@app.put('/user/{username}')
async def update_user(user: models.User, current_user: Annotated[models.User, Depends(auth_user)]) -> JSONResponse:
    if user.username != current_user.username:
        return JSONResponse(status_code=403, content={'status': 'error', 'error': f'Unable to update `{user.username}` from `{current_user.username}'})
    if not await is_user_exists(user.username):
        return JSONResponse(status_code=404, content={'status': 'error', 'error': 'User not found'})
    password_hash = authsdk.internal.get_password_hash(user.password)
    query = 'UPDATE "user" SET username = $1, password_hash = $2, first_name = $3, last_name = $4, email = $5, phone = $6 WHERE username = $1;'
    await app.db_pool.execute(query, user.username, password_hash, user.first_name, user.last_name, user.email, user.phone)
    return JSONResponse(status_code=200, content={'status': f'User `{user.username}` was updated'})


@app.delete('/user/{username}')
async def delete_user(username: str, current_username: Annotated[str, Depends(auth_username)]) -> JSONResponse:
    if username != current_username:
        return JSONResponse(status_code=403, content={'status': 'error', 'error': f'Unable to delete `{username}` from `{current_username}'})
    if not await is_user_exists(username):
        return JSONResponse(status_code=404, content={'status': 'error', 'error': 'User not found'})
    query = 'DELETE FROM "user" WHERE username = $1;'
    await app.db_pool.execute(query, username)
    return JSONResponse(status_code=200, content={'status': f'User `{username}` was deleted'})


async def is_user_exists(username) -> bool:
    query = 'SELECT username FROM "user" WHERE username = $1'
    record = await app.db_pool.fetchrow(query, username)
    return record is not None


# Auth

@app.post('/signup', status_code=status.HTTP_307_TEMPORARY_REDIRECT, responses={307: {'description': 'Redirected to /user'}})
async def signup(user: models.User):
    return RedirectResponse('/user')


@app.post('/login', status_code=status.HTTP_307_TEMPORARY_REDIRECT, responses={307: {'description': 'Redirected to /token'}})
async def login(form_data: Annotated[OAuth2PasswordRequestForm, Depends()]):
    return RedirectResponse('/token')


@app.post("/token", response_model=authsdk.models.TokenResponse)
async def token(form_data: Annotated[OAuth2PasswordRequestForm, Depends()]) -> authsdk.models.TokenResponse:
    user = await db.get_user(app.db_pool, form_data.username)
    if user is None:
        raise HTTPBadCredentials
    if not authsdk.internal.verify_password(form_data.password, user.password_hash):
        raise HTTPBadCredentials
    return authsdk.models.TokenResponse(
        access_token=authsdk.internal.create_access_token(subject=user.username),
        refresh_token=authsdk.internal.create_refresh_token(subject=user.username),
    )
