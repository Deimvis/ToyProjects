# E-commerce Service-Oriented App

_System Design hw3_\
Implements:

* 5 Services
  * User Service
  * Payment Service
  * Order Service
  * Warehouse Service
  * Delivery Service
* Auth SDK — Python package that provides unified and simple authentication interface
* 2 Phase Commit for Order Service over Warehouse and Delivery Services

## Functionality

* User Service — Provides CRUD API to user data and is responsible for auth
* Payment Service — Allows to view, withdraw and replenish balance
* Order Service — Allows to create a new order and to view an order history
* Warehouse Service — Allows to update warehouse related info for items and manage item reservations
* Delivery Service — Allows to update and manage couriers info

## Getting Started

```bash
docker-compose up --build
```

* User Service doc: <http://127.0.0.1:8000/docs>
* Payment Service doc: <http://127.0.0.1:8001/docs>
* Order Service doc: <http://127.0.0.1:8002/docs>
* Warehouse Service doc: <http://127.0.0.1:8003/docs>
* Delivery Service doc: <http://127.0.0.1:8004/docs>

## Run tests

```bash
./scripts/docker/run_tests
```

## Documentation

* [Architecture](docs/architecture.md)
* [Order Creation Algorithm (2PC)](docs/order_creation_algorithm.md)
