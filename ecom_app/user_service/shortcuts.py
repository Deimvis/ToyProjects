from fastapi import Depends, HTTPException, status
from typing import Annotated

import authsdk
import models


HTTPBadCredentials = HTTPException(
    status_code=status.HTTP_401_UNAUTHORIZED,
    detail='Incorrect username or password',
    headers={'WWW-Authenticate': 'Bearer'},
)


async def auth_user(token: Annotated[str, Depends(authsdk.conf.OAUTH2_SCHEME)]) -> models.User:
    user = await authsdk.get_current_user(token, bad_creds_exc=HTTPBadCredentials)
    return models.User.model_validate(user.model_dump())


async def auth_username(token: Annotated[str, Depends(authsdk.conf.OAUTH2_SCHEME)]) -> models.User:
    return await authsdk.get_current_username(token, bad_creds_exc=HTTPBadCredentials)
