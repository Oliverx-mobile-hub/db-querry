# Tasks：数据库查询工具

**Input**：`/specs/001-database-query-tool/` 下的 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`quickstart.md`、`contracts/openapi.yaml`

**Tests**：项目宪章要求关键后端路径必须有测试覆盖；本任务清单包含后端安全、API、metadata 和前端核心流程测试任务。

**Organization**：按 7 个 phase 拆分，先完成共享基础设施，再按可独立验证的功能增量推进。

## Format：`[ID] [P?] [Area] Description`

- **[P]**：可并行执行，前提是不同文件且无直接依赖。
- **[Area]**：任务所属区域，例如 `Setup`、`Store`、`API`、`Metadata`、`SQL`、`LLM`、`Frontend`、`Verify`。
- 每个任务必须包含明确文件路径。

## Phase 1：项目初始化

**Purpose**：建立 Go 后端和 Vue 前端的基本工程骨架，保证后续任务有稳定落点。

- [X] T001 [Setup] 初始化 Go module，创建 `backend/go.mod`，module 名称使用仓库匹配的 Go module 名。
- [X] T002 [Setup] 创建后端入口 `backend/cmd/server/main.go`，启动基础 HTTP server。
- [X] T003 [P] [Setup] 创建后端目录结构：`backend/internal/api/`、`backend/internal/config/`、`backend/internal/dbstore/`、`backend/internal/metadata/`、`backend/internal/pgconn/`、`backend/internal/query/`、`backend/internal/sqlguard/`、`backend/internal/llm/`。
- [X] T004 [Setup] 初始化前端 Vite + Vue 3 + TypeScript 工程文件：`frontend/package.json`、`frontend/index.html`、`frontend/vite.config.ts`、`frontend/tsconfig.json`。
- [X] T005 [P] [Setup] 创建前端基础入口：`frontend/src/main.ts`、`frontend/src/App.vue`、`frontend/src/api/client.ts`、`frontend/src/api/types.ts`。
- [X] T006 [P] [Setup] 添加前端样式入口：`frontend/src/styles/design-tokens.css` 和 `frontend/src/styles/global.css`。
- [X] T007 [Setup] 配置 `.gitignore`，确保忽略构建产物、依赖目录、临时文件和本地 SQLite 文件路径相关产物。

**Checkpoint**：后端和前端都具备最小可启动结构，尚不要求完整业务能力。

---

## Phase 2：基础设施

**Purpose**：实现所有后续功能依赖的配置、存储、路由、错误和测试基础。

**Critical**：本阶段完成前，不应开始业务 API 实现。

- [X] T008 [Config] 实现配置加载 `backend/internal/config/config.go`，读取监听地址、SQLite 路径默认值 `~/.db_querry/db_querrt.db`、`openai_api_key`。
- [X] T009 [Store] 实现 SQLite 路径解析和目录创建 `backend/internal/dbstore/store.go`。
- [X] T010 [Store] 实现 SQLite migration `backend/internal/dbstore/migrations.go`，创建 `db_connections`、`metadata_snapshots`、`generated_sql_drafts`。
- [X] T011 [P] [Store] 为 SQLite migration 添加测试 `backend/internal/dbstore/migrations_test.go`，验证表和索引可创建。
- [X] T012 [API] 实现统一 API envelope 和错误模型 `backend/internal/api/responses.go`，字段名使用 camelCase。
- [X] T013 [API] 实现基础 router 和 CORS middleware `backend/internal/api/routes.go`，允许所有 origin。
- [X] T014 [P] [API] 添加 CORS 和 error envelope 测试 `backend/internal/api/routes_test.go`。
- [X] T015 [P] [Model] 定义共享 Go 类型 `backend/internal/api/types.go`，覆盖 DbSummary、Metadata、QueryResult、SqlValidationResult、GeneratedSqlDraft。
- [X] T016 [P] [Frontend] 定义 TypeScript API 类型 `frontend/src/api/types.ts`，与 `contracts/openapi.yaml` 对齐。
- [X] T017 [Frontend] 实现 API client `frontend/src/api/client.ts`，处理 envelope、错误状态和 JSON 请求。

**Checkpoint**：SQLite、配置、CORS、统一响应、前后端类型基础就绪。

---

## Phase 3：数据库连接与 metadata

**Goal**：用户能添加 PostgreSQL 连接，系统采集并保存 metadata，前端能浏览连接和表/视图信息。

**Independent Test**：调用 `PUT /api/v1/dbs/local` 添加连接，再调用 `GET /api/v1/dbs` 和 `GET /api/v1/dbs/local` 获取脱敏连接信息和 metadata。

### Tests

- [ ] T018 [P] [Metadata] 编写 metadata collector 单元测试 `backend/internal/metadata/collector_test.go`，覆盖 schema/table/view/column/primary key/comment 映射。
- [ ] T019 [P] [API] 编写连接 API handler 测试 `backend/internal/api/dbs_handlers_test.go`，覆盖成功、连接失败、URL 脱敏和 db not found。
- [X] T020 [P] [Store] 编写连接和 metadata store 测试 `backend/internal/dbstore/store_test.go`，覆盖 upsert connection、insert snapshot、latest snapshot。

### Implementation

- [X] T021 [Store] 实现连接存储方法 `backend/internal/dbstore/store.go`：upsert、list、get、update metadata status。
- [X] T022 [Store] 实现 metadata snapshot 存储方法 `backend/internal/dbstore/store.go`：insert latest、get latest。
- [X] T023 [PG] 实现 PostgreSQL connector `backend/internal/pgconn/connector.go`，支持连接测试和只在后端使用完整 URL。
- [X] T024 [Metadata] 定义 metadata Go model `backend/internal/metadata/model.go`，与 `data-model.md` 的 JSON schema 对齐。
- [X] T025 [Metadata] 实现 PostgreSQL metadata collector `backend/internal/metadata/collector.go`，采集 schema、table、view、column、nullable、primaryKey、comment。
- [X] T026 [API] 实现 `GET /api/v1/dbs` handler `backend/internal/api/dbs_handlers.go`。
- [X] T027 [API] 实现 `PUT /api/v1/dbs/{name}` handler `backend/internal/api/dbs_handlers.go`，完成连接测试、metadata 采集和 SQLite 保存。
- [X] T028 [API] 实现 `GET /api/v1/dbs/{name}` handler `backend/internal/api/dbs_handlers.go`，返回最新 metadata。
- [X] T029 [Frontend] 实现数据库连接面板 `frontend/src/components/ConnectionPanel.vue`，支持连接列表、添加连接 dialog、状态 tag。
- [X] T030 [Frontend] 实现 metadata 浏览组件 `frontend/src/components/MetadataExplorer.vue`，展示 schema、table/view、columns。
- [X] T031 [Frontend] 在 `frontend/src/App.vue` 集成连接列表和 metadata explorer，选中连接后加载 metadata。

**Checkpoint**：连接和 metadata 功能可独立演示，不依赖 SQL 查询或 LLM。

---

## Phase 4：SQL 校验与查询执行

**Goal**：用户能提交只读 SELECT 查询，后端执行前强制校验并默认限制 1000 行，前端展示 JSON 表格结果。

**Independent Test**：提交合法 SELECT、无 LIMIT SELECT、语法错误、非 SELECT、多语句 SQL，验证响应和数据库未执行危险语句。

### Tests

- [X] T032 [P] [SQL] 编写 SQL Guard 测试 `backend/internal/sqlguard/validator_test.go`，覆盖合法 SELECT、无 LIMIT、INSERT、UPDATE、DELETE、DROP、ALTER、TRUNCATE、CREATE、REPLACE、GRANT、REVOKE、EXEC、多语句、语法错误。
- [ ] T033 [P] [Query] 编写查询执行测试 `backend/internal/query/executor_test.go`，覆盖 JSON 友好类型转换、空结果、超时错误脱敏。
- [ ] T034 [P] [API] 编写 query API handler 测试 `backend/internal/api/query_handlers_test.go`，覆盖 envelope、默认 LIMIT、校验失败和 db not found。

### Implementation

- [X] T035 [SQL] 实现 SQL Guard 类型 `backend/internal/sqlguard/validator.go`，返回 `SqlValidationResult`。
- [X] T036 [SQL] 实现单语句 SELECT 校验 `backend/internal/sqlguard/validator.go`，拒绝非 SELECT、多语句和 parser 错误。
- [X] T037 [SQL] 实现禁止语句和高风险语法策略 `backend/internal/sqlguard/validator.go`，至少覆盖 spec 中列出的危险语句和 `SELECT ... INTO`。
- [X] T038 [SQL] 实现默认 `LIMIT 1000` 处理 `backend/internal/sqlguard/validator.go`，响应中标记 `limitApplied` 和 `limit`。
- [X] T039 [Query] 实现 query executor `backend/internal/query/executor.go`，按连接名读取 URL，带 timeout 执行只读查询。
- [X] T040 [Query] 实现查询结果转换 `backend/internal/query/result.go`，输出 columns、rows、rowCount、durationMs、limitApplied、limit、empty。
- [X] T041 [API] 实现 `POST /api/v1/dbs/{name}/query` handler `backend/internal/api/query_handlers.go`。
- [X] T042 [Frontend] 实现 SQL 编辑组件 `frontend/src/components/QueryEditor.vue`，包含 SQL 输入、执行按钮、校验/执行状态。
- [X] T043 [Frontend] 实现结果表格 `frontend/src/components/ResultTable.vue`，根据 JSON columns/rows 动态渲染 Element Plus table。
- [X] T044 [Frontend] 在 `frontend/src/App.vue` 集成 QueryEditor 和 ResultTable。

**Checkpoint**：手写 SQL 查询闭环可用，核心安全路径有测试覆盖。

---

## Phase 5：自然语言生成 SQL

**Goal**：用户能基于 metadata 输入自然语言，后端调用 OpenAI 生成 SQL 草稿，返回预览和校验状态，但不自动执行。

**Independent Test**：对已有 metadata 的连接提交 prompt，返回 SQL 草稿；LLM 返回非法 SQL 时响应不可执行状态。

### Tests

- [ ] T045 [P] [LLM] 编写 OpenAI client mock 测试 `backend/internal/llm/openai_test.go`，覆盖成功、超时、非 JSON 输出、缺失 API key。
- [ ] T046 [P] [API] 编写 natural query handler 测试 `backend/internal/api/natural_query_handlers_test.go`，覆盖 `prompt`、兼容 `promt`、metadata 缺失、LLM 输出非法 SQL。
- [X] T047 [P] [Store] 编写 generated SQL draft store 测试 `backend/internal/dbstore/generated_sql_drafts_test.go`。

### Implementation

- [X] T048 [LLM] 实现 OpenAI client `backend/internal/llm/openai.go`，从 `openai_api_key` 读取 key，不向前端暴露。
- [X] T049 [LLM] 实现 metadata context 构造 `backend/internal/llm/prompt_builder.go`，限制上下文大小并明确只读 SQL 输出要求。
- [X] T050 [LLM] 实现 LLM 结构化输出解析 `backend/internal/llm/sql_draft.go`，只接受包含 sql/explanation 的 JSON。
- [X] T051 [Store] 实现 generated SQL draft 持久化 `backend/internal/dbstore/store.go`。
- [X] T052 [API] 实现 `POST /api/v1/dbs/{name}/query/natural` handler `backend/internal/api/natural_query_handlers.go`，支持 `prompt` 并兼容 `promt`。
- [X] T053 [API] 将 LLM 生成 SQL 接入 SQL Guard `backend/internal/api/natural_query_handlers.go`，返回 validation，不执行查询。
- [X] T054 [Frontend] 实现自然语言面板 `frontend/src/components/NaturalLanguagePanel.vue`，支持 prompt 输入、生成 SQL、展示解释和校验状态。
- [X] T055 [Frontend] 在 `frontend/src/App.vue` 集成自然语言生成结果，将生成 SQL 填入 QueryEditor 供用户执行。

**Checkpoint**：自然语言生成 SQL 可用，且无法绕过 SQL Guard。

---

## Phase 6：前端 dashboard 与视觉风格

**Goal**：按 MotherDuck-inspired / brutalist dashboard 风格完成可用界面，不保留默认 Element Plus 后台模板感。

**Independent Test**：打开前端，能完成连接、浏览 metadata、输入 SQL、自然语言生成和查看结果；视觉符合 plan 中 token 和布局要求。

### Tests

- [X] T056 [P] [Frontend] 编写 API client 单元测试 `frontend/tests/unit/apiClient.test.ts`，覆盖 success/failure envelope。
- [ ] T057 [P] [Frontend] 编写组件状态测试 `frontend/tests/unit/queryFlow.test.ts`，覆盖执行 SQL 成功、错误和空结果状态。
- [X] T058 [P] [Frontend] 编写 metadata explorer 测试 `frontend/tests/unit/metadataExplorer.test.ts`。

### Implementation

- [X] T059 [Frontend] 在 `frontend/src/styles/design-tokens.css` 定义奶油背景、白色面板、黑色 2px 边框、蓝色主按钮、黄色强调、小圆角和硬投影 token。
- [X] T060 [Frontend] 在 `frontend/src/styles/global.css` 覆盖 Element Plus button/input/table/dialog/tag 样式，保持 brutalist dashboard 风格。
- [X] T061 [Frontend] 重构 `frontend/src/App.vue` 为三栏 dashboard：左侧连接区、中间查询区、右侧 metadata explorer、下方结果表。
- [X] T062 [Frontend] 为 `ConnectionPanel.vue` 应用硬边框、小圆角、uppercase 标签和状态色。
- [X] T063 [Frontend] 为 `QueryEditor.vue` 应用深色 SQL 编辑区和清晰执行状态。
- [X] T064 [Frontend] 为 `NaturalLanguagePanel.vue` 应用黄色强调和 SQL 草稿预览样式。
- [X] T065 [Frontend] 为 `ResultTable.vue` 应用硬边框表格、空状态和错误状态样式。
- [X] T066 [Frontend] 补齐响应式布局 `frontend/src/styles/global.css`，确保窄屏下不发生文本和面板重叠。

**Checkpoint**：前端功能和视觉风格达到 plan 要求。

---

## Phase 7：集成验证与收尾

**Purpose**：确认实现与 spec/plan/contracts 一致，补齐文档和端到端验证。

- [X] T067 [Verify] 对照 `specs/001-database-query-tool/contracts/openapi.yaml` 检查后端所有 endpoint、字段名和错误 envelope。
- [X] T068 [Verify] 运行后端测试 `go test ./...`，记录并修复失败。
- [X] T069 [Verify] 运行前端测试 `npm test` 或项目实际测试命令，记录并修复失败。
- [ ] T070 [Verify] 按 `specs/001-database-query-tool/quickstart.md` 手动验证添加连接、查询列表、metadata、SQL 查询、禁止语句、自然语言生成 SQL。
- [X] T071 [Security] 检查 API 响应和前端代码，确认不暴露完整 DB URL、密码、`openai_api_key`、堆栈或原始连接详情。
- [X] T072 [Security] 检查 CORS、无认证、只读 SQL、timeout、默认 LIMIT 组合风险，确认与宪章一致。
- [X] T073 [Docs] 更新 `specs/001-database-query-tool/quickstart.md` 中实际启动命令、端口和已知限制。
- [X] T074 [Docs] 如实现偏离 `data-model.md` 或 `contracts/openapi.yaml`，先更新对应设计文档再提交实现。
- [ ] T075 [Verify] 生成最终验收记录，列出已运行测试、未运行测试及原因。

**Checkpoint**：功能可按 quickstart 验证，关键安全路径有测试和人工确认。

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 项目初始化**：无依赖。
- **Phase 2 基础设施**：依赖 Phase 1，阻塞所有业务 phase。
- **Phase 3 数据库连接与 metadata**：依赖 Phase 2。
- **Phase 4 SQL 校验与查询执行**：依赖 Phase 2；Query executor 依赖 Phase 3 的连接存储能力。
- **Phase 5 自然语言生成 SQL**：依赖 Phase 3 metadata 和 Phase 4 SQL Guard。
- **Phase 6 前端 dashboard 与视觉风格**：依赖 Phase 2 API client；可在 Phase 3-5 后端接口稳定后逐步集成。
- **Phase 7 集成验证与收尾**：依赖目标功能 phase 完成。

### Parallel Opportunities

- T003、T005、T006 可并行。
- T011、T014、T015、T016 可并行。
- Phase 3 的测试 T018-T020 可并行，随后实现 T021-T028。
- Phase 4 的测试 T032-T034 可并行，随后实现 SQL Guard 和 Query Executor。
- Phase 5 的测试 T045-T047 可并行。
- Phase 6 的组件样式任务 T062-T065 可在 design tokens 完成后并行。

### MVP Path

最小可用版本建议按以下顺序完成：

1. Phase 1。
2. Phase 2。
3. Phase 3。
4. Phase 4。
5. Phase 7 中 T067、T068、T070、T071。

此时用户可以添加数据库、浏览 metadata、执行安全 SELECT 并查看结果。Phase 5 和完整 Phase 6 可作为下一增量继续。

## Notes

- 不实现 authentication；如果后续需要用户身份或权限，必须先更新 spec 和 plan。
- 前端不得复制后端 SQL 安全规则，只能做轻量提示。
- 所有新增 API 字段保持 camelCase。
- `promt` 仅作为输入兼容别名，正式类型和响应统一使用 `prompt`。
- SQLite 文件包含敏感连接 URL，任何日志、错误和 API 响应都不得输出完整 URL。

