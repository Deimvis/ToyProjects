import asyncpg
import os


DB_POOL = None


INIT_DB_SCRIPT = '''
CREATE TABLE IF NOT EXISTS "courier" (
    id VARCHAR(512) PRIMARY KEY,
    status VARCHAR(32)
);

CREATE TABLE IF NOT EXISTS "transaction_delivery_info" (
    "transaction_id" VARCHAR(512) PRIMARY KEY,
    assigned_courier_id VARCHAR(512) references courier(id),
    status VARCHAR(32)
)
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
