from enum import Enum
from pydantic import BaseModel
from typing import Dict, Self


class Item(BaseModel):
    item_id: str
    in_stock: int


class ReservedItem(BaseModel):
    transaction_id: str
    item_id: str
    creation_ts: int
    expiration_ts: int


class TransactionInfo(BaseModel):
    class Status(str, Enum):
        RUNNING = 'running'
        COMMITED = 'commited'
        ROLLBACKED = 'rollbacked'

    transaction_id: str
    status: str