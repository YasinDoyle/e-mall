# 电商系统完整开发计划（后端补全 + 前端双工程）

## 架构说明

- 后端：现有 Go + Gin 项目（`e-mall/`）
- 用户端前端：`web/` — 独立 Vite 工程，Nginx 独立部署
- 管理后台：`web-admin/` — 独立 Vite 工程，独立子域名部署
- 开发时两端均通过 `vite.config.ts` proxy `/api/v1` → `localhost:8080`，不依赖后端 `static/`

---

## Part 1 — 后端功能补全

**Phase B1 — Admin 角色体系** `约2天`

1. `repository/db/model/user.go` 新增 `IsAdmin bool`，DB migrate
2. 新建 `AdminAuthMiddleware`（校验 JWT + IsAdmin）
3. 分类管理 CRUD：`/api/v1/admin/category`
4. 轮播图管理 CRUD：`/api/v1/admin/carousel`
5. 公告管理 CRUD：`/api/v1/admin/notice`（`Notice` 模型已有，补 DAO + Service + API）
6. 用户列表/封禁：`/api/v1/admin/users`
7. 商品模型加 `Status` 字段，管理员审核上架接口

**Phase B2 — 商品评价系统** `约1-2天`

1. 新建 `model/review.go`：`{UserID, ProductID, OrderID, Rating uint, Content, Images}`
2. 新建 `dao/review.go`、`service/review.go`、`api/v1/review.go`
3. 路由：`POST reviews/create`（仅 OrderType=已收货可调用）、`GET product/reviews`、`DELETE admin/review/:id`

**Phase B3 — 优惠券系统** `约2-3天`

1. 新建 `model/coupon.go`（类型/折扣/门槛/库存/过期时间）、`model/user_coupon.go`（领券记录）
2. 路由：`POST admin/coupon/create`、`POST coupon/claim`、`GET coupon/list`
3. 修改 `service/order.go` `OrderCreate`，增加优惠券核销逻辑

**Phase B4 — 退款/售后** `约1-2天，依赖 B5`

1. `repository/db/model/order.go` 新增 `RefundStatus int`、`RefundReason string`、`TrackingNo string`
2. `consts/order.go` 新增 `OrderTypeRefundRequested`、`OrderTypeRefunded`
3. 路由：`POST orders/refund/request`（用户）、`POST admin/orders/refund/approve`（管理员）
4. 退款成功后调用 B5 的退款接口

**Phase B5 — 微信/支付宝支付网关** `约3-5天，可并行 B1-B4`

SDK：`github.com/go-pay/gopay`（同时支持微信V3 和支付宝）

架构：充值 → 钱包余额 → 扣款（保留现有余额支付逻辑）

1. `config/config.go` 新增两个配置块：
   - `WechatPay`：AppID、MchID、ApiV3Key、SerialNo、PrivateKey 路径
   - `Alipay`：AppID、PrivateKey、AlipayPublicKey
2. 新建 `service/payment_gateway.go`：
   - `WechatPayRecharge(orderNum, amount)` → Native Pay → 返回二维码 URL
   - `AlipayRecharge(orderNum, amount)` → PC/WAP Pay → 返回跳转链接
   - `WechatPayRefund()` / `AlipayRefund()`
3. 回调路由（不挂 JWT 中间件）：
   - `POST /api/v1/pay/wechat/notify`：验签 → 更新 `user.Money` → 发 MQ 消息
   - `POST /api/v1/pay/alipay/notify`：验签 → 更新 `user.Money`
4. 充值发起路由（挂 JWT）：
   - `POST /api/v1/recharge/wechat`、`POST /api/v1/recharge/alipay`
   - `GET /api/v1/recharge/status`（前端轮询充值结果）

**Phase B6 — Admin 统计 API** `约1天`

1. `GET /api/v1/admin/stats/overview`：今日订单数、总销售额、注册用户数
2. `GET /api/v1/admin/stats/orders`：时间段内订单趋势（ECharts 折线图数据）

---

## Part 2 — Vue3 用户端 `web/`

**技术栈：** Vite 5 · Vue3 · TypeScript · Element Plus · Pinia · Vue Router 4 · Axios

---

**Phase F1 — 工程初始化** `0.5天`

1. `web/` 下 `npm create vite@latest . -- --template vue-ts`
2. 安装：Element Plus、Pinia、Vue Router 4、Axios
3. `vite.config.ts`：dev proxy `/api/v1` → `http://localhost:8080`
4. 目录：`src/api/`、`src/views/`、`src/components/`、`src/stores/`、`src/types/`
5. Axios 封装：baseURL、JWT 自动注入、401 跳登录、统一 Error toast
6. 按后端路由分模块封装全部 API 函数（`api/user.ts`、`api/product.ts` 等）

**Phase F2 — 用户认证** `1天，依赖 F1`

1. 登录页 `/login`、注册页 `/register`、邮箱验证页 `/valid-email`
2. Pinia `userStore`（token + userInfo，localStorage 持久化）
3. 全局路由守卫（未登录跳 `/login`）
4. 顶部 `NavBar.vue`（头像下拉、购物车角标、退出登录）

> **可验收：** 注册 → 登录 → 顶部显示昵称 → 退出，流程完整

**Phase F3 — 首页 & 商品浏览** `2-3天，依赖 F1，可并行 F2`

1. 首页 `/`：轮播图、分类导航、商品网格、秒杀预告入口
2. 商品列表页 `/products`：分页 + 分类 Tag 筛选
3. 商品详情页 `/product/:id`：多图切换、加入购物车、收藏、评价列表占位
4. 搜索页 `/search`（ES 接口）
5. 公共组件：`ProductCard.vue`、`Pagination.vue`

> **可验收：** 首页刷出商品 → 点进详情 → 搜索找到商品

**Phase F4 — 购物车 & 下单** `2天，依赖 F3`

1. 购物车页 `/cart`：数量增减、删除、全选、价格汇总
2. 结算页 `/checkout`：收货地址选择（无地址引导新建）、优惠券占位、订单摘要
3. 支付页 `/payment`：余额支付（调 `/api/v1/paydown`）+ 微信/支付宝入口占位
4. 订单成功页

> **可验收：** 加购 → 结算 → 余额支付 → 订单成功，主链路跑通

**Phase F5 — 用户中心** `2天，依赖 F2，可并行 F3`

1. 用户中心 `/user` 布局（左侧导航）
2. 个人资料：编辑昵称/邮箱、头像上传
3. 订单列表 `/user/orders`：状态 Tab（待支付/待发货/已发货/已完成）、确认收货、删除
4. 订单详情 `/user/orders/:id`：物流状态时间轴、写评价入口占位
5. 收货地址 CRUD
6. 收藏夹：列表、取消收藏
7. 钱包页：余额展示 + 充值入口占位

**Phase F6 — 秒杀专场** `1天，依赖 F3`

1. 秒杀列表页 `/flash-sale`：商品卡片 + 倒计时组件
2. 秒杀详情 `/flash-sale/:id`：库存显示、抢购按钮（loading + 防重复提交）

**Phase F7 — 评价 & 优惠券** `1天，依赖 F5 + B2/B3 完成`

1. 订单详情页接入写评价：星级 + 文字 + 图片上传
2. 用户中心"我的优惠券"：可用/已用/已过期三 Tab
3. 结算页接入优惠券选择弹窗

**Phase F8 — 支付宝/微信充值** `1天，依赖 B5 完成`

1. 钱包页充值弹窗：
   - **微信**：展示二维码，每 2 秒轮询 `/recharge/status`
   - **支付宝**：新页签跳转，返回后自动刷新余额
2. 支付页接入真实支付方式

---

## Part 3 — Vue3 Admin 后台 `web-admin/`

**技术栈：** 同用户端 + ECharts · 独立部署 · 独立子域名

**Phase F9 — Admin 后台** `5-7天，依赖 B1 完成`

1. 独立 Vite 工程初始化（加 ECharts）
2. Admin 专属登录页（`isAdmin` 校验）
3. Dashboard：统计卡片（订单量/销售额/用户数）+ ECharts 折线图（对接 B6）
4. 商品管理：列表/新增/编辑（多图上传）/审核上架/下架
5. 分类管理 CRUD
6. 轮播图管理 CRUD（含图片上传预览）
7. 用户管理：列表/封禁/解封
8. 订单管理：全量列表/退款申请审批（通过/拒绝）
9. 优惠券管理：创建/下线
10. 秒杀管理：新增/编辑秒杀商品
11. 公告管理 CRUD

---

## 整体节奏

```
Week 1:   F1 + F2 + F3         → 首页、商品浏览可看效果
Week 2:   F4 + F5              → 购物主流程跑通
Week 3:   F6 + F7              → 秒杀 + 评价 + 优惠券
后端穿插:  B1-B3 完成 → F7 联调
           B5 完成   → F8 联调
Week 4+:  F9 Admin 后台
```

| 部分                | 预估工时      |
| ------------------- | ------------- |
| 后端补全（B1-B6）   | 10-15天       |
| 用户端前端（F1-F8） | 10-12天       |
| Admin 后台（F9）    | 5-7天         |
| **合计**            | **约25-34天** |

---

## 关键文件索引

| 文件                           | 变更内容                                   |
| ------------------------------ | ------------------------------------------ |
| `repository/db/model/user.go`  | 新增 `IsAdmin bool`                        |
| `repository/db/model/order.go` | 新增 `RefundStatus`、`TrackingNo`          |
| `service/payment.go`           | 现有余额支付逻辑保留                       |
| `service/payment_gateway.go`   | 新建，gopay SDK 封装                       |
| `config/config.go`             | 新增 WechatPay / Alipay 配置块             |
| `routes/router.go`             | 新增 `/admin` 路由组、`/pay/*/notify` 回调 |
| `web/vite.config.ts`           | proxy 配置（dev 环境）                     |
| `web-admin/vite.config.ts`     | proxy 配置（dev 环境）                     |
