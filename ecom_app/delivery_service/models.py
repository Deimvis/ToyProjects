from enum import Enum
from pydantic import BaseModel
from typing import Dict, Self


class Courier(BaseModel):
    class Status(str, Enum):
        AVAILABLE = 'available'
        RESERVED = 'reserved'
        BUSY = 'busy'

    id: str
    status: Status


class TransactionInfo(BaseModel):
    class Status(str, Enum):
        RUNNING = 'running'
        COMMITED = 'commited'
        ROLLBACKED = 'rollbacked'

    transaction_id: str
    status: Status
    assigned_courier_id: str