# interview_db 本地 MySQL 测试库

Docker 只负责启动 MySQL 8.4。启动后数据库是空的，表结构和种子数据由你手动导入：

- `init/01-schema.sql`：20 张业务表、索引、外键、约束和 4 个分析视图。
- `init/02-seed.sql`：招聘需求、候选人、申请、面试、评分和 Offer 等合成数据。

所有公司名、姓名、电话、邮箱和履历均为虚构数据，仅供本地测试。

## 第一步：启动空 MySQL

前提：Docker Desktop 已启动。

```powershell
cd C:\Users\Oliver-x\Desktop\db-querry\docker\interview-db
Copy-Item .env.example .env
notepad .env
```

把 `.env` 中的两个示例密码替换掉，然后执行：

```powershell
docker compose config
docker compose pull
docker compose up -d
docker compose ps
```

影响范围：

- 拉取官方镜像 `mysql:8.4`。
- 创建容器 `interview-db-mysql`。
- 创建持久化卷 `interview-db-data`。
- 创建空数据库 `interview_db`。
- 仅监听本机 `127.0.0.1:3307`，可通过 `.env` 修改端口。

等待 `docker compose ps` 显示 `healthy`。需要排查启动问题时执行：

```powershell
docker compose logs --tail 200 mysql
```

此时只有空数据库，还没有招聘业务表和种子数据。

## 第二步：手动导入 SQL

先把 SQL 文件复制进容器：

```powershell
docker compose cp .\init\01-schema.sql mysql:/tmp/01-schema.sql
docker compose cp .\init\02-seed.sql mysql:/tmp/02-seed.sql
```

使用 root 进入 MySQL，命令会交互式询问 `.env` 中的 root 密码：

```powershell
docker compose exec mysql mysql -uroot -p
```

在 MySQL 提示符中按顺序执行：

```sql
SOURCE /tmp/01-schema.sql;
SOURCE /tmp/02-seed.sql;
```

不要颠倒顺序，也不要在同一个数据库中重复执行 `02-seed.sql`，否则唯一键会冲突。

也可以在 DBeaver 或 DataGrip 中以 root 连接后，依次打开并执行这两个文件。

## 第三步：验证导入结果

仍在 MySQL 提示符中执行：

```sql
USE interview_db;

SELECT 'candidates' AS entity, COUNT(*) AS row_count FROM candidates
UNION ALL SELECT 'applications', COUNT(*) FROM applications
UNION ALL SELECT 'interviews', COUNT(*) FROM interview_sessions
UNION ALL SELECT 'feedback', COUNT(*) FROM interview_feedback
UNION ALL SELECT 'competency_scores', COUNT(*) FROM feedback_competency_scores
UNION ALL SELECT 'offers', COUNT(*) FROM offers;
```

预期结果：

| entity | row_count |
| --- | ---: |
| candidates | 160 |
| applications | 220 |
| interviews | 297 |
| feedback | 526 |
| competency_scores | 5260 |
| offers | 44 |

再验证几个业务查询：

```sql
SELECT status, COUNT(*) AS application_count
FROM applications
GROUP BY status
ORDER BY application_count DESC;

SELECT *
FROM v_recruiting_funnel
ORDER BY total_applications DESC, requisition_no
LIMIT 10;

SELECT scheduled_start, candidate_name, job_title, interview_type, mode, status
FROM v_interview_schedule
WHERE scheduled_start >= '2026-07-23'
ORDER BY scheduled_start
LIMIT 20;
```

## 连接信息

| 参数 | 值 |
| --- | --- |
| Host | `127.0.0.1` |
| Port | `3307`，或 `.env` 中的 `MYSQL_PORT` |
| Database | `interview_db` |
| User | `interview_reader` |
| Password | `.env` 中的 `MYSQL_READER_PASSWORD` |
| Charset | `utf8mb4` |
| Time zone | `+08:00` |

导入 `01-schema.sql` 后，`interview_reader` 只拥有 `interview_db` 的 `SELECT` 权限：

```powershell
docker compose exec mysql mysql -uinterview_reader -p interview_db
```

```sql
SHOW GRANTS;
SELECT COUNT(*) FROM v_application_overview;
```

当前 `db-querry` 后端只实现了 PostgreSQL connector，暂时不能直接连接这个 MySQL
测试库。可先通过 MySQL CLI、DBeaver、DataGrip 或其他 MySQL 客户端使用。

## 停止和重新测试

停止但保留数据：

```powershell
docker compose stop
```

重新启动已有数据：

```powershell
docker compose start
```

如果要清空整个数据库并重新测试导入流程：

```powershell
docker compose down -v
docker compose up -d
```

警告：`docker compose down -v` 会永久删除 `interview-db-data` 卷中的全部数据。
重新启动后会再次得到一个空的 `interview_db`，然后重新执行手动导入步骤。
