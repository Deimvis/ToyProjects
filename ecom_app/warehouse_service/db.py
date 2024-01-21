import asyncpg
import os


DB_POOL = None


INIT_DB_SCRIPT = '''
CREATE TABLE IF NOT EXISTS "item_wh_info" (
    item_id VARCHAR(512) PRIMARY KEY,
    in_stock BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "reserved_item" (
    transaction_id VARCHAR(512) PRIMARY KEY,
    item_id VARCHAR(512) NOT NULL,
    creation_ts BIGINT,
    expiration_ts BIGINT
);

CREATE TABLE IF NOT EXISTS "transaction_wh_info" (
    "transaction_id" VARCHAR(512) PRIMARY KEY,
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
