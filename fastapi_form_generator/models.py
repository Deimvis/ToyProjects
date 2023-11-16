from utils import WithFormGen
from enum import Enum
from pydantic import BaseModel, Field
from typing import List

@WithFormGen(method_name='form', by_alias=True)
class Movie(BaseModel):
    class Person(BaseModel):
        name: str
        age: int

    class Genre(str, Enum):
        COMEDY = 'comedy'

    name: str = Field('title')
    genres: List[Genre] = Field(default_factory=list)
    director: Person
