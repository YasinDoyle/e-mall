## Plan: 补齐阶段 5/6 缺失能力

这份学习指南从阶段 5 开始，确实把“想做的能力”写成了“已经实现的能力”。按当前仓库代码看，阶段 5 目前只有本地事务支付、普通库存扣减、Redis 记录订单超时；阶段 6 目前只有秒杀数据初始化、列表缓存、详情读取，离真正的秒杀链路还差很远。

### 建议的落地顺序

1. 先改文档基线  
先把 [doc/学习指南.md](doc/学习指南.md) 改成“已实现 / 待实现”两栏。否则后面边学边改时，文档会一直误导你。  
知识点：代码审计、文档和实现对齐。  
学习方式：对照 [cmd/main.go](cmd/main.go)、[service/order.go](service/order.go)、[service/payment.go](service/payment.go)、[service/flash_sale.go](service/flash_sale.go)，逐条确认“代码真的做了什么”。

2. 先补订单支付幂等和状态机  
先在 [service/payment.go](service/payment.go) 里补“已支付不能重复扣款”的检查，再在 [repository/db/model/order.go](repository/db/model/order.go) 和 [repository/db/dao/order.go](repository/db/dao/order.go) 里补支付时间、支付状态或等价字段，并把待发货、待收货、已收货的流转补完整。  
知识点：幂等、状态机、事务边界。  
学习方式：先学“重复点击支付”和“回调重放”为什么会重复扣钱，再看幂等键、唯一约束、状态机这三种常见做法。

3. 再补库存并发控制，主路径优先乐观锁  
当前扣库存只是先查再改，见 [service/payment.go](service/payment.go)。建议先在 [repository/db/model/product.go](repository/db/model/product.go) 增加 version 字段，或者直接在 [repository/db/dao/product.go](repository/db/dao/product.go) 里做条件更新，保证库存扣减是原子的。悲观锁可以作为对照实现补进去，但不建议先拿它做主路径。  
知识点：乐观锁、悲观锁、条件更新、防超卖。  
学习方式：先理解“为什么先读再写会超卖”，再分别实现一个乐观锁版和一个 select for update 版，对比吞吐和侵入性。

4. 订单支付后再接 RabbitMQ，但只承担异步副作用  
RabbitMQ 现在只是配置留坑，代码没接，见 [cmd/main.go](cmd/main.go)、[config/config.go](config/config.go)、[config/locales/config.yaml](config/locales/config.yaml)、[go.mod](go.mod)。建议在 [repository/rabbitmq](repository/rabbitmq) 新增基础连接、生产者和最小消费者，然后在 [service/payment.go](service/payment.go) 事务提交成功后发“订单已支付”事件。不要把支付成功本身依赖在 MQ 上。  
知识点：消息队列、事务后事件、最终一致性。  
学习方式：先分清本地事务解决“库内一致性”，MQ 解决“异步解耦”，它们不是一个层面的能力。

5. 先把秒杀缓存模型补完整，再做预热  
当前 [service/flash_sale.go](service/flash_sale.go) 只把列表写进 Redis，详情 key 和库存 key 其实没有完整初始化。建议先补 [repository/cache/key.go](repository/cache/key.go) 的库存 key、详情 key、用户购买标记 key，再在 [service/flash_sale.go](service/flash_sale.go) 里统一做初始化、详情回源、缓存回填和预热。  
知识点：缓存预热、cache aside、热点 key 设计。  
学习方式：先搞清楚列表缓存、详情缓存、库存缓存为什么不能混用，再学预热什么时候值得做。

6. 秒杀核心再补 Redis 预减库存、防重和限流  
在 [service/flash_sale.go](service/flash_sale.go) 里引入 Lua，把“检查库存、检查重复购买、扣减库存、写用户标记”合成一个原子脚本；同时在 [repository/cache/common.go](repository/cache/common.go) 或附近封装脚本执行；限流可以先做简单的 Redis 固定窗口。  
知识点：Redis 原子性、Lua、预减库存、防重复下单、限流。  
学习方式：先单独学 DECR、SETNX、EXPIRE，再学为什么秒杀要把多步逻辑收敛到 Lua 脚本。

7. 最后再接 Kafka 做秒杀异步下单  
当前仓库没有 Kafka 依赖也没有初始化，见 [cmd/main.go](cmd/main.go) 和 [go.mod](go.mod)。建议在 Redis 预减库存和防重稳定后，再新增 [repository/kafka](repository/kafka) 做 producer/consumer，把秒杀成功资格写入消息，消费者再落正式订单。  
知识点：削峰填谷、异步下单、消费幂等、补偿。  
学习方式：先理解“Kafka 解决的是吞吐，不是正确性的第一步”，再学 at least once 和消费幂等。

8. 最后补验证和压测  
给支付幂等、库存并发、秒杀 Lua、防重购买、Kafka 消费幂等补最小测试；再准备并发请求脚本，验证库存不会负数、订单数不会超过库存。  
知识点：并发测试、集成测试、压测。  
学习方式：重点盯“多 goroutine 下库存和订单数是否一致”，不用一开始追求高覆盖率。

### 关键文件

[service/payment.go](service/payment.go) 是阶段 5 的核心入口，要补幂等、库存并发控制、事务后事件。  
[repository/db/model/order.go](repository/db/model/order.go)、[repository/db/dao/order.go](repository/db/dao/order.go)、[repository/db/model/product.go](repository/db/model/product.go)、[repository/db/dao/product.go](repository/db/dao/product.go) 是订单状态和库存锁的主要落点。  
[service/flash_sale.go](service/flash_sale.go)、[repository/db/dao/flash_sale.go](repository/db/dao/flash_sale.go)、[repository/db/model/flash_sale.go](repository/db/model/flash_sale.go)、[repository/cache/key.go](repository/cache/key.go)、[repository/cache/common.go](repository/cache/common.go) 是阶段 6 的主战场。  
[cmd/main.go](cmd/main.go)、[config/config.go](config/config.go)、[config/locales/config.yaml](config/locales/config.yaml)、[go.mod](go.mod) 是 RabbitMQ 和 Kafka 基础设施接入点。

### 推荐学习重点

如果你是“边补项目边学”，我建议顺序固定成这条线：

1. 幂等和订单状态机  
2. 乐观锁防超卖  
3. RabbitMQ 事务后异步事件  
4. 缓存预热和 cache aside  
5. Redis Lua 做秒杀预减库存  
6. 防重复购买和限流  
7. Kafka 异步下单和消费幂等

这样学的原因是：先把“正确性”补齐，再补“高并发”，最后再补“吞吐和解耦”。否则你会很容易陷进 Kafka、Lua、MQ 这些词，但底层的一致性问题其实还没站稳。