# Data Model：数据库查询工具

**日期**：2026-07-15

**关联计划**：[plan.md](./plan.md)

## SQLite 数据库

SQLite 文件路径固定为：

```text
~/.db_querry/db_querrt.db
```

后端启动时必须创建目录和数据库文件，并执行迁移。SQLite 文件包含目标数据库连接 URL，属于敏感文件，不得由前端读取。

## 表：`db_connections`

保存用户添加的数据库连接。

```sql
CREATE TABLE IF NOT EXISTS db_connections (
  name TEXT PRIMARY KEY,
  database_type TEXT NOT NULL,
  url TEXT NOT NULL,
  display_dsn TEXT NOT NULL,
  metadata_status TEXT NOT NULL,
  metadata_error TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
```

字段说明：

- `name`：连接名称，来自 `/api/v1/dbs/{name}`，作为稳定标识。
- `database_type`：v1 固定为 `postgres`。
- `url`：完整 DB URL，仅后端读取。
- `display_dsn`：脱敏展示字符串，例如 `postgres://postgres@localhost:5432/postgres`。
- `metadata_status`：`pending`、`ready`、`failed`。
- `metadata_error`：脱敏后的最近一次 metadata 错误。
- `created_at` / `updated_at`：RFC3339 时间字符串。

约束：

- API 响应不得返回 `url`。
- `name` 只允许安全字符，建议 `^[A-Za-z0-9_-]{1,64}$`。

## 表：`metadata_snapshots`

保存每次采集得到的 metadata 快照。v1 查询时默认使用最新一条。

```sql
CREATE TABLE IF NOT EXISTS metadata_snapshots (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  db_name TEXT NOT NULL,
  metadata_json TEXT NOT NULL,
  object_count INTEGER NOT NULL,
  warning_json TEXT NOT NULL DEFAULT '[]',
  created_at TEXT NOT NULL,
  FOREIGN KEY (db_name) REFERENCES db_connections(name) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_metadata_snapshots_db_name_created_at
  ON metadata_snapshots(db_name, created_at DESC);
```

字段说明：

- `db_name`：关联 `db_connections.name`。
- `metadata_json`：结构化 metadata JSON。
- `object_count`：table/view 数量。
- `warning_json`：脱敏采集警告数组。
- `created_at`：采集时间。

## 表：`generated_sql_drafts`

保存自然语言生成 SQL 的草稿记录，便于调试和后续复用。

```sql
CREATE TABLE IF NOT EXISTS generated_sql_drafts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  db_name TEXT NOT NULL,
  prompt TEXT NOT NULL,
  sql TEXT NOT NULL,
  explanation TEXT NOT NULL DEFAULT '',
  referenced_objects_json TEXT NOT NULL DEFAULT '[]',
  validation_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (db_name) REFERENCES db_connections(name) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_generated_sql_drafts_db_name_created_at
  ON generated_sql_drafts(db_name, created_at DESC);
```

字段说明：

- `prompt`：用户自然语言输入，正式字段名使用 `prompt`。
- `sql`：LLM 生成的 SQL 草稿。
- `explanation`：简短说明。
- `referenced_objects_json`：LLM 声称引用的 table/view。
- `validation_json`：SQL Guard 对草稿的校验结果。

## Metadata JSON Schema

`metadata_json` 顶层结构：

```json
{
  "databaseType": "postgres",
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

类型约束：

- `databaseType`：v1 固定 `postgres`。
- `schemas[].name`：schema 名。
- `objects[].type`：`table` 或 `view`。
- `objects[].columns[]`：按 `ordinal` 升序。
- `comment` 可为空字符串，不使用 `null`。

## API Model

### Response Envelope

```ts
type ApiSuccess<T> = {
  success: true
  data: T
  error: null
}

type ApiFailure = {
  success: false
  data: null
  error: ApiError
}

type ApiResponse<T> = ApiSuccess<T> | ApiFailure
```

### ApiError

```ts
type ApiError = {
  code:
    | 'invalidRequest'
    | 'dbNotFound'
    | 'dbConnectionFailed'
    | 'metadataCollectionFailed'
    | 'sqlParseFailed'
    | 'sqlValidationFailed'
    | 'queryExecutionFailed'
    | 'llmUnavailable'
    | 'llmOutputInvalid'
    | 'internalError'
  message: string
  details: Record<string, unknown>
}
```

### DbSummary

```ts
type DbSummary = {
  name: string
  databaseType: 'postgres'
  displayDsn: string
  metadataStatus: 'pending' | 'ready' | 'failed'
  metadataUpdatedAt: string | null
}
```

### SqlValidationResult

```ts
type SqlValidationResult = {
  valid: boolean
  executable: boolean
  statementType: 'select' | 'unknown'
  normalizedSql: string
  limitApplied: boolean
  limit: number | null
  errors: Array<{
    code: string
    message: string
  }>
}
```

### QueryExecutionResult

```ts
type QueryExecutionResult = {
  columns: Array<{
    name: string
    dataType: string
  }>
  rows: Array<Record<string, unknown>>
  rowCount: number
  durationMs: number
  limitApplied: boolean
  limit: number | null
  empty: boolean
}
```

### GeneratedSqlDraft

```ts
type GeneratedSqlDraft = {
  prompt: string
  sql: string
  explanation: string
  referencedObjects: string[]
  validation: SqlValidationResult
}
```

## Go 类型映射

Go 后端应使用明确 struct 表达 API request/response、metadata、query result 和错误。

建议：

- 时间在 API 中序列化为 RFC3339 字符串。
- 查询结果行可以先从 driver 值收窄到 JSON 友好类型：`string`、`float64`、`int64`、`bool`、`nil`、RFC3339 时间字符串。
- 不在核心业务边界长期传递 `map[string]any`；如果用于动态查询行，应在 query executor 边界封装。

## 生命周期

### 添加或更新连接

1. upsert `db_connections`，状态设为 `pending`。
2. 连接 PostgreSQL 并采集 metadata。
3. 插入 `metadata_snapshots`。
4. 更新 `db_connections.metadata_status` 为 `ready`。
5. 如果失败，更新为 `failed` 并保存脱敏错误。

### 查询 metadata

1. 查询 `db_connections`。
2. 查询最新 `metadata_snapshots`。
3. 返回 connection summary + metadata JSON。

### 自然语言生成 SQL

1. 查询最新 metadata。
2. 调用 LLM。
3. SQL Guard 校验。
4. 插入 `generated_sql_drafts`。
5. 返回草稿和校验结果。
