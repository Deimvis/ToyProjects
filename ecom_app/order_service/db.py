import asyncpg
import os


DB_POOL = None


INIT_DB_SCRIPT = '''
CREATE TABLE IF NOT EXISTS "item" (
    id VARCHAR(512) PRIMARY KEY,
    cost BIGINT
);
INSERT INTO "item" (id, cost) VALUES ('cheap_item_id', 100) ON CONFLICT (id) DO NOTHING;
INSERT INTO "item" (id, cost) VALUES ('expensive_item_id', 1000) ON CONFLICT (id) DO NOTHING;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "order" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(512),
    item_id VARCHAR(512) REFERENCES item(id),
    courier_id VARCHAR(512) REFERENCES courier(id),
    creation_ts BIGINT DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT,
    successful BOOLEAN
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
