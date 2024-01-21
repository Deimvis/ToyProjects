import os
from fastapi.security import OAuth2PasswordBearer
from passlib.context import CryptContext


USER_SERVICE_ADDR=os.getenv('USER_SERVICE_ADDR')
assert USER_SERVICE_ADDR is not None, 'USER_SERVICE_ADDR env variable is not set'
JWT_HS256_KEY=os.getenv('JWT_HS256_KEY')
assert JWT_HS256_KEY is not None, 'JWT_HS256_KEY env varible is not set'
ALGORITHM = 'HS256'
DEFAULT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES = 30
DEFAULT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES = 7 * 24 * 60  # 7 days

OAUTH2_SCHEME = OAuth2PasswordBearer(tokenUrl='token', scheme_name='JWT')
PWD_CONTEXT = CryptContext(schemes=['bcrypt'], deprecated='auto')
