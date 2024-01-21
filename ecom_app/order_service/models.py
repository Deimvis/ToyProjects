from pydantic import BaseModel
from typing import Dict, Self


class Item(BaseModel):
    id: str
    cost: int


class Order(BaseModel):
    id: str
    username: str
    item_id: str
    creation_ts: int
    successful: bool

    @staticmethod
    def from_record(record: Dict) -> Self:
        print(record)
        return Order(
            id=record['id'].hex,
            username=record['username'],
            item_id=record['item_id'],
            creation_ts=record['creation_ts'],
            successful=record['successful'],
        )


class CreateOrderRequest(BaseModel):
    item_id: str
