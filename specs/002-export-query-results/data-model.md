# Data Model：导出当前查询结果为 CSV/JSON

**日期**：2026-07-20

**关联计划**：[plan.md](./plan.md)

## 概览

本功能不新增后端数据库表，不新增 SQLite 字段，也不新增持久化实体。

导出功能围绕前端当前内存中的 `QueryResult` 派生出 CSV 文本、JSON 文本和下载文件名。

## 现有实体：QueryResult

前端已有类型：

```ts
export type QueryResult = {
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

### 字段说明

- `columns`：列定义。CSV 表头和列顺序必须以此字段为准。
- `rows`：当前查询结果行。CSV 和 JSON 的主要数据来源。
- `rowCount`：当前返回行数。JSON 导出保留。
- `durationMs`：查询耗时。JSON 导出保留。
- `limitApplied`：后端是否自动应用 LIMIT。JSON 导出保留。
- `limit`：LIMIT 值。JSON 导出保留。
- `empty`：是否为空结果。用于禁用导出按钮。
- `validation`：可选 SQL 校验结果。JSON 导出保留。

## 新增前端模型：ExportFormat

```ts
export type ExportFormat = 'csv' | 'json'
```

### 规则

- `csv`：导出表格数据。
- `json`：导出完整 `QueryResult`。

## 新增前端模型：ExportContext

```ts
export type ExportContext = {
  dbName: string | null
  result: QueryResult | null
  loading: boolean
  now: Date
}
```

### 字段说明

- `dbName`：当前选中的数据库名；文件名使用该值，空值时使用 `query`。
- `result`：当前查询结果；为空时不可导出。
- `loading`：当前是否正在查询；为 true 时不可导出。
- `now`：生成文件名所需时间；测试可注入固定时间。

## 派生模型：ExportFile

```ts
export type ExportFile = {
  filename: string
  mimeType: string
  content: string
}
```

### 规则

- CSV：
  - `filename`: `{safeDbName}-query-{YYYYMMDD-HHmmss}.csv`
  - `mimeType`: `text/csv;charset=utf-8`
  - `content`: CSV 文本，第一行为表头。
- JSON：
  - `filename`: `{safeDbName}-query-{YYYYMMDD-HHmmss}.json`
  - `mimeType`: `application/json;charset=utf-8`
  - `content`: `JSON.stringify(result, null, 2)`

## CSV 序列化规则

- 表头顺序等于 `QueryResult.columns` 顺序。
- 每行按 `column.name` 读取对应值。
- `null` 和 `undefined` 输出为空字段。
- 数字和布尔值使用 `String(value)`。
- 字符串原样作为字段值参与转义。
- 日期或时间字符串按当前 JSON 值输出，不额外转换时区。
- 字段包含逗号、双引号、CR、LF 时，必须使用双引号包裹。
- 字段内双引号替换为两个双引号。
- 行分隔符使用 `\r\n`，便于 Excel / Windows 表格工具识别。

## 文件名规则

```text
{safeDbName}-query-{YYYYMMDD-HHmmss}.{csv|json}
```

- `safeDbName`：将数据库名中非 `[A-Za-z0-9_-]` 的字符替换为 `-`。
- 如果数据库名为空，使用 `query`。
- 时间戳使用浏览器本地时间。

## 状态规则

导出按钮启用条件：

```text
result != null AND result.empty == false AND loading == false
```

其他状态均禁用。
