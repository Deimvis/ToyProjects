import os


PAYMENT_SERVICE_ADDR = os.getenv('PAYMENT_SERVICE_ADDR')
assert PAYMENT_SERVICE_ADDR is not None
WAREHOUSE_SERVICE_ADDR = os.getenv('WAREHOUSE_SERVICE_ADDR')
assert WAREHOUSE_SERVICE_ADDR is not None
DELIVERY_SERVICE_ADDR = os.getenv('DELIVERY_SERVICE_ADDR')
assert DELIVERY_SERVICE_ADDR is not None
