import asyncpg
import os


DB_POOL = None


INIT_DB_SCRIPT = '''
CREATE TABLE IF NOT EXISTS "user_balance" (
    username VARCHAR(512) PRIMARY KEY,
    balance BIGINT CHECK (balance >= 0)
);
'''


async def init_db(pool: asyncpg.Pool):
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
