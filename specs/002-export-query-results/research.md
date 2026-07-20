# Research：导出当前查询结果为 CSV/JSON

**日期**：2026-07-20

**关联计划**：[plan.md](./plan.md)

## 决策 1：v1 导出在前端完成，不新增后端导出 API

**Decision**：CSV/JSON 导出使用当前前端持有的 `QueryResult`，通过浏览器 Blob 下载，不新增 `/api/v1/.../export` 后端接口。

**Rationale**：需求是“导出当前查询结果”，当前结果已经由后端查询 API 返回给前端。前端本地转换可以避免重新执行 SQL，避免引入新的后端安全边界，也不会绕过默认 LIMIT。

**Alternatives considered**：

- 新增后端 `POST /api/v1/dbs/{name}/query/export`：适合导出完整大结果集或后端流式下载，但 v1 不需要重新执行 SQL。
- 将导出历史保存到 SQLite：当前没有审计、历史下载或共享需求，增加存储复杂度。

## 决策 2：JSON 导出完整当前 QueryResult

**Decision**：JSON 文件导出完整当前 `QueryResult`，包括 `columns`、`rows`、`rowCount`、`durationMs`、`limitApplied`、`limit`、`empty` 和可选 `validation`。

**Rationale**：完整 `QueryResult` 能保留列定义、行数据、执行耗时、LIMIT 状态和 SQL 校验上下文，更适合调试和后续程序处理。

**Alternatives considered**：

- 只导出 rows：文件更简洁，但丢失列类型、行数、LIMIT 和 validation 信息。
- 导出 `{columns, rows}`：比 rows 更完整，但仍丢失执行元信息。

## 决策 3：CSV 按 columns 顺序导出 rows

**Decision**：CSV 表头来自 `QueryResult.columns[].name`，每行按相同顺序读取 `row[column.name]`。

**Rationale**：JavaScript 对象 key 顺序不应作为业务契约；后端已经返回 columns，前端必须使用这个稳定顺序生成 CSV。

**Alternatives considered**：

- 从第一行对象 key 推断列顺序：空结果无法推断，且不同记录字段顺序可能不一致。
- 让用户选择列：可以作为后续高级功能，v1 不需要。

## 决策 4：CSV 转义由项目内 TypeScript utility 实现

**Decision**：新增小型 `exportResults` utility，手写 CSV 字段序列化和文件名生成，不引入第三方库。

**Rationale**：需求范围小，CSV 规则明确；手写 utility 便于测试，避免增加依赖和打包体积。

**Alternatives considered**：

- 引入 PapaParse 或类似 CSV 库：功能完整但超出 v1 需求。
- 在组件内直接拼字符串：会让 `ResultTable` 变重，测试也不够聚焦。

## 决策 5：禁用状态由 ResultTable 的输入状态决定

**Decision**：`ResultTable` 接收 `result`、`loading` 和 `dbName`，当 `loading=true`、`result=null` 或 `result.empty=true` 时禁用导出按钮。

**Rationale**：按钮状态应由明确 props 决定，符合宪章中“前端状态清晰、可预测”的要求。

**Alternatives considered**：

- 点击后弹错误：用户体验差，且不如禁用状态直观。
- 组件内部读取全局查询状态：会增加隐式状态依赖。

## 决策 6：文件名使用安全化数据库名和本地时间戳

**Decision**：文件名格式为 `{safeDbName}-query-{YYYYMMDD-HHmmss}.{csv|json}`，数据库名中的非字母数字、下划线、短横线字符统一替换为 `-`。

**Rationale**：用户要求文件名包含数据库名和时间；安全化处理可避免路径分隔符、空格和特殊字符造成下载异常。

**Alternatives considered**：

- 使用原始数据库名：可能包含文件系统不友好的字符。
- 只使用固定文件名：会覆盖用户下载历史，不利于识别来源。

## 决策 7：后端日志输出到 stdout，使用完整时间戳和非敏感字段

**Decision**：后端启动日志必须输出到 stdout，记录非敏感运行配置摘要，例如 `base_url`、`model`、`wire_api`、`key_loaded` 和监听地址；禁止输出完整 API key、数据库密码或完整连接 URL。

**Rationale**：K8s 和云厂商通常从容器 stdout/stderr 采集日志。用户提供的参考实现使用 logrus、stdout、DebugLevel、完整时间戳格式 `2006-01-02 15:04:05` 和 `Debug/Info/Warning/Error` wrapper，符合部署采集需求。

**Alternatives considered**：

- 只用默认 `log.Printf`：实现简单，但字段化程度弱；短期可用，后续部署可升级为 logging wrapper。
- 写入本地文件：容器环境不利于采集，也会增加文件轮转问题。
- 输出 key 或完整 DSN 方便调试：违反宪章的敏感信息保护要求。

## 决策 8：后端 API contract 不变

**Decision**：本功能不改变现有 `/api/v1/dbs/{name}/query` 响应结构，导出以现有 `QueryResult` 为 contract。

**Rationale**：导出当前结果是前端本地行为；保持后端 API 不变可以降低回归风险。

**Alternatives considered**：

- 在 query response 中增加 export metadata：当前文件名所需信息前端已有，不需要后端新增字段。
