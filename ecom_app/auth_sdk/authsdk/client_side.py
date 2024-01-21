import httpx
import logging
from pydantic import ValidationError
from jose import JWTError, jwt

import authsdk.models as models
from authsdk.conf import *  # noqa
from authsdk.exceptions import BadCredentials


async def get_current_username(token: str, bad_creds_exc=BadCredentials()) -> str:
    try:
        payload_data = jwt.decode(token, JWT_HS256_KEY, algorithms=[ALGORITHM], options={'verify_exp': False})
        payload = models.TokenPayload(**payload_data)
        if payload.sub is None:
            logging.debug('"sub" is not found in token payload')
            raise bad_creds_exc
        if payload.is_expired():
            logging.debug('Token is expired')
            raise bad_creds_exc
    except (JWTError, ValidationError) as error:
        logging.debug(f'Failed to decode token: {error}')
        raise bad_creds_exc
    return payload.sub


async def get_current_user(token: str, bad_creds_exc=BadCredentials()) -> models.User:
    username = await get_current_username(token, bad_creds_exc=bad_creds_exc)
    url = 'http://{user_service_addr}/user/{username}'.format(user_service_addr=USER_SERVICE_ADDR, username=username)
    async with httpx.AsyncClient() as client:
        resp = await client.get(url)
    logging.debug(f'Get user with username="{username}" -> status_code={resp.status_code}')
    match resp.status_code:
        case 200:
            user = models.User.model_validate(resp.json())
        case 404:
            raise bad_creds_exc
        case _:
            raise bad_creds_exc
    return user
