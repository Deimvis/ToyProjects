import asyncpg
import os
import models


DB_POOL = None


INIT_DB_SCRIPT = '''
CREATE TABLE IF NOT EXISTS "user" (
    username VARCHAR(512) PRIMARY KEY,
    password_hash VARCHAR(128),
    first_name VARCHAR(256),
    last_name VARCHAR(256),
    email VARCHAR(512),
    phone VARCHAR(128)
);
'''


async def init_db(pool):
    async with pool.acquire() as con:
        await con.execute(INIT_DB_SCRIPT)


async def create_pool():
    global DB_POOL
    assert DB_POOL is None, 'db pool has already been created'
    DB_POOL = await asyncpg.create_pool(
        user=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        database=os.getenv('POSTGRES_DB'),
        host=os.getenv('POSTGRES_HOST'),
        port=int(os.getenv('POSTGRES_PORT', '5432')),
    )
    return DB_POOL


async def get_user(db_pool: asyncpg.Pool, username: str) -> models.User | None:
    """ Be careful, it fetches all user info including password_hash """
    query = 'SELECT username, password_hash, first_name, last_name, email, phone FROM "user" WHERE username = $1;'
    record = await db_pool.fetchrow(query, username)
    if record is None:
        return None
    return models.User(**record)