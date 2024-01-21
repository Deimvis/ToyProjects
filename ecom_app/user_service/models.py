from datetime import datetime
from pydantic import BaseModel
from typing import Any, Dict, Self


class User(BaseModel):
    username: str
    password: str | None = None
    password_hash: str | None = None
    first_name: str
    last_name: str
    email: str
    phone: str
