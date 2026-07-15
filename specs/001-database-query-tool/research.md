# Research：数据库查询工具

**日期**：2026-07-15

**关联计划**：[plan.md](./plan.md)

## 决策 1：v1 仅支持 PostgreSQL / PostgreSQL 兼容数据库

**Decision**：v1 连接目标限定为 PostgreSQL 或 PostgreSQL 兼容数据库，DB URL 使用 `postgres://` 或 `postgresql://`。

**Rationale**：需求明确提到“根据 PostgreSQL 的功能查询系统表和视图的信息”。先固定一个方言可以让 metadata 采集、SQL parser、只读校验和错误提示保持一致。

**Alternatives considered**：

- 同时支持 MySQL / SQLite / SQL Server：会显著扩大 metadata 查询、SQL 方言和 parser 复杂度，不适合第一个 feature。
- 抽象多数据库 connector：当前没有第二种数据库的明确需求，过早抽象会增加实现成本。

## 决策 2：metadata 采集使用 PostgreSQL catalog + information_schema

**Decision**：metadata collector 以 `pg_catalog` 为主，必要时结合 `information_schema`，采集 schema、table、view、column、data type、nullable、primary key、comment。

**Rationale**：`pg_catalog` 能提供更完整的 PostgreSQL 信息，例如 comment、primary key 和对象类型；`information_schema` 可作为标准化补充。

**Alternatives considered**：

- 只用 `information_schema`：跨数据库标准更强，但 PostgreSQL 注释、对象分类和约束信息不够完整。
- 只让 LLM 从原始 catalog 文本推断：LLM 输出不可靠，不能作为唯一结构化来源。

## 决策 3：LLM 只参与生成 SQL 和可选 metadata 整理，不作为可信来源

**Decision**：LLM 生成的 SQL 和 metadata JSON 都视为不可信输入，必须解析为严格 Go 类型并经过后端校验。

**Rationale**：项目宪章要求所有生成 SQL 都必须当作不可信输入。LLM 适合生成草稿和说明，不适合作为权限边界。

**Alternatives considered**：

- LLM 直接返回可执行 SQL 并执行：违反只读安全原则。
- LLM 负责完整 metadata 结构定义：可用性和一致性难以保障，应由后端结构化 collector 兜底。

## 决策 4：SQL 校验使用 PostgreSQL 方言 parser + 策略层

**Decision**：SQL Guard 分两层：先用 PostgreSQL 方言 parser 验证语法和 AST，再用策略层限制为单条只读 SELECT。

**Rationale**：单纯字符串匹配无法可靠处理注释、字符串字面量、CTE、子查询、大小写和多语句。parser 能提供语法级判断，策略层负责产品安全规则。

**Alternatives considered**：

- 正则或关键字黑名单：实现快但安全性不足。
- 将 SQL 发给数据库 `EXPLAIN` 做验证：仍可能触发数据库解析风险，并且不能替代本地策略判断。

## 决策 5：无 LIMIT 时默认应用 `LIMIT 1000`

**Decision**：如果合法 SELECT 查询没有 top-level `LIMIT`，后端在执行前应用 `LIMIT 1000`，并在响应中返回 `limitApplied=true` 和 `limit=1000`。

**Rationale**：需求明确指定默认限制，且无认证场景下必须控制结果规模。

**Alternatives considered**：

- 前端提示用户补 LIMIT：不可靠，不能作为安全边界。
- 数据库 cursor 分页：更复杂，可作为后续优化。

## 决策 6：SQLite 文件路径固定为 `~/.db_querry/db_querrt.db`

**Decision**：后端启动时解析用户 home 目录，创建 `~/.db_querry/`，并初始化 `db_querrt.db`。

**Rationale**：这是用户明确指定的路径。路径中的 `db_querrt.db` 按原样保留，避免实现和需求不一致。

**Alternatives considered**：

- 使用项目目录内 SQLite：不符合用户指定路径。
- 使用环境变量覆盖路径：可作为后续增强，但 v1 固定默认路径。

## 决策 7：DB URL v1 存 SQLite，但只允许后端读取

**Decision**：v1 将完整 DB URL 存储在 SQLite 的 `db_connections.url` 字段；任何 API 响应只能返回脱敏后的 `displayDsn`。

**Rationale**：用户明确要求连接字符串存储在 SQLite。当前项目无认证，本地工具优先完成可用闭环。

**Risk**：SQLite 文件中包含数据库凭据。后续可以升级为本地加密、操作系统密钥管理或只保存环境变量引用。

**Mitigation**：

- 前端和 API 永不返回完整 URL。
- 错误脱敏。
- quickstart 中提醒 SQLite 文件包含敏感信息。

## 决策 8：CORS v1 允许所有 origin

**Decision**：后端对 `/api/v1/*` 启用 CORS，允许所有 origin、常用方法和 `Content-Type`。

**Rationale**：用户明确要求允许所有 origin，适合本地开发阶段。

**Risk**：无认证 + 全开放 CORS 会放大本机服务被其他网页调用的风险。

**Mitigation**：

- 只读 SQL 安全、查询超时、默认 LIMIT、错误脱敏必须严格执行。
- 后续若部署到共享环境，必须先更新 spec 和 plan，再收紧 CORS 和访问控制。

## 决策 9：自然语言 API 字段使用 `prompt`，兼容 `promt`

**Decision**：正式 contract 使用 `prompt`；后端 DTO 可兼容读取 `promt`。

**Rationale**：用户草案中写作 `promt`，但正式 API 不应扩散拼写错误。兼容读取可以降低早期调用失败率。

**Alternatives considered**：

- 只使用 `promt`：长期维护成本高。
- 只接受 `prompt`：与用户草案不完全兼容。

## 决策 10：前端采用 MotherDuck-inspired / brutalist dashboard 风格

**Decision**：前端设计 token 固定为奶油背景、白色面板、黑色 2px 硬边框、小圆角、蓝色主按钮、黄色强调、uppercase 标签、深色 SQL 编辑区。

**Rationale**：用户明确指定参考会话 `019f64f6-a164-7fe1-8f9c-bced42126aef` 的风格总结。

**Alternatives considered**：

- Element Plus 默认后台风格：与用户指定风格不一致。
- 纯 Tailwind 自定义 UI：当前技术约束要求 Vue 3 + Element Plus。
