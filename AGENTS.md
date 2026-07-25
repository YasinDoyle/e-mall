# AGENTS.md

本文件是给 Codex/Agent 使用的项目工作规范。新会话开始时，先读本文件，再读当前任务相关 plan 和代码。

## 项目定位

- 这是一个电商商城项目，目标不是只堆功能，而是逐步形成平台型商城的业务闭环。
- 总路线图在 `.github/prompts/plan-eMallFullDevelopmentPlan.prompt.md`。
- P1 平台商业闭环执行计划在 `docs/superpowers/plans/2026-07-24-p1-platform-business-loop.md`。
- 当前推荐演进顺序：先准生产电商主链路，再可运营商业闭环，最后生产级架构练习。
- 商城账号采用统一用户体系。用户不是注册时区分买家/卖家，而是通过商家入驻获得卖家能力。

## 工作方式

- 动代码前先读现有实现，优先沿用当前项目分层：`api/v1`、`service`、`repository/db/dao`、`repository/db/model`、`types`、`routes`、`consts`。
- 手工编辑文件使用 `apply_patch`。
- 不要回滚用户已有改动。工作区可能是脏的，先用 `git status --short` 看清楚。
- 大功能先拆 plan，再按小任务实现。P1 后续优先按 `docs/superpowers/plans/2026-07-24-p1-platform-business-loop.md` 执行。
- Review 时必须注意未跟踪文件。`git diff <base>` 不会显示 untracked 文件，若代码引用了 untracked 新文件，提交后可能编译失败。

## Git 和敏感文件

- 不要提交本地配置、密钥、日志、数据库/缓存数据。
- `config/locales/config.yaml` 是本地配置，已经进入 `.gitignore` 后仍可能因为被 Git 跟踪而出现在状态里；需要用 `git rm --cached config/locales/config.yaml` 停止跟踪。
- `data/redis` 是 Redis 本地数据，若已经被跟踪，需要用 `git rm -r --cached data/redis` 停止跟踪。
- 如果敏感配置已经推送到远程，先更换密钥/密码，再考虑清理 Git 历史。
- 不要使用 `git reset --hard`、`git checkout --` 等会丢改动的命令，除非用户明确要求。

## 后端规范

- 后端使用 Go + Gin + GORM + MySQL。
- GORM 全局配置启用了 `SingularTable: true`，表名通常是单数，如 `user`、`order`、`cart`、`favorite`。
- `NewDBClient(ctx)` 必须返回带 `NewDB: true` 的 session，避免查询状态串到下一次查询。
- 如果 GORM 查询使用表别名，不要写 `Model(&model.X{}).Joins("AS x ...")` 后依赖自动软删除条件；GORM 可能生成原表名的 `deleted_at` 条件，导致 `Unknown column '<table>.deleted_at'`。
- 使用别名查询时优先：

```go
db.Table("cart AS c").
    Joins("LEFT JOIN product AS p ON c.product_id = p.id").
    Where("c.deleted_at IS NULL")
```

- `order` 是 MySQL 关键字，手写 SQL 查询表结构时要用反引号：`DESC \`order\`;`。
- 钱包/支付相关逻辑要特别注意幂等、金额校验、重复回调、退款冲正。
- 支付密码不应在注册时设置。当前流程是注册后在钱包页设置支付密码。
- 卖家上架商品必须满足商家审核通过和已设置支付密码；下架商品不应被支付密码拦住。
- P1 商家边界生效后，发布商品、上架商品、Admin 审核商品都应校验卖家身份。

## 后端测试

- 后端改动后默认运行：

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./...
```

- 只改 DAO 时可先运行：

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./repository/db/dao
```

- 新增或修改业务规则时优先补小的纯逻辑测试，例如 service 校验函数、状态机函数、金额计算函数。
- 对 GORM SQL 生成问题，可使用 MySQL Dialector + `DryRun: true` + `DisableAutomaticPing: true`，避免依赖真实数据库。

## 用户端 web 规范

- 用户端前端在 `web/`，技术栈是 Vue 3 + TypeScript + Vite + Element Plus + Pinia + Vue Router + Axios。
- 修改 `web/` 后运行：

```bash
cd web
npm run build
```

- 用户端是商城主站，应该全屏铺开。不要在全局 `#app` 上设置固定宽度、居中边框或全局 `text-align: center`。
- 页面内容可以在各自页面里使用内容最大宽度，但根容器必须 `width: 100%`。
- 登录态从 Pinia `userStore` 获取，接口通过 `web/src/utils/request.ts` 统一带 token 和处理错误。
- 注册流程使用邮箱验证码：发送验证码、邮箱验证码、密码、确认密码。不在注册页设置支付密码。
- 前端接口字段以后端 JSON 字段为准。订单相关字段统一使用后端需要的 `order_id`，不要混用 `id` 作为订单操作字段。
- 图片上传用 `FormData`。通常不需要手写 multipart boundary；浏览器会根据 `FormData` 自动生成。

## 管理后台 web-admin 规范

- 管理后台在 `web-admin/`，技术栈和用户端类似，额外使用 ECharts。
- 修改 `web-admin/` 后运行：

```bash
cd web-admin
npm run build
```

- 后台是运营工具，不做营销落地页风格。界面应偏工作台：信息密度适中、表格/筛选/操作清晰。
- Admin 接口挂 `/api/v1/admin/*`，需要登录和 `IsAdmin` 权限。
- 新增后台模块时要同时考虑列表筛选、详情查看、审批/操作结果提示和错误回显。
- 商品、商家、订单、售后、结算这些后台模块要优先保证状态清晰，不要只给按钮没有状态解释。

## 产品和业务规则

- 平台目标是类似淘宝/京东的平台型商城，但按阶段演进。
- 账号统一：买家和卖家是同一个用户账号的不同能力。
- 商家入驻是 P1 的基础。通过 `SellerProfile` 表示店铺资料、审核状态、拒绝原因、通过时间。
- 平台盈利第一阶段只做订单佣金，不急着做广告费、活动报名费、推荐坑位等复杂收费。
- 订单履约必须闭环：下单、支付、卖家发货、买家收货、售后退款、商家结算。
- 资金相关必须有流水：用户流水、商家流水、平台流水。每次支付、退款、佣金入账、结算出账都要可追溯。
- 秒杀不能只做按钮和接口，要明确容量边界：并发请求数、库存一致性、防重复下单、异步建单结果查询。
- 搜索和推荐属于生产级专题。先让主链路和数据闭环跑通，再做 ES 规模化、埋点、推荐召回和排序。

## 当前 P1 进度提示

- P1 Task 1-4 后端基础已开始实现：商家资料、用户申请、Admin 审核、商品链路商家校验。
- P1 后续优先做：
  - Task 5：用户端卖家中心。
  - Task 6：Admin 商家审核页面。
  - Task 7：佣金与资金流水。
  - Task 8：结算后台操作。
- 继续 P1 前先检查 `git status --short`。当前可能存在未跟踪的新后端文件，提交或审查时不要漏掉。

## 常见坑

- `.gitignore` 不会自动忽略已经被 Git 跟踪的文件，需要 `git rm --cached`。
- `git diff <base>` 不显示 untracked 文件，review 和提交前必须看 `git status --short`。
- GORM `WithContext` 不等于重置 session；需要 `NewDB: true` 避免 statement 残留。
- `Model + alias join + soft delete` 容易生成错误的 `<table>.deleted_at` 条件。
- MySQL 关键字表名如 `order` 手写 SQL 必须加反引号。
- 发送邮件如果报 `dial tcp :465`，通常是 SMTP host/port 配置不完整。
- 本地静态头像默认使用 `avatar.jpg`，注意大小写，不要再引用 `avatar.JPG`。
