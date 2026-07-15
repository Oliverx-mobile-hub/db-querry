# 实现计划：数据库查询工具

**分支**：`001-database-query-tool` | **日期**：2026-07-15 | **规格**：[spec.md](./spec.md)

**输入**：来自 `/specs/001-database-query-tool/spec.md` 的功能规格，以及本轮补充约束：

- 后端使用 Go 1.23+。
- 前端使用 Vue 3 + Element Plus + TypeScript。
- UI 风格参考会话 `019f64f6-a164-7fe1-8f9c-bced42126aef` 总结的 MotherDuck-inspired / brutalist dashboard 风格。
- OpenAI API key 来自环境变量 `openai_api_key`。
- 应用 SQLite 数据库路径为 `~/.db_querry/db_querrt.db`。
- 后端 API 支持 CORS，允许所有 origin。
- v1 API 路径使用 `/api/v1/dbs` 命名空间。

## 摘要

本功能实现一个本地运行的数据库查询工具：Go 后端负责连接 PostgreSQL、采集 table/view metadata、保存连接和 metadata 到 SQLite、校验并执行只读 SQL、调用 LLM 生成 SQL 草稿；Vue 3 前端负责连接管理、metadata 浏览、SQL 编辑、自然语言输入、结果表格和错误展示。

核心技术路线：

1. 后端启动时初始化 `~/.db_querry/db_querrt.db` 和必要 SQLite 表。
2. 用户通过 `PUT /api/v1/dbs/{name}` 添加或更新数据库连接。
3. 后端连接 PostgreSQL，读取系统 catalog / information_schema，生成结构化 metadata JSON 并保存。
4. 用户查询 metadata 时通过 `GET /api/v1/dbs/{name}` 获取已保存的 metadata。
5. 用户提交 SQL 时，后端先用 PostgreSQL 方言 parser 和只读策略校验，再执行；缺少 `LIMIT` 时默认限制 1000 行。
6. 用户提交自然语言时，后端把 metadata 作为上下文传给 LLM，生成 SQL 草稿，再通过同一 SQL 校验链路返回给前端。

## 技术上下文

**语言 / 版本**：

- 后端：Go 1.23+。
- 前端：TypeScript + Vue 3。

**主要依赖**：

- 后端 HTTP：Go 标准库 `net/http` 或轻量 router；优先保持简单，除非实现中确实需要更强路由能力。
- 后端 SQLite：Go SQLite driver。
- 后端 PostgreSQL：Go PostgreSQL driver。
- SQL 解析：PostgreSQL 方言 parser，用于解析 AST 并识别语句类型、多语句和禁止语法。
- LLM：OpenAI API，API key 从 `openai_api_key` 环境变量读取。
- 前端：Vue 3、Element Plus、Vite、TypeScript。

**存储**：

- 应用本地 SQLite：`~/.db_querry/db_querrt.db`。
- 目标业务数据库：v1 支持 PostgreSQL / PostgreSQL 兼容连接。
- 连接 URL 和 metadata 都保存在 SQLite；前端永远不接触完整连接 URL。

**测试**：

- 后端：Go unit test + handler test + SQLite 临时库集成测试。
- 前端：Vitest + Vue Test Utils，覆盖核心组件和 API 状态流。
- 关键路径优先：SQL 只读校验、禁止语句检测、默认 LIMIT、错误脱敏、metadata 存储、API envelope。

**目标平台**：

- 本地 Web 应用：Go API 服务 + Vite/Vue 前端。
- 默认运行在开发者本机，可连接本机或网络可达 PostgreSQL。

**项目类型**：

- Web application：`backend/` + `frontend/`。

**性能目标**：

- metadata 查询优先从 SQLite 返回，避免每次页面加载都连接目标数据库。
- SQL 查询默认最多 1000 行，避免无界结果集。
- 单次查询必须有超时控制，具体默认值在实现任务中固定。
- LLM 上下文必须限制 metadata 大小，避免超出模型上下文或造成高延迟。

**约束**：

- 所有数据库访问、SQL 校验、SQL 生成编排和查询执行必须在 Go 后端完成。
- 所有 API 响应 JSON 字段使用 camelCase。
- 后端 CORS 允许所有 origin。
- 当前不引入 authentication。
- 前端不得保存数据库凭据、OpenAI key 或高权限配置。
- 用户提交 SQL 和 LLM 生成 SQL 都视为不可信输入。

**规模 / 范围**：

- v1 支持多个命名数据库连接。
- v1 metadata 范围：schema、table、view、column、data type、nullable、primary key、comment。
- v1 查询结果范围：只读 SELECT，默认最多 1000 行。

## 宪章检查

*GATE：Phase 0 research 前必须通过；Phase 1 design 后复查。*

- **只读数据库安全**：通过。计划要求所有手写 SQL 和 LLM SQL 都经过服务端 parser + policy 校验；禁止非 SELECT、多语句、修改性语句和执行语句；缺省 LIMIT 1000。
- **Go 后端负责业务逻辑**：通过。连接、metadata、SQL 校验、查询执行、LLM 编排、错误脱敏都在 Go 后端。
- **Vue 3 + Element Plus 薄客户端**：通过。前端只负责交互、展示和调用 API，不复制后端安全规则。
- **API 契约优先**：通过。plan 固定 v1 endpoint、统一响应 envelope 和错误模型；详细字段在 `contracts/` 中定义。
- **严格类型**：通过。Go 使用 request/response struct、service model、error type；Vue 使用 TypeScript API 类型和组件状态类型。
- **测试关键路径**：通过。后端 SQL 安全、错误脱敏和执行边界列为必须测试项。

当前无宪章违规项。

## 项目结构

### 文档结构

```text
specs/001-database-query-tool/
├── spec.md
├── plan.md
├── research.md          # 后续 Phase 0 输出
├── data-model.md        # 后续 Phase 1 输出
├── quickstart.md        # 后续 Phase 1 输出
├── contracts/           # 后续 OpenAPI / API examples
└── tasks.md             # 后续 Phase 2 输出
```

### 源码结构

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   ├── routes.go
│   │   └── responses.go
│   ├── config/
│   │   └── config.go
│   ├── dbstore/
│   │   ├── migrations.go
│   │   └── store.go
│   ├── llm/
│   │   └── openai.go
│   ├── metadata/
│   │   ├── collector.go
│   │   └── model.go
│   ├── pgconn/
│   │   └── connector.go
│   ├── query/
│   │   ├── executor.go
│   │   └── result.go
│   └── sqlguard/
│       ├── validator.go
│       └── validator_test.go
└── tests/
    ├── api/
    └── integration/

frontend/
├── src/
│   ├── api/
│   │   ├── client.ts
│   │   └── types.ts
│   ├── components/
│   │   ├── ConnectionPanel.vue
│   │   ├── MetadataExplorer.vue
│   │   ├── QueryEditor.vue
│   │   ├── NaturalLanguagePanel.vue
│   │   └── ResultTable.vue
│   ├── styles/
│   │   ├── design-tokens.css
│   │   └── global.css
│   ├── App.vue
│   └── main.ts
└── tests/
    └── unit/
```

**结构决策**：采用现有 `backend/` 和 `frontend/` 双目录。后端按职责拆分为 API、配置、SQLite store、PostgreSQL connector、metadata、SQL guard、query executor、LLM client。前端按 dashboard 工作区组件拆分，不引入复杂全局状态库，先用 Vue 组合式状态和局部组件状态实现。

## API 计划

### 统一响应结构

所有 API 返回统一 envelope：

```json
{
  "success": true,
  "data": {},
  "error": null
}
```

错误响应：

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "sqlValidationFailed",
    "message": "SQL 只能包含一条 SELECT 查询",
    "details": {}
  }
}
```

### Endpoint

#### 获取所有已连接数据库

`GET /api/v1/dbs`

返回已保存数据库列表，不包含完整 DB URL 和密码。

响应 data：

```json
{
  "dbs": [
    {
      "name": "local",
      "databaseType": "postgres",
      "displayDsn": "postgres://postgres@localhost:5432/postgres",
      "metadataStatus": "ready",
      "metadataUpdatedAt": "2026-07-15T16:00:00+08:00"
    }
  ]
}
```

#### 添加或更新数据库

`PUT /api/v1/dbs/{name}`

请求：

```json
{
  "url": "postgres://postgres:postgre@localhost:5432/postgres"
}
```

行为：

- 校验 `{name}` 合法性，作为连接的稳定标识。
- 在后端测试连接。
- 保存或更新连接 URL。
- 采集 metadata 并保存到 SQLite。
- 返回连接概要和 metadata 状态。

#### 获取一个数据库的 metadata

`GET /api/v1/dbs/{name}`

返回指定数据库连接的 metadata 快照。

响应 data：

```json
{
  "name": "local",
  "metadataStatus": "ready",
  "metadataUpdatedAt": "2026-07-15T16:00:00+08:00",
  "schemas": []
}
```

#### 查询某个数据库

`POST /api/v1/dbs/{name}/query`

请求：

```json
{
  "sql": "SELECT * FROM users"
}
```

行为：

- 解析 SQL。
- 拒绝多语句和非 SELECT。
- 拒绝禁止关键字和可能副作用语法。
- 无 `LIMIT` 时应用 `LIMIT 1000`。
- 执行查询并返回 JSON 结果。

#### 根据自然语言生成 SQL

`POST /api/v1/dbs/{name}/query/natural`

请求：

```json
{
  "prompt": "查询用户表的所有信息"
}
```

兼容性：

- 用户草案中字段写作 `promt`。实现时请求 DTO 应优先使用正确字段 `prompt`，同时可兼容读取 `promt`，避免早期调用方失败。
- 响应和前端类型统一使用 `prompt`。

行为：

- 从 SQLite 读取该数据库 metadata。
- 构造 LLM 上下文。
- 调用 OpenAI API。
- 将 LLM 输出解析为 SQL 草稿。
- 对生成 SQL 执行同一 SQL 校验。
- 返回 SQL 草稿、说明、引用对象和校验状态，不自动执行。

## 数据交互模型

### 总体数据流

```mermaid
flowchart LR
    UI["Vue 3 + Element Plus 前端"] --> API["Go API /api/v1"]
    API --> Store["SQLite ~/.db_querry/db_querrt.db"]
    API --> PG["目标 PostgreSQL"]
    API --> Guard["SQL Guard"]
    API --> LLM["OpenAI API"]
    PG --> API
    Store --> API
    Guard --> API
    LLM --> API
    API --> UI
```

### 添加数据库

1. 前端提交 `PUT /api/v1/dbs/{name}`。
2. 后端解析 URL，仅在服务端保存完整 URL。
3. 后端连接 PostgreSQL 并执行轻量连接测试。
4. 后端读取系统 catalog / information_schema。
5. 后端将原始 metadata 归一化为严格 Go 类型。
6. 如需 LLM 辅助整理说明文本，LLM 输出必须再解析回严格类型。
7. 后端保存连接和 metadata 快照到 SQLite。
8. 前端收到连接状态并刷新左侧数据库列表和 metadata 面板。

### SQL 查询

1. 前端提交 SQL 到 `POST /api/v1/dbs/{name}/query`。
2. 后端读取连接 URL。
3. SQL Guard 解析并校验 SQL。
4. 如果没有 LIMIT，生成受限 SQL 或执行时强制限制 1000 行。
5. Query Executor 使用只读上下文和超时执行查询。
6. 后端将数据库结果转换为 JSON 友好类型。
7. 前端用 Element Plus table 展示结果。

### 自然语言生成 SQL

1. 前端提交自然语言 `prompt`。
2. 后端读取该连接 metadata。
3. 后端裁剪 metadata 上下文，避免过大。
4. 后端调用 OpenAI API，key 从 `openai_api_key` 读取。
5. 后端解析 LLM 输出，只接受结构化 SQL 草稿。
6. SQL Guard 校验生成 SQL。
7. 前端展示 SQL 草稿和校验状态，由用户决定是否执行。

## SQLite 数据模型初稿

详细字段和约束在 `data-model.md` 固化。plan 阶段采用以下方向：

### `db_connections`

- `name`：主键，来自 `/api/v1/dbs/{name}`。
- `database_type`：v1 固定为 `postgres`。
- `url`：完整连接 URL，仅后端读取。
- `display_dsn`：脱敏后的展示字符串。
- `metadata_status`：`pending` / `ready` / `failed`。
- `created_at`、`updated_at`。

### `metadata_snapshots`

- `id`：主键。
- `db_name`：关联 `db_connections.name`。
- `metadata_json`：结构化 metadata JSON。
- `object_count`。
- `warning_json`：脱敏采集警告。
- `created_at`。

### `generated_sql_drafts`

- `id`：主键。
- `db_name`。
- `prompt`。
- `sql`。
- `validation_json`。
- `created_at`。

v1 可以不保存查询历史；如果实现中需要查询历史，必须先更新 data model 和 tasks。

## Metadata JSON 模型初稿

```json
{
  "schemas": [
    {
      "name": "public",
      "objects": [
        {
          "schema": "public",
          "name": "users",
          "type": "table",
          "comment": "用户表",
          "columns": [
            {
              "name": "id",
              "dataType": "integer",
              "nullable": false,
              "primaryKey": true,
              "ordinal": 1,
              "comment": "用户 ID"
            }
          ]
        }
      ]
    }
  ]
}
```

## SQL 安全策略

- 只允许单条 SELECT 查询。
- 禁止多语句。
- 禁止 DML / DDL / DCL / procedure / execute 类语句。
- 禁止 `SELECT ... INTO`。
- 对 CTE、子查询、窗口函数保持支持，只要 parser 和策略判断为只读。
- 函数调用默认允许常见只读函数；高风险函数策略需在 research 阶段确认。
- 无 LIMIT 时默认应用 `LIMIT 1000`。
- 所有 SQL 错误返回脱敏后的用户消息。

## LLM 方案

- OpenAI API key 读取环境变量 `openai_api_key`。
- 后端启动时不把 key 返回给任何 API。
- Prompt 输入包括：用户自然语言、数据库 metadata 摘要、只读 SQL 要求、输出格式要求。
- LLM 输出必须是结构化 JSON，至少包含 `sql` 和 `explanation`。
- 后端必须将 LLM 输出作为不可信输入处理，解析失败或 SQL 校验失败时返回不可执行状态。
- 自然语言接口只生成 SQL 草稿，不执行查询。

## 前端设计计划

### 视觉风格

参考会话 `019f64f6-a164-7fe1-8f9c-bced42126aef` 中总结的 MotherDuck-inspired / brutalist dashboard：

- 页面背景：奶油色 `#F4EFEA`。
- 面板背景：白色 `#FFFFFF`。
- 主文字：深灰 `#383838`。
- 强边框：黑色 `2px solid #000000`。
- 圆角：2px 左右的小圆角。
- 主按钮：蓝色底、黑色边框、深色文字。
- 强调色：黄色 `#FFDE00`。
- 表头、标签和小标题使用粗体 uppercase 风格。
- SQL 编辑器区域可使用深色主题，作为页面中唯一明显深色区域。
- 避免柔和卡片阴影，必要时使用硬投影。

### 页面结构

- 左侧连接列表：显示数据库连接、状态、添加入口。
- 中间主工作区：SQL 编辑器、自然语言输入、生成 SQL 预览、执行按钮。
- 右侧 metadata explorer：schema、table/view、columns。
- 底部或主区下半部分：查询结果表格、空状态、错误状态。

### Element Plus 使用

- `el-button`：主操作、刷新 metadata、执行 SQL。
- `el-input` / textarea：DB URL、自然语言 prompt。
- `el-table`：查询结果。
- `el-tree` 或 `el-collapse`：metadata explorer。
- `el-dialog`：添加数据库连接。
- `el-alert`：错误和校验提示。
- `el-tag`：metadata 状态、SQL 校验状态。

## 错误模型

错误码采用 camelCase：

- `invalidRequest`
- `dbNotFound`
- `dbConnectionFailed`
- `metadataCollectionFailed`
- `sqlParseFailed`
- `sqlValidationFailed`
- `queryExecutionFailed`
- `llmUnavailable`
- `llmOutputInvalid`
- `internalError`

错误消息面向用户，不能包含：

- 完整 DB URL。
- 密码。
- OpenAI key。
- 堆栈。
- 原始驱动错误中的敏感连接详情。

## Phase 0：Research

需要在 `research.md` 中确认：

- PostgreSQL metadata 查询 SQL：选择 catalog / information_schema 的组合。
- SQL parser 选择：是否使用 PostgreSQL AST parser，以及对复杂 SELECT 的兼容边界。
- SQLite driver 选择和本地文件路径创建策略。
- OpenAI 调用方式和结构化输出策略。
- 高风险 SQL 函数策略。
- DB URL 是否直接保存在 SQLite，或是否需要本地加密；如果选择直接保存，必须记录风险和后续改进路径。

## Phase 1：Design

输出：

- `data-model.md`：SQLite 表结构、Go model、metadata JSON schema、query result schema。
- `contracts/`：OpenAPI 或 endpoint markdown examples，覆盖所有 `/api/v1/dbs` endpoint。
- `quickstart.md`：本地启动步骤、环境变量、示例 PostgreSQL 连接、示例查询。
- 前端设计 token 初稿：`frontend/src/styles/design-tokens.css`。

## Phase 2：Tasks

后续 `tasks.md` 应按可独立交付的顺序拆分：

1. 后端项目初始化和配置。
2. SQLite store 和 migration。
3. 数据库连接管理 API。
4. PostgreSQL metadata collector。
5. SQL Guard。
6. Query Executor。
7. Natural language SQL generation。
8. 前端项目初始化和 design tokens。
9. 连接管理 UI。
10. Metadata explorer。
11. SQL editor + result table。
12. Natural language SQL panel。
13. 测试和 quickstart 验证。

## 复杂度跟踪

当前无宪章违规。

| 复杂点 | 为什么需要 | 简化边界 |
| --- | --- | --- |
| SQL parser + policy guard | 仅靠字符串匹配无法可靠判断只读 SQL、多语句和语法正确性 | v1 只支持 PostgreSQL 方言 SELECT |
| SQLite metadata cache | metadata 需要复用，LLM 生成 SQL 不能每次重新连接目标数据库采集 | v1 只保存最新 snapshot，不做复杂版本比较 |
| LLM 生成 SQL | 产品核心能力之一 | 只生成草稿，不自动执行，必须通过 SQL Guard |
