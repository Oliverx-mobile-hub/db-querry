# Research：MySQL 数据库查询支持

**日期**：2026-07-22

**关联计划**：[plan.md](./plan.md)

## 决策 1：MySQL 作为第二种数据库类型，与 PostgreSQL 并存

**Decision**：`databaseType` 扩展为 `postgres | mysql`，所有连接、metadata、SQL Guard、query executor 和 LLM prompt 都按数据库类型分发。

**Rationale**：当前系统内部已有 `database_type` 存储字段，但代码硬编码 `postgres`。把 databaseType 变成真实分发维度，可以支持 MySQL 且保护 PostgreSQL 回归。

**Alternatives considered**：

- 只从 URL scheme 推断类型：对 MySQL DSN 不稳定，前端也无法清楚展示数据库类型。
- 创建新的 `/api/v1/mysql/...` API：会破坏现有 API 风格，增加前端分支。

## 决策 2：PUT 连接请求增加可选 databaseType

**Decision**：`PUT /api/v1/dbs/{name}` 请求体支持 `{ "url": "...", "databaseType": "mysql" }`。旧请求不传时保持兼容，由 URL 推断或默认 `postgres`。

**Rationale**：显式 databaseType 可以避免 MySQL URL/DSN 格式歧义，同时保持旧 PostgreSQL 用法不变。

**Alternatives considered**：

- 前端只让用户输入 URL：实现简单但用户不清楚当前连接类型。
- 根据端口 3306/5432 推断：不可靠，端口可以自定义。

## 决策 3：MySQL metadata 使用 information_schema

**Decision**：MySQL metadata collector 使用 `information_schema.tables`、`information_schema.columns`、`information_schema.table_constraints`、`information_schema.key_column_usage`。

**Rationale**：MySQL 的 `information_schema` 可以稳定采集 table/view、column、data type、nullable、primary key 和 comment，满足当前 UI 展示需求。

**Alternatives considered**：

- 使用 `SHOW TABLES` / `SHOW FULL COLUMNS`：实现快，但需要多轮查询，结构化和排序不如 information_schema。
- 用 LLM 整理 metadata：不可靠，不能作为数据来源。

## 决策 4：MySQL 查询使用 database/sql + go-sql-driver/mysql

**Decision**：新增 MySQL connector/query executor 使用 `database/sql` 和 `github.com/go-sql-driver/mysql`，并把浏览器输入的 MySQL URL 转换为 driver DSN。

**Rationale**：Go MySQL driver 成熟稳定；`database/sql` 的 `Rows.Columns` 和 `ColumnTypes` 可以支持现有 QueryResult。

**Alternatives considered**：

- 直接让用户输入 driver DSN：不如 URL 直观，也不符合现有 DB URL 输入习惯。
- 统一所有数据库都走 database/sql：PostgreSQL 现有 pgx 路径已经可用，没必要在本功能中重写。

## 决策 5：SQL Guard 按数据库类型分发

**Decision**：SQL Guard 接口需要感知 `databaseType`。PostgreSQL 继续走现有规则；MySQL 增加 MySQL 方言校验，至少保证单条只读 SELECT、拒绝修改性语句和多语句、默认 LIMIT 1000。

**Rationale**：MySQL 和 PostgreSQL 引号、函数、LIMIT、CTE 支持和 parser 行为不同，不能长期混用同一套规则。

**Alternatives considered**：

- 继续使用现有通用正则校验：短期可运行，但会误判 MySQL 方言，且难以支撑后续扩展。
- 依赖数据库执行错误作为校验：不能替代安全边界。

## 决策 6：LLM prompt 显式包含 databaseType 和 MySQL 方言要求

**Decision**：当 `metadata.databaseType=mysql` 时，LLM system/user prompt 必须明确要求生成 MySQL SELECT SQL；生成 SQL 仍需要 MySQL SQL Guard 校验。

**Rationale**：当前 LLM prompt 主要按 PostgreSQL 场景设计。MySQL 需要不同标识符引用和函数习惯。

**Alternatives considered**：

- 只给 metadata 不说明数据库类型：模型可能生成 PostgreSQL 方言 SQL。
- LLM 输出后直接执行：违反只读安全原则。

## 决策 7：前端保留现有 UI 风格，只增加数据库类型选择

**Decision**：`ConnectionPanel` 增加数据库类型选择，默认 PostgreSQL；选择 MySQL 时 placeholder 切换为 MySQL URL 示例，连接列表展示 databaseType。

**Rationale**：用户要求参考之前样式；最小 UI 改动可以降低回归风险。

**Alternatives considered**：

- 为 MySQL 单独创建连接页面：复杂且不符合当前 dashboard 工作流。
- 只靠 URL 让用户猜类型：可用性差。

## 决策 8：Chrome 验证作为实现阶段的前端验收

**Decision**：后端测试通过后，启动后端和前端，并使用 Google Chrome 插件验证 `interview_db` 添加、metadata、手写 SQL、自然语言 SQL、结果展示。

**Rationale**：MySQL 是端到端功能，单元测试不能完全覆盖真实 UI 和浏览器下载/交互状态。

**Alternatives considered**：

- 只跑 Vitest：无法验证真实后端、MySQL 和 UI 集成。
- 手工口头验证：不可复现。

## 决策 9：日志继续使用 stdout 单行结构化文本

**Decision**：继续使用 `backend/internal/logging/logging.go` 的 stdout 单行格式：`time="..." level=info msg="..."`，MySQL 流程新增日志必须脱敏。

**Rationale**：现有日志已满足 K8s/云日志采集的基本要求，继续复用比引入新日志框架更简单。

**Alternatives considered**：

- 引入 logrus：用户参考实现使用 logrus，但当前项目已有轻量封装；为 MySQL 支持单独引入依赖收益有限。
- 写文件日志：容器场景不利于采集。
