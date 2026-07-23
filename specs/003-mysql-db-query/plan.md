# 实现计划：MySQL 数据库查询支持

**分支**：`003-mysql-db-query` | **日期**：2026-07-22 | **规格**：[spec.md](./spec.md)

**输入**：来自 `/specs/003-mysql-db-query/spec.md` 的功能规格，以及本轮补充约束：

- 支持新增 MySQL 数据库查询能力。
- 添加数据库连接示例：`interview_db`。
- 支持 MySQL metadata 展示、手写 SQL 查询、自然语言生成 SQL、执行并显示结果。
- 实现 MySQL 相关后端测试；后端测试通过后，再启动后端和前端，使用 Google Chrome 插件验证前端基本流程。
- 保留现有 PostgreSQL 能力不退化。

## 摘要

本功能在现有 PostgreSQL 查询工具基础上增加 MySQL 数据库类型。用户可以添加名为 `interview_db` 的 MySQL 连接，系统按数据库类型采集 metadata、执行只读 SQL、生成自然语言 SQL 草稿并显示查询结果。

核心技术路线：

1. API 保持 `/api/v1/dbs` 命名空间，`PUT /api/v1/dbs/{name}` 增加可选 `databaseType` 字段；旧请求继续兼容。
2. SQLite 已有 `database_type` 字段，后续实现需要把 `postgres`/`mysql` 作为真实分发依据。
3. 后端连接、metadata、query executor、SQL guard 按 `databaseType` 分发。
4. MySQL metadata 使用 `information_schema` 采集 schema/table/view/column/type/nullability/primary key/comment。
5. MySQL SQL 只允许单条只读 `SELECT`，缺少 top-level `LIMIT` 时默认加 `LIMIT 1000`。
6. LLM 生成 SQL 时把 `metadata.databaseType=mysql` 作为上下文，要求输出 MySQL 方言 SQL；生成后仍经过后端 SQL Guard。
7. 前端数据库连接弹窗增加数据库类型选择，UI 风格保持现有 MotherDuck-inspired / brutalist dashboard 风格。
8. 实现完成后先运行后端测试，通过后再启动后端/前端并用 Google Chrome 验证 `interview_db` 的连接、查询、自然语言生成和结果展示。

## 技术上下文

**语言 / 版本**：

- 后端：Go 1.23+。
- 前端：TypeScript + Vue 3。

**主要依赖**：

- 后端 HTTP：现有 Go 标准库 `net/http`。
- 后端 SQLite：现有 SQLite driver 和 `db_connections.database_type` 字段。
- PostgreSQL：保留现有 pgx 路径。
- MySQL：新增 Go MySQL driver，优先使用 `github.com/go-sql-driver/mysql`。
- SQL Guard：抽象为数据库类型感知的 validator；MySQL v1 使用 MySQL 方言 parser 或保守策略组合，必须拒绝非 SELECT、多语句和修改性语句。
- LLM：沿用现有 OpenAI-compatible client 和 `LLM_*` 配置。
- 前端：Vue 3、Element Plus、Vite、TypeScript。
- 浏览器验证：使用 Google Chrome 插件验证前端流程。

**存储**：

- 应用本地 SQLite：`~/.db_querry/db_querrt.db`。
- 连接记录：继续写入 `db_connections`，`database_type` 支持 `postgres` 和 `mysql`。
- metadata：继续写入 `metadata_snapshots.metadata_json`，`metadata.databaseType` 为 `mysql` 时表示 MySQL metadata。
- SQL 草稿：继续写入 `generated_sql_drafts`。

**测试**：

- 后端：Go unit test + handler test + metadata collector test + query executor test。
- 前端：Vitest + Vue Test Utils。
- 手动验证：后端测试通过后启动后端/前端，用 Google Chrome 完成 `interview_db` 基本流程。

**目标平台**：

- 本地 Web 应用：Go API 服务 + Vite/Vue 前端。
- MySQL 目标：本机或网络可达 MySQL 8.x / MySQL 兼容服务。

**项目类型**：

- Web application：`backend/` + `frontend/`。

**性能目标**：

- metadata 优先复用 SQLite 最新快照。
- MySQL 查询默认最多 1000 行。
- 查询与 metadata 采集必须有超时控制。
- LLM 上下文必须限制 metadata 大小，避免提示词过大。

**约束**：

- 前端不得接触完整数据库 URL、数据库密码或 LLM key。
- SQL 校验和执行必须在后端完成。
- MySQL 支持不能破坏现有 PostgreSQL 行为。
- API JSON 字段继续使用 camelCase。
- 当前仍不引入 authentication。
- 错误消息必须脱敏，不得泄露连接字符串、密码、key 或堆栈信息。

## 宪章检查

*GATE：Phase 0 research 前必须通过；Phase 1 design 后复查。*

- **只读数据库安全**：通过。MySQL 手写 SQL 和 LLM SQL 都必须经过后端 SQL Guard，只允许单条 SELECT，并默认 LIMIT。
- **Go 后端负责业务逻辑**：通过。连接、metadata、SQL 校验、查询执行、LLM 编排和错误脱敏都在 Go 后端。
- **Vue 3 + Element Plus 薄客户端**：通过。前端只增加数据库类型选择和状态展示，不实现数据库安全规则。
- **API 契约优先**：通过。API 变化先记录在 `contracts/openapi.yaml`，保持统一 envelope 和 camelCase。
- **严格类型**：通过。Go 和 TypeScript 都需要扩展 `databaseType` 为 `postgres | mysql`。
- **测试关键路径**：通过。计划覆盖 MySQL metadata、SQL Guard、query executor、API handler、NL2SQL 和 Chrome 验证。

当前无宪章违规项。

## 项目结构

### 文档结构

```text
specs/003-mysql-db-query/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
    └── openapi.yaml
```

### 源码结构

```text
backend/
├── cmd/server/main.go
├── internal/api/
│   ├── dbs_handlers.go
│   ├── interfaces.go
│   ├── query_handlers.go
│   ├── natural_query_handlers.go
│   └── types.go
├── internal/config/
├── internal/dbstore/
├── internal/llm/
├── internal/logging/
├── internal/metadata/
│   ├── collector.go
│   ├── mysql_collector.go
│   └── postgres_collector.go
├── internal/pgconn/
├── internal/mysqlconn/
├── internal/query/
│   ├── executor.go
│   ├── mysql_executor.go
│   └── postgres_executor.go
└── internal/sqlguard/
    ├── validator.go
    ├── mysql_validator.go
    └── postgres_validator.go

frontend/
├── src/api/
│   ├── client.ts
│   └── types.ts
├── src/components/
│   ├── ConnectionPanel.vue
│   ├── MetadataExplorer.vue
│   ├── QueryEditor.vue
│   ├── NaturalLanguagePanel.vue
│   └── ResultTable.vue
└── tests/unit/
```

**结构决策**：保持现有目录，不引入新的应用层框架。新增 MySQL 时优先通过接口分发和小模块扩展，避免把 PostgreSQL 和 MySQL 逻辑混在同一个文件中继续膨胀。

## API 计划

- `GET /api/v1/dbs`：响应中的 `databaseType` 支持 `postgres` 和 `mysql`。
- `PUT /api/v1/dbs/{name}`：请求体增加可选 `databaseType`，旧请求可由 URL 推断或默认 PostgreSQL；MySQL 推荐显式传 `databaseType: "mysql"`。
- `GET /api/v1/dbs/{name}`：metadata 中 `databaseType` 支持 `mysql`。
- `POST /api/v1/dbs/{name}/query`：不改路径，根据连接记录的 `databaseType` 执行对应 SQL Guard 和 query executor。
- `POST /api/v1/dbs/{name}/query/natural`：不改路径，根据 metadata.databaseType 生成对应方言 SQL。

## 日志定义

后端日志继续使用 `backend/internal/logging/logging.go` 的 stdout 单行格式：

```text
time="YYYY-MM-DD HH:mm:ss" level=info msg="..."
```

启动阶段必须至少输出：

```text
time="2026-07-22 10:00:00" level=info msg="llm config loaded: base_url=https://api2.codexcn.com/v1 model=gpt-5.5 wire_api=responses key_loaded=true"
time="2026-07-22 10:00:00" level=info msg="db-querry backend listening on :8080"
```

MySQL 相关实现新增日志时必须遵守：

- 输出到 stdout，方便 K8s 和云厂商采集。
- 使用固定时间格式 `2006-01-02 15:04:05`。
- 使用 `level=info|error` 等等级字段。
- 使用 `msg="..."` 包裹消息。
- 允许输出 `databaseType=mysql`、`dbName=interview_db`、`metadataStatus=ready`、`queryDurationMs=...`。
- 禁止输出完整 DB URL、密码、LLM key、原始 SQL 中可能包含的敏感字面量。

## Phase 0：Research 输出

见 [research.md](./research.md)。

## Phase 1：Design 输出

- 数据模型：[data-model.md](./data-model.md)
- API 契约：[contracts/openapi.yaml](./contracts/openapi.yaml)
- 验证指南：[quickstart.md](./quickstart.md)

## Complexity Tracking

当前无宪章违规项，无需复杂性豁免。
