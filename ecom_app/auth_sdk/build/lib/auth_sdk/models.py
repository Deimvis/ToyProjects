from datetime import datetime
from pydantic import BaseModel


class User(BaseModel):
    username: str
    first_name: str
    last_name: str
    email: str
    phone: str



class TokenResponse(BaseModel):
    access_token: str
    refresh_token: str


class TokenPayload(BaseModel):
    sub: str
    exp: int

    def is_expired(self) -> bool:
        return self.exp < int(datetime.utcnow().timestamp())
