# Data Model：MySQL 数据库查询支持

**日期**：2026-07-22

**关联计划**：[plan.md](./plan.md)

## 概览

本功能复用现有 SQLite 表结构，并把已有 `database_type` 字段从当前实际单值 `postgres` 扩展为 `postgres | mysql`。

## 数据库类型：DatabaseType

```ts
type DatabaseType = 'postgres' | 'mysql'
```

Go 中建议使用明确类型：

```go
type DatabaseType string

const (
  DatabaseTypePostgres DatabaseType = "postgres"
  DatabaseTypeMySQL    DatabaseType = "mysql"
)
```

## 实体：db_connections

现有 SQLite 表：

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

### 新规则

- `database_type` 支持 `postgres` 和 `mysql`。
- `name` 示例：`interview_db`。
- `url` 保存完整连接字符串，只能后端读取。
- `display_dsn` 必须脱敏，不能包含密码。
- API 响应只能返回 `displayDsn`，不能返回完整 `url`。

## 实体：MetadataDocument

```ts
type MetadataDocument = {
  databaseType: 'postgres' | 'mysql'
  schemas: MetadataSchema[]
}
```

### MySQL 映射规则

- MySQL database 映射为 `MetadataSchema.name`。
- `information_schema.tables.table_type = 'BASE TABLE'` 映射为 `type: 'table'`。
- `information_schema.tables.table_type = 'VIEW'` 映射为 `type: 'view'`。
- `information_schema.columns.column_name` 映射为 `MetadataColumn.name`。
- `column_type` 或 `data_type` 映射为 `dataType`。
- `is_nullable = 'YES'` 映射为 `nullable=true`。
- primary key 从 `table_constraints` + `key_column_usage` 判断。
- `table_comment` 和 `column_comment` 映射为 comment。

## 实体：QueryResult

结构保持不变：

```ts
type QueryResult = {
  columns: Array<{ name: string; dataType: string }>
  rows: Array<Record<string, unknown>>
  rowCount: number
  durationMs: number
  limitApplied: boolean
  limit: number | null
  empty: boolean
  validation?: SqlValidationResult
}
```

### MySQL 值转换规则

- `NULL` -> `null`
- `[]byte` 文本 -> `string`
- `time.Time` -> RFC3339 string
- `decimal` -> string 或 number，必须保持 JSON 安全
- `json` 字段 -> 如果 driver 返回 `[]byte`，先以 string 返回，后续可考虑解析
- `blob` -> 不直接暴露原始二进制；v1 可返回 base64 或脱敏提示，具体在实现阶段确认

## 实体：GeneratedSQLDraft

结构保持不变：

```ts
type GeneratedSqlDraft = {
  prompt: string
  sql: string
  explanation: string
  referencedObjects: string[]
  validation: SqlValidationResult
}
```

### MySQL 规则

- LLM 生成 SQL 时必须使用 MySQL 方言。
- 生成结果必须经过 MySQL SQL Guard。
- `referencedObjects` 建议使用 `database.table` 或 `table`，由 metadata 可定位即可。

## API 请求模型：PutDbRequest

```ts
type PutDbRequest = {
  url: string
  databaseType?: 'postgres' | 'mysql'
}
```

### 兼容规则

- 旧请求只传 `url` 时保持兼容。
- `databaseType=mysql` 时后端按 MySQL 连接、metadata 和查询路径处理。
- 如果 `databaseType` 与 URL scheme 明显冲突，后端应返回 `invalidRequest`。

## 状态规则

- `metadataStatus`: `pending | ready | failed`
- `connectionStatus`: `online | offline | unknown`

MySQL 离线时，前端不应展示旧 metadata 伪装为 ready。
