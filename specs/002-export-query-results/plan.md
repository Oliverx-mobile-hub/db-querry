# 实现计划：导出当前查询结果为 CSV/JSON

**分支**：`002-export-query-results` | **日期**：2026-07-20 | **规格**：[spec.md](./spec.md)

**输入**：来自 `/specs/002-export-query-results/spec.md` 的功能规格，以及本轮补充约束：

- Results 面板增加 `Export CSV`、`Export JSON` 两个按钮。
- 无结果、空结果、查询中时按钮禁用。
- JSON 导出完整当前 `QueryResult`。
- CSV 按 `QueryResult.columns` 顺序导出 `rows`，并正确处理逗号、双引号、换行、null、日期和数字。
- 文件名包含数据库名和时间，例如 `local-query-20260720-143000.csv`。
- 单元测试覆盖 CSV 转义和按钮可用状态。
- 后端和前端继续参考现有实现风格；API 契约继续沿用 `/api/v1`、camelCase 和统一响应 envelope。
- 后端日志需要输出到启动日志，方便本地查看、K8s 部署和云厂商日志采集；格式参考用户提供的 logrus stdout/time/full timestamp 风格。

## 摘要

本功能为查询结果增加本地导出能力。用户在执行 SQL 并获得非空结果后，可以从 Results 面板导出 CSV 或 JSON 文件。

技术路线：

1. 前端在 `ResultTable` 接收当前 `QueryResult`、当前数据库名和查询 loading 状态。
2. 前端提供 `Export CSV`、`Export JSON` 两个按钮，禁用状态由当前结果状态决定。
3. JSON 导出完整当前 `QueryResult`，不重新请求后端。
4. CSV 使用纯 TypeScript utility 从 `columns + rows` 生成，严格按列顺序输出，处理 RFC 4180 常见转义。
5. 浏览器侧用 Blob + object URL 触发下载，文件名由数据库名和本地时间戳生成。
6. 后端不新增导出 API；现有查询 API、SQL 安全和 LIMIT 行为保持不变。
7. 后端日志在实现阶段继续向 stdout 输出启动配置摘要，使用结构化字段或接近 logrus text formatter 的可采集格式，禁止输出密钥。

## 技术上下文

**语言 / 版本**：

- 后端：Go 1.23+。
- 前端：TypeScript + Vue 3。

**主要依赖**：

- 前端：Vue 3、Element Plus、Vite、TypeScript、Vitest、Vue Test Utils。
- 后端：Go 标准库、现有 API / query / config / logging 结构。
- 导出：优先使用浏览器原生 `Blob`、`URL.createObjectURL`、临时 `<a download>`，不引入第三方 CSV 库。
- 日志：后端可继续使用标准库 `log` 或在实现阶段引入轻量 logging wrapper；如引入 logrus，需说明收益，并保持 stdout 输出。

**存储**：

- 不新增持久化存储。
- 导出文件仅由浏览器本地下载生成。
- 不写入 SQLite，不保存导出历史。

**测试**：

- 前端：Vitest + Vue Test Utils。
- 必测：CSV 转义、JSON payload、文件名格式、安全化数据库名、按钮禁用状态。
- 后端：如调整启动日志，运行 `go test ./...`，并验证日志不泄露 key。

**目标平台**：

- 本地 Web 应用：Go API 服务 + Vite/Vue 前端。
- 浏览器下载行为覆盖现代 Chromium 浏览器；不要求兼容 IE。

**项目类型**：

- Web application：`backend/` + `frontend/`。

**性能目标**：

- 导出当前前端内存中的结果，不重新执行 SQL。
- 对当前默认 LIMIT 结果规模，CSV/JSON 生成应在用户感知上即时完成。
- 导出过程不得阻塞后续查询状态更新。

**约束**：

- 不绕过后端 SQL 只读校验和默认 LIMIT。
- 前端不得保存或导出数据库连接 URL、LLM key、后端环境变量或高权限配置。
- 所有新增 TypeScript 类型必须明确。
- Results 面板视觉继续沿用现有 MotherDuck-inspired / brutalist dashboard 风格：硬边框、小圆角、Element Plus button、紧凑信息密度。
- 后端日志必须输出到 stdout，便于 K8s / 云日志采集；日志不得打印 API key、数据库密码或完整连接 URL。

**规模 / 范围**：

- v1 只导出当前 ResultTable 中已有结果。
- v1 支持 CSV 和 JSON 两种格式。
- v1 不支持导出全部未分页数据、不支持服务端异步导出、不支持导出历史记录。

## 宪章检查

*GATE：Phase 0 research 前必须通过；Phase 1 design 后复查。*

- **只读数据库安全**：通过。导出只使用后端已返回的当前结果，不新增 SQL 执行路径，不绕过 SQL Guard。
- **Go 后端负责业务逻辑**：通过。查询执行、SQL 校验和敏感配置仍由 Go 后端负责；导出只是前端对已返回结果的本地文件转换。
- **Vue 3 + Element Plus 薄客户端**：通过。前端负责展示和下载交互，不承担数据库访问或安全校验。
- **API 契约优先**：通过。v1 不新增后端 export API，contracts 明确无 API 变更；现有 query response 是导出数据源。
- **严格类型**：通过。新增 exporter utility、props、button state、测试 fixture 均使用 TypeScript 类型。
- **测试关键路径**：通过。CSV 转义和按钮状态列为必须测试；后端如改日志需跑 Go 测试。

当前无宪章违规项。

## 项目结构

### 文档结构

```text
specs/002-export-query-results/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
    └── openapi.yaml
```

### 源码结构

```text
frontend/
├── src/
│   ├── components/
│   │   └── ResultTable.vue
│   ├── utils/
│   │   └── exportResults.ts
│   └── api/
│       └── types.ts
└── tests/
    └── unit/
        ├── exportResults.test.ts
        └── ResultTable.test.ts

backend/
├── cmd/
│   └── server/
│       └── main.go
└── internal/
    └── config/
```

**Structure Decision**：导出功能主要落在前端 `ResultTable.vue` 和 `frontend/src/utils/exportResults.ts`；后端不新增业务 API，只保留或整理启动日志，确保运行时能看到 base_url、model、wire_api、key_loaded 等非敏感信息。

## Phase 0：Research 输出

见 [research.md](./research.md)。

## Phase 1：Design 输出

- 数据模型：[data-model.md](./data-model.md)
- API / contract：[contracts/openapi.yaml](./contracts/openapi.yaml)
- 验证指南：[quickstart.md](./quickstart.md)

## Complexity Tracking

当前无宪章违规项，无需复杂性豁免。
