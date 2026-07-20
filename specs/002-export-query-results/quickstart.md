# Quickstart：导出当前查询结果为 CSV/JSON

**日期**：2026-07-20

**关联计划**：[plan.md](./plan.md)

## 前置条件

- 后端已经可以启动并连接测试 PostgreSQL。
- 前端已经可以运行并展示查询结果。
- `backend/env/.env` 中大模型配置不影响本功能，但保持现有配置即可。

## 启动服务

后端：

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\backend
C:\Users\Oliver-x\sdk\go1.26.2\bin\go.exe run ./cmd/server
```

前端：

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\frontend
npm.cmd run dev
```

## 手动验证

1. 打开前端页面。
2. 选择一个已连接且 online 的数据库。
3. 执行包含特殊字符的查询，例如：

```sql
SELECT
  1 AS id,
  'Alice, Bob' AS comma_text,
  'He said "Hello"' AS quote_text,
  E'line1\nline2' AS multiline_text,
  NULL AS null_value,
  42 AS amount,
  now()::text AS created_at;
```

4. 查询成功后，确认 Results 面板显示 `Export CSV` 和 `Export JSON` 两个按钮，并且按钮可点击。
5. 点击 `Export CSV`：
   - 文件名应类似 `local-query-20260720-143000.csv`。
   - CSV 第一行是表头。
   - 逗号、双引号、换行和 null 不应导致串列。
6. 点击 `Export JSON`：
   - 文件名应类似 `local-query-20260720-143000.json`。
   - 文件可以被 JSON parser 解析。
   - JSON 包含完整 `QueryResult`，包括 `columns`、`rows`、`rowCount`、`durationMs`、`limitApplied`、`limit`、`empty`、`validation`。

## 禁用状态验证

1. 刷新页面，未执行查询时，两个导出按钮应禁用。
2. 执行查询过程中，两个导出按钮应禁用。
3. 执行返回空结果的查询：

```sql
SELECT * FROM users WHERE 1 = 0;
```

结果为空时，两个导出按钮应禁用。

## 自动化测试建议

前端：

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\frontend
npm.cmd test
```

应覆盖：

- CSV 字段转义：逗号、双引号、换行、null、数字、布尔值。
- 文件名格式：数据库名安全化和时间戳。
- JSON 导出：完整 `QueryResult`。
- `ResultTable` 按钮状态：无结果、空结果、查询中、非空结果。

后端：

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\backend
C:\Users\Oliver-x\sdk\go1.26.2\bin\go.exe test ./...
```

如实现阶段调整后端启动日志，应确认：

- 日志输出到 stdout。
- 启动日志包含非敏感配置摘要，例如 `base_url`、`model`、`wire_api`、`key_loaded`。
- 日志不包含 API key、数据库密码或完整 DB URL。
