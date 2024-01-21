from pydantic import BaseModel, field_validator


class UserBalance(BaseModel):
    username: str
    balance: int


class ReplenishRequest(BaseModel):
    amount: int

    @field_validator('amount')
    @classmethod
    def amount_is_valid(cls, v):
        assert v >= 0, 'Can\'t replenish negative amount of money'
        return v


class WithdrawRequest(BaseModel):
    amount: int

    @field_validator('amount')
    @classmethod
    def amount_is_valid(cls, v):
        assert v >= 0, 'Can\'t withdraw negative amount of money'
        return v

