from datetime import datetime, timedelta
from jose import jwt

import authsdk.models as models
from authsdk.conf import *  # noqa


def verify_password(plain_password, hashed_password):
    return PWD_CONTEXT.verify(plain_password, hashed_password)


def get_password_hash(password):
    return PWD_CONTEXT.hash(password)


def create_access_token(subject: str, expires_delta: timedelta = timedelta(minutes=DEFAULT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES)):
    payload = models.TokenPayload(
        sub=subject,
        exp=int((datetime.utcnow() + expires_delta).timestamp()),
    )
    return jwt.encode(payload.model_dump(), JWT_HS256_KEY, algorithm=ALGORITHM)


def create_refresh_token(subject: str, expires_delta: timedelta = timedelta(minutes=DEFAULT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES)):
    payload = models.TokenPayload(
        sub=subject,
        exp=int((datetime.utcnow() + expires_delta).timestamp()),
    )
    return  jwt.encode(payload.model_dump(), JWT_HS256_KEY, algorithm=ALGORITHM)
