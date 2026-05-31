# e-mall

一个基于 Gin + GORM 的 Go 电商后端示例项目，当前已经包含用户、商品、购物车、订单支付、秒杀、ES 搜索、Redis 缓存、RabbitMQ 事件、Kafka 异步下单和基础链路追踪能力。

## 当前状态

目前仓库里已经实现的重点能力：

- 用户注册、登录、资料更新、关注、头像上传、邮箱验证
- 商品、分类、轮播图、购物车、收藏夹、收货地址等常规电商模块
- 订单创建、列表、详情、删除、发货、收货
- 支付幂等、库存并发控制、支付成功 RabbitMQ 事件发布
- 秒杀缓存预热、Redis Lua 预减库存、防重复秒杀、Kafka 异步建单
- Elasticsearch 商品搜索
- Jaeger tracing 基础接入，HTTP -> DB/Redis/MQ 链路可继续向下传递

## 项目结构

```text
cmd/                程序入口
api/v1/             Handler 层
service/            业务层
repository/db/dao/  数据访问层
repository/cache/   Redis
repository/es/      Elasticsearch
repository/rabbitmq/RabbitMQ 发布
repository/kafka/   Kafka 生产和消费
middleware/         JWT、CORS、Tracing
config/             配置和 SQL 初始化脚本
doc/                学习指南、验证说明
scripts/            烟测和压测脚本
```

## 本地启动

### 1. 启动依赖

```powershell
docker-compose up -d mysql redis rabbitmq kafka jaeger elasticsearch kibana
```

如果你暂时只想验证订单支付和秒杀链路，最小依赖可以先起：

```powershell
docker-compose up -d mysql redis rabbitmq kafka
```

### 2. 启动服务

```powershell
cd cmd
go run .\main.go
```

默认 HTTP 端口来自 [config/locales/config.yaml](config/locales/config.yaml)，当前是 `:5001`。

### 3. 初始化顺序

服务启动时会依次执行：

```text
InitConfig -> InitMysql -> InitCache -> InitRabbitMQ -> InitES -> InitKafka -> InitTrack -> Router
```

入口代码在 [cmd/main.go](cmd/main.go)。

## 关键依赖

- MySQL: 主业务数据存储，默认端口 `3307`
- Redis: 缓存、订单超时记录、秒杀库存和用户防重标记
- RabbitMQ: 支付成功事件发布
- Kafka: 秒杀异步建单
- Elasticsearch: 商品搜索
- Jaeger: 链路追踪

主要默认配置见 [config/locales/config.yaml](config/locales/config.yaml)。

## 验证与压测

仓库里已经补了最基础的验证脚本：

- [doc/验证与压测.md](doc/验证与压测.md)
- [scripts/smoke-verify.ps1](scripts/smoke-verify.ps1)
- [scripts/flash-sale-load.ps1](scripts/flash-sale-load.ps1)

它们分别用于：

- 支付幂等与库存烟测
- 秒杀接口并发请求验证

## 学习入口

如果你是拿这个项目练后端分层、事务、缓存、消息队列和秒杀链路，建议直接看：

- [doc/学习指南.md](doc/学习指南.md)
- [routes/router.go](routes/router.go)
- [service/payment.go](service/payment.go)
- [service/flash_sale.go](service/flash_sale.go)
- [repository/kafka/common.go](repository/kafka/common.go)
- [repository/rabbitmq/common.go](repository/rabbitmq/common.go)

## 当前已知限制

- RabbitMQ 目前只补了生产端，没有对应业务消费端
- 秒杀下单是异步流程，当前还没有单独的结果查询接口
- 验证脚本更偏手工验证，不是完整 CI 自动化测试
- `go run .\main.go` 之前需要先把依赖起起来，否则会因为下游服务不可用而启动失败

## 说明

这个仓库更适合用来学习和演示：

- Gin 分层接口设计
- GORM 事务和条件更新
- Redis 缓存和 Lua 原子脚本
- RabbitMQ / Kafka 的最小接入方式
- tracing 上下文传播

如果你要继续完善，比较自然的下一步是：

1. 补秒杀异步结果查询接口
2. 给支付和秒杀链路补更细的 service span
3. 增加更标准的压测脚本，例如 k6
