import logging
import os
import psycopg2
import pytest
import requests
import time
from typing import List


USER_SERVICE_ADDR = os.getenv('USER_SERVICE_ADDR')
assert USER_SERVICE_ADDR is not None
PAYMENT_SERVICE_ADDR = os.getenv('PAYMENT_SERVICE_ADDR')
assert PAYMENT_SERVICE_ADDR is not None
ORDER_SERVICE_ADDR = os.getenv('ORDER_SERVICE_ADDR')
assert ORDER_SERVICE_ADDR is not None
WAREHOUSE_SERVICE_ADDR = os.getenv('WAREHOUSE_SERVICE_ADDR')
assert WAREHOUSE_SERVICE_ADDR is not None
assert WAREHOUSE_SERVICE_ADDR != ''
DELIVERY_SERVICE_ADDR = os.getenv('DELIVERY_SERVICE_ADDR')
assert DELIVERY_SERVICE_ADDR is not None
TEST_USER_PREFIX = 'TEST_de3699d0b11ac8dd_'


@pytest.fixture
def setup():
    code = None
    while code != 200:
        try:
            code = requests.get(f'http://{USER_SERVICE_ADDR}/health').status_code
        except:
            pass
        time.sleep(0.3)
    con = psycopg2.connect(
        user=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        dbname=os.getenv('POSTGRES_DB'),
        host=os.getenv('POSTGRES_HOST'),
        port=int(os.getenv('POSTGRES_PORT', '5432')),
    )
    con.autocommit = True
    remove_test_usernames(con)
    yield con
    remove_test_usernames(con)


def test_simple_scenario(setup):
    username = f'{TEST_USER_PREFIX}_somename'
    password = 'password'


    resp = requests.get(f'http://{USER_SERVICE_ADDR}/user/{username}')
    assert resp.status_code == 404
    logging.info('User doesn\'t exist')


    resp = requests.post(f'http://{USER_SERVICE_ADDR}/user', json={
        "username": username,
        "password": password,
        "first_name": "Why",
        "last_name": "Do",
        "email": "You",
        "phone": "Care"
    })
    assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
    logging.info('User is created')


    headers = {'Content-Type': 'application/x-www-form-urlencoded'}
    data='&'.join(map(lambda i: f'{i[0]}={i[1]}', {
        "grant_type": "password",
        "username": username,
        "password": password,
    }.items()))
    resp = requests.post(f'http://{USER_SERVICE_ADDR}/token', headers=headers, data=data)
    assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
    access_token = resp.json()['access_token']
    refresh_token = resp.json()['refresh_token']
    logging.info('Signed in - Received acess token')


    auth_headers = {'Authorization': f'Bearer {access_token}'}
    def check_balance(expected: int):
        resp = requests.get(f'http://{PAYMENT_SERVICE_ADDR}/balance', headers=auth_headers)
        assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
        assert resp.json()['balance'] == expected
    check_balance(0)
    logging.info(f'Balance = 0')

    resp = requests.post(f'http://{PAYMENT_SERVICE_ADDR}/balance/{username}/replenish', headers=auth_headers, json={
        'amount': 500,
    })
    assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
    check_balance(500)
    logging.info('Replenish balance: 500')
    logging.info('Balance = 500 (+500)')


    def check_orders(expected_items: List[int], successful_list: List[bool]):
        resp = requests.get(f'http://{ORDER_SERVICE_ADDR}/orders', headers=auth_headers)
        assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
        orders = resp.json()
        orders.sort(key=lambda order: order['creation_ts'])
        assert len(orders) == len(expected_items)
        for order, expected_item_id, successful in zip(resp.json(), expected_items, successful_list):
            assert order['username'] == username
            assert order['item_id'] == expected_item_id
            assert order['successful'] == successful
    check_orders([], [])
    logging.info('Orders: []')


    def update_warehouse(item_id: str, in_stock: int):
        resp = requests.post(f'http://{WAREHOUSE_SERVICE_ADDR}/items/{item_id}', json={
            'item_id': item_id,
            'in_stock': in_stock,
        })
        assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'

    def check_warehouse(item_id: str, in_stock: int):
        resp = requests.get(f'http://{WAREHOUSE_SERVICE_ADDR}/items/{item_id}')
        assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
        item = resp.json()
        assert item['item_id'] == item_id
        assert item['in_stock'] == in_stock

    logging.info(f'http://{WAREHOUSE_SERVICE_ADDR}/health')
    resp = requests.get(f'http://{WAREHOUSE_SERVICE_ADDR}/health')
    assert resp.status_code == 200
    update_warehouse('cheap_item_id', 1)
    check_warehouse('cheap_item_id', 1)
    logging.info('')
    logging.info('Warehouse update: +item(id="cheap_item_id", in_stock=1)')
    logging.info('Warehouse: "cheap_item_id", 1')


    def update_courier(courier_id: str, status: str):
        resp = requests.post(f'http://{DELIVERY_SERVICE_ADDR}/couriers/{courier_id}', json={
            'id': courier_id,
            'status': status,
        })
        assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'

    def check_courier(courier_id: str, status: str):
        resp = requests.get(f'http://{DELIVERY_SERVICE_ADDR}/couriers/{courier_id}')
        assert resp.status_code == 200,  f'code={resp.status_code}, text={resp.text}'
        courier = resp.json()
        assert courier['id'] == courier_id
        assert courier['status'] == status

    for courier_id in ['Alice', 'Bob', 'Charlie', 'Daniel']:
        update_courier(courier_id, 'available')
        check_courier(courier_id, 'available')
        logging.info(f'Couriers update: +{courier_id} (available)')


    resp = requests.post(f'http://{ORDER_SERVICE_ADDR}/orders', headers=auth_headers, json={
        'item_id': 'cheap_item_id',
    })
    assert resp.status_code == 200, f'code={resp.status_code}, text={resp.text}'
    courier_id = resp.json()['courier_id']
    check_courier(courier_id, 'busy')
    check_warehouse('cheap_item_id', 0)
    check_orders(['cheap_item_id'], [True])
    check_balance(400)
    logging.info('')
    logging.info('Buy item(id="cheap_item_id", cost=100)')
    logging.info('Warehouse: "cheap_item_id", 0')
    logging.info(f'Delivery: {courier_id} (busy)')
    logging.info('Orders: ["cheap_item_id"(good)]')
    logging.info('Balance = 400 (-100)')


    update_warehouse('expensive_item_id', 1)
    check_warehouse('expensive_item_id', 1)
    logging.info('')
    logging.info('Warehouse update: item(id="expensive_item_id", in_stock=1)')
    logging.info('Warehouse: "cheap_item_id", 0 | "expensive_item_id", 1')


    resp = requests.post(f'http://{ORDER_SERVICE_ADDR}/orders', headers=auth_headers, json={
        'item_id': 'expensive_item_id',
    })
    assert resp.status_code == 403
    check_warehouse('expensive_item_id', 1)
    check_orders(['cheap_item_id', 'expensive_item_id'], [True, False])
    check_balance(400)
    logging.info('')
    logging.info('Buy item(id="expensive_item_id", cost=1000) -> ERROR (not enough money)')
    logging.info('Warehouse: "cheap_item_id", 0 | "expensive_item_id", 1')
    logging.info('Orders: ["cheap_item_id"(good), "expensive_item_id"(bad)]')
    logging.info('Balance = 400 (-0)')


    resp = requests.post(f'http://{ORDER_SERVICE_ADDR}/orders', headers=auth_headers, json={
        'item_id': 'cheap_item_id',
    })
    assert resp.status_code == 403
    check_warehouse('cheap_item_id', 0)
    check_orders(['cheap_item_id', 'expensive_item_id', 'cheap_item_id'], [True, False, False])
    check_balance(400)
    logging.info('')
    logging.info('Buy item(id="cheap_item_id", cost=100) -> ERROR (not enough in stock)')
    logging.info('Warehouse: "cheap_item_id", 0 | "expensive_item_id", 1')
    logging.info('Orders: ["cheap_item_id"(good), "expensive_item_id"(bad), "cheap_item_id"(bad)]')
    logging.info('Balance = 400 (-0)')


    for courier_id in ['Alice', 'Bob', 'Charlie', 'Daniel']:
        update_courier(courier_id, 'busy')
    update_warehouse('cheap_item_id', 1)
    check_warehouse('cheap_item_id', 1)
    logging.info('')
    for courier_id in ['Alice', 'Bob', 'Charlie', 'Daniel']:
        logging.info(f'Couriers update: {courier_id} (busy)')
    logging.info('Warehouse update: item(id="cheap_item_id", in_stock=1)')
    logging.info('Warehouse: "cheap_item_id", 1 | "expensive_item_id", 1')


    resp = requests.post(f'http://{ORDER_SERVICE_ADDR}/orders', headers=auth_headers, json={
        'item_id': 'cheap_item_id',
    })
    assert resp.status_code == 403
    check_warehouse('cheap_item_id', 1)
    check_orders(['cheap_item_id', 'expensive_item_id', 'cheap_item_id', 'cheap_item_id'], [True, False, False, False])
    check_balance(400)
    logging.info('')
    logging.info('Buy item(id="cheap_item_id", cost=100) -> ERROR (no available courier)')
    logging.info('Warehouse: "cheap_item_id", 1 | "expensive_item_id", 1')
    logging.info('Orders: ["cheap_item_id"(good), "expensive_item_id"(bad), "cheap_item_id"(bad), "cheap_item_id"(bad)]')
    logging.info('Balance = 400 (-0)')


def remove_test_usernames(con: psycopg2.extensions.connection):
    query = f'''
DELETE FROM "user" WHERE username LIKE '{TEST_USER_PREFIX}%';
DELETE FROM "user_balance" WHERE username LIKE '{TEST_USER_PREFIX}%';
DELETE FROM "order" WHERE username LIKE '{TEST_USER_PREFIX}%';
'''
    with con.cursor() as cursor:
        cursor.execute(query)
