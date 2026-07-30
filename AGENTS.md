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
- 不在运行时代码里做旧版兼容处理，不写隐式读取旧 key、旧字段、旧存储结构并自动迁移的兼容补丁。旧版升级需要额外做数据迁移脚本、一次性升级脚本或明确的升级任务；业务代码保持当前版本模型纯净。

## 任务边界和验收 SOP

当执行 P1/Post-P1 或其他阶段性任务时，不能只按技术模块是否接入来判断完成，必须按用户可见能力和验收标准判断完成。尤其是任务名包含“支持”“架构”“闭环”“国际化”“多账号”“支付”“实时”等词时，默认存在隐含业务验收，不允许用“基础设施已完成”替代“能力已完成”。

### 1. 开始任务前：先澄清 Done 的定义

- 读取总计划、当前阶段计划、相关代码后，先把任务拆成四类内容：
  - 已有能力：当前代码已经支持什么。
  - 目标能力：本任务结束后用户/管理员/系统必须能做什么。
  - 非目标：明确不在本阶段解决什么，避免后续误判。
  - 验收标准：用具体场景描述如何证明完成。
- 如果计划文字存在多种解释，先向用户确认，不要自行选择较窄解释后直接标记完成。
- 对“支持 X”这类表述，必须追问或主动补齐 X 的实际使用场景。例如“支持多账号”要区分：
  - 只保存多套账号，不互相覆盖。
  - 同一 tab 内切换账号。
  - 多个 tab 同时使用不同账号。
- 对“国际化/i18n”这类表述，必须区分：
  - 只接入 i18n 基础设施和默认语言。
  - 提供两套语言资源。
  - 页面有语言切换入口。
  - 前端组件库和后端错误消息跟随语言切换。

### 2. 写计划时：每个任务必须包含验收用例

- 每个任务至少写 3 类验收：
  - 正向主流程：用户按预期操作能完成目标。
  - 隔离/边界流程：不同账号、不同状态、不同角色、不同 tab、不同语言等不会串数据。
  - 回归流程：旧数据、旧 localStorage key、旧接口响应、旧订单/资金状态的升级或处理策略必须明确；不要在运行时代码里隐式兼容旧结构。
- 验收用例必须尽量写成可执行句子，例如：
  - “tab A 使用账号 A 下单时，请求 token 仍为账号 A；tab B 使用账号 B 操作时，请求 token 仍为账号 B。”
  - “切换到 en-US 后，用户端导航、后台菜单、Element Plus 组件文案、业务错误提示都显示英文，并且后端返回的 status code 不变。”
- 如果某个能力只完成了基础设施，任务状态必须写成 `Partial` 或 `Reopened`，并列出剩余验收项；不能写 `Completed`。
- `Completed` 只能用于计划中的目标能力和验收标准都已实现、已验证、且已说明已知限制的任务。

### 3. 实施中：持续对照验收标准

- 每完成一个子能力，回到计划核对它满足的是哪条验收标准。
- 如果实现过程中发现计划漏了关键用户场景，必须先更新计划，再继续实现。
- 如果为了控制范围只做基础设施，必须在计划和最终回复里明确“当前不支持什么”，并把缺口记录为本阶段剩余任务或后续阶段任务。
- 不要因为测试/构建通过就认为业务完成。测试/构建只证明实现没有明显编译或自动化回归问题，不证明任务验收完整。

### 4. 完成前：执行完成度审查

- 最终回复和提交前，逐条回答：
  - 计划里的每条目标是否都完成？
  - 用户可见入口是否存在？
  - 默认数据、旧数据、异常数据的处理策略是否明确？
  - 多角色、多账号、多 tab、多语言、多状态等边界是否验证？
  - 是否存在“架构已接入但用户无法使用”的情况？
  - 是否有明确未完成项？如果有，任务不能标记 `Completed`。
- 对前端体验类任务，必须从用户操作路径验证，而不只是看组件或 store 是否存在。
- 对资金/支付/结算类任务，必须验证金额流向、流水、幂等和状态转移，不只看接口返回成功。
- 对通知/实时/订阅类任务，必须验证创建数据、推送信号、前端刷新、断线/回退策略。

### 5. 状态和汇报口径

- 汇报时区分：
  - “基础设施已完成”：底层能力或代码框架已经有了。
  - “主流程可用”：最常见路径可操作。
  - “验收完成”：计划列出的主流程、边界流程、回归流程都通过。
- 最终回复必须包含已验证命令或手工验收场景，也必须包含已知不支持项。
- 用户指出“没有做完”时，先复核任务边界和验收标准，不要只从代码实现角度辩解；如果确实是边界定义不清，立即更新计划/SOP，再补实现。

## Git 和敏感文件

- 不要提交本地配置、密钥、日志、数据库/缓存数据。
- `config/locales/config.yaml` 是本地配置，已经进入 `.gitignore` 后仍可能因为被 Git 跟踪而出现在状态里；需要用 `git rm --cached config/locales/config.yaml` 停止跟踪。
- `data/redis` 是 Redis 本地数据，若已经被跟踪，需要用 `git rm -r --cached data/redis` 停止跟踪。
- 如果敏感配置已经推送到远程，先更换密钥/密码，再考虑清理 Git 历史。
- 不要使用 `git reset --hard`、`git checkout --` 等会丢改动的命令，除非用户明确要求。

## 后端规范

- 后端使用 Go + Gin + GORM + MySQL。
- 后端按 `api/v1 -> application/usecase -> domain/event + repository/dao` 的方向组织跨模块业务。`api/v1` 可以调用 `application` 或单一模块的 `service`；涉及支付、退款、结算、资金、通知、索引等跨模块流程时，优先放到 `application` 编排。
- `service` 和 `service` 之间不能直接互相调用。禁止在 `service/*.go` 的业务方法里调用其他 `GetXxxSrv()`；需要跨模块强一致协作时放到 `application`，需要通知、搜索索引、邮件等副作用时发布 `domain/event` 事件。
- `application` 和 `domain/event` 不允许 import `service` 包。`application` 负责事务边界和用例编排，可直接调用 DAO、纯领域 helper、事件发布器；`domain/event` handler 只能调用 DAO、repository adapter 或基础设施客户端。
- 当前轻量事件总线是进程内同步发布，只用于 Post-P1 解耦通知/索引等副作用；可靠 outbox、RabbitMQ 重试、死信队列和补偿平台仍属于 A4，不要在 Post-P1 中提前扩大范围。
- 支付成功、退款审核、结算打款这类资金强一致流程必须由 `application` 统一控制 GORM 事务；通知、ES 索引、RabbitMQ 发布等副作用应在事务成功后触发。
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
- 前端新增或修改任何用户可见文本、按钮、菜单、表头、占位符、弹窗、Toast/Message、校验提示、空状态、状态标签时，必须同时写入 `zh-CN` 和 `en-US` locale 资源，并在页面中通过 i18n key 使用；不要在 `.vue/.ts` 业务代码里硬编码中文提示。业务数据本身，如商品名、分类名、店铺名、用户输入内容，不属于此约束。

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
