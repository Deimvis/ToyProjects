# Order Creation Algorithm

## Explanation

_Выбран Two-phase commit над Saga pattern_

​	Самое значимое преимущество Saga pattern заключается в его возможности развязать сервисы по времени и тем самым эффективно реализовать длительные транзакции, упростить интерфейс общения и увеличить throughput, но при этом сложность реализации значительно вырастает, архитектура также усложняется (за счет добавления очередей) и появляется длительная неконсистентность между данными (eventual consistency).

​	Для данной системы Saga pattern это полноценный overkill: предполагается, что транзакция будет выполняться относительно быстро, не ожидается чрезмерно высокий RPS, и соответственно незачем жертвовать простотой и консистентностью.

## Scenarios

### Normal Flow

<img src="images/2pc_normal.png" style="zoom:15%">

### Not enough in stock

<img src="images/2pc_not_enough_in_stock.png" style="zoom:15%">

### No available courier

<img src="images/2pc_no_available_courier.png" style="zoom:15%">



