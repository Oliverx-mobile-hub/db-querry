# Quickstart：MySQL 数据库查询支持

**日期**：2026-07-22

**关联计划**：[plan.md](./plan.md)

## 前置条件

- 本地或网络可达 MySQL 8.x。
- 已创建测试数据库 `interview_db`。
- 已准备至少一张测试表，例如 `candidates`。
- `backend/env/.env` 中 LLM 配置可用，用于自然语言生成 SQL。

## 建议测试数据

```sql
CREATE DATABASE IF NOT EXISTS interview_db;
USE interview_db;

CREATE TABLE IF NOT EXISTS candidates (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(120) NOT NULL,
  role VARCHAR(120) NOT NULL,
  years_experience INT NOT NULL,
  city VARCHAR(80),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO candidates (name, role, years_experience, city)
VALUES
  ('Alice Zhang', 'Backend Engineer', 5, 'Shanghai'),
  ('Bob Li', 'Frontend Engineer', 3, 'Hangzhou');
```

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

## 后端测试

实现完成后先运行：

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\backend
C:\Users\Oliver-x\sdk\go1.26.2\bin\go.exe test ./...
```

后端测试通过后再进行 Chrome 前端验证。

## API 快速验证

添加 MySQL 连接：

```powershell
$body = @{
  databaseType = "mysql"
  url = "mysql://root:password@localhost:3306/interview_db"
} | ConvertTo-Json

Invoke-RestMethod -Method Put `
  -Uri "http://localhost:8080/api/v1/dbs/interview_db" `
  -ContentType "application/json" `
  -Body $body
```

执行手写 SQL：

```powershell
$body = @{ sql = "SELECT id, name, role, years_experience FROM candidates" } | ConvertTo-Json

Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/api/v1/dbs/interview_db/query" `
  -ContentType "application/json" `
  -Body $body
```

自然语言生成 SQL：

```powershell
$body = @{ prompt = "查询所有候选人的姓名、岗位和工作年限" } | ConvertTo-Json

Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/api/v1/dbs/interview_db/query/natural" `
  -ContentType "application/json" `
  -Body $body
```

## Google Chrome 前端验证

后端测试通过后：

1. 打开 Google Chrome。
2. 访问前端地址，例如 `http://localhost:5173`。
3. 点击 `Add`。
4. 选择数据库类型 `MySQL`。
5. 名称填写 `interview_db`。
6. URL 填写 MySQL 连接字符串。
7. 保存后确认数据库列表显示 `interview_db`，状态为 online，metadata 为 ready。
8. 点击 `interview_db`，确认左侧或旁侧 metadata 面板能看到 `candidates` 表。
9. 在 SQL Editor 输入：

```sql
SELECT id, name, role, years_experience FROM candidates;
```

10. 点击运行，确认 Results 表格显示数据。
11. 在 Natural SQL 输入：

```text
查询所有候选人的姓名、岗位和工作年限
```

12. 生成 SQL 后执行，确认结果显示。

## 日志检查

后端启动时应看到：

```text
time="2026-07-22 10:00:00" level=info msg="llm config loaded: base_url=https://api2.codexcn.com/v1 model=gpt-5.5 wire_api=responses key_loaded=true"
time="2026-07-22 10:00:00" level=info msg="db-querry backend listening on :8080"
```

MySQL 连接和查询阶段的新增日志如果实现，必须脱敏：

```text
time="2026-07-22 10:00:05" level=info msg="db connection saved: dbName=interview_db databaseType=mysql"
time="2026-07-22 10:00:06" level=info msg="metadata collected: dbName=interview_db databaseType=mysql objectCount=1"
```

不得出现密码、完整 DB URL、LLM key。
