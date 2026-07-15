# Quickstart：数据库查询工具

**日期**：2026-07-15

**关联计划**：[plan.md](./plan.md)

## 前置条件

- Go 1.23+。
- Node.js 20+。
- 一个可访问的 PostgreSQL 数据库。
- OpenAI API key，放在环境变量 `openai_api_key`。

## 环境变量

PowerShell：

```powershell
$env:openai_api_key = "sk-..."
```

后端必须从 `openai_api_key` 读取 OpenAI API key。前端不得读取或保存该 key。

## 本地数据文件

应用 SQLite 数据库固定为：

```text
~/.db_querry/db_querrt.db
```

注意：该文件会保存目标数据库连接 URL。它是敏感文件，不应提交到 Git，也不应暴露给前端。

## 启动后端

计划中的后端入口：

```powershell
cd backend
go run ./cmd/server
```

默认行为：

- 初始化 `~/.db_querry/db_querrt.db`。
- 启动 `/api/v1` API。
- 对所有 origin 启用 CORS。

## 启动前端

计划中的前端入口：

```powershell
cd frontend
npm install
npm run dev
```

前端使用 Vue 3 + Element Plus + TypeScript，并按 MotherDuck-inspired / brutalist dashboard 风格实现。

## API 验证流程

### 1. 添加数据库

```powershell
Invoke-RestMethod `
  -Method Put `
  -Uri "http://localhost:8080/api/v1/dbs/local" `
  -ContentType "application/json" `
  -Body '{"url":"postgres://postgres:postgre@localhost:5432/postgres"}'
```

预期：

- 返回 `success: true`。
- 返回连接名称和 metadata 状态。
- 不返回完整 DB URL。

### 2. 获取数据库列表

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/dbs"
```

预期：

- `data.dbs` 包含 `local`。
- `displayDsn` 是脱敏字符串。

### 3. 获取 metadata

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/dbs/local"
```

预期：

- 返回 `schemas`。
- 每个 table/view 包含 columns。

### 4. 执行 SQL

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/dbs/local/query" `
  -ContentType "application/json" `
  -Body '{"sql":"SELECT * FROM users"}'
```

预期：

- 如果 SQL 合法，返回 `columns`、`rows`、`rowCount`。
- 如果 SQL 没有 `LIMIT`，响应中 `limitApplied` 为 `true`，`limit` 为 `1000`。
- 如果表不存在，返回脱敏后的查询错误。

### 5. 验证禁止语句

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/dbs/local/query" `
  -ContentType "application/json" `
  -Body '{"sql":"DROP TABLE users"}'
```

预期：

- 返回 `success: false`。
- `error.code` 为 `sqlValidationFailed` 或同类安全错误。
- 数据库不会执行该语句。

### 6. 自然语言生成 SQL

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/dbs/local/query/natural" `
  -ContentType "application/json" `
  -Body '{"prompt":"查询用户表的所有信息"}'
```

预期：

- 返回 SQL 草稿。
- 返回校验状态。
- 不自动执行 SQL。

兼容性：

- 后端可以兼容读取错误拼写字段 `promt`。
- 正式 API、响应和前端类型统一使用 `prompt`。

## 前端验收路径

1. 打开前端页面。
2. 添加数据库连接。
3. 左侧连接列表出现新连接。
4. 右侧 metadata explorer 展示 schema/table/view/columns。
5. 在 SQL 编辑区输入 SELECT 查询并执行。
6. 下方结果表格展示 JSON 结果。
7. 在自然语言输入区输入问题，预览生成 SQL。
8. 点击执行生成 SQL，结果表格更新。

## 必须验证的安全行为

- 前端网络响应中看不到完整 DB URL。
- OpenAI key 不出现在任何前端代码或响应里。
- 非 SELECT SQL 被拒绝。
- 多语句 SQL 被拒绝。
- 无 LIMIT 的 SELECT 默认限制 1000 行。
- 错误响应不包含堆栈和密钥。
