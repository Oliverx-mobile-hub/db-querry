# Tasks：MySQL 数据库查询支持

**Input**：`/specs/003-mysql-db-query/` 下的 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`quickstart.md`、`contracts/openapi.yaml`

**Tests**：本功能要求覆盖 MySQL 连接与 metadata、SQL Guard、查询执行、自然语言生成 SQL、PostgreSQL 回归、日志格式和 Chrome 端到端验证。

**Organization**：按 7 个 phase 拆分，并按用户故事组织，确保每个故事都可以独立实现和验证。

## Format：`[ID] [P?] [Story] Description`

- **[P]**：可并行执行，前提是不同文件且无直接依赖。
- **[Story]**：任务所属用户故事，例如 `US1`、`US2`、`US3`、`US4`。
- 每个任务必须包含明确文件路径。

---

## Phase 1：基础类型与依赖

**Purpose**：把 MySQL 支持所需的依赖、类型和前后端共享模型先对齐。

- [ ] T001 [P] [Setup] 在 `backend/go.mod` 和 `backend/go.sum` 中增加 `github.com/go-sql-driver/mysql` 依赖，确保后端可以编译 MySQL driver。
- [ ] T002 [P] [Setup] 在 `backend/internal/api/types.go` 中把 `DatabaseType`、`DBConnectionRecord`、`DBSummary`、`MetadataDocument` 收敛成明确类型，并补充 `postgres` / `mysql` 常量。
- [ ] T003 [P] [Setup] 在 `frontend/src/api/types.ts` 中将 `DbSummary.databaseType`、`MetadataDocument.databaseType`、`PutDbRequest.databaseType` 扩展为 `postgres | mysql`，并同步 `frontend/src/api/client.ts` 的请求/响应类型引用。

**Checkpoint**：前后端共享类型已经支持 MySQL，后续实现可以按数据库类型分发。

---

## Phase 2：基础设施与分发

**Purpose**：建立 MySQL 连接、metadata、SQL Guard、查询执行的分发骨架，并保持日志和存储兼容。

- [ ] T004 [P] [Foundation] 新增 `backend/internal/mysqlconn/connector.go`，实现 MySQL URL -> DSN 转换、连接测试入口和错误脱敏。
- [ ] T005 [P] [Foundation] 在 `backend/internal/metadata/collector.go` 中加入 `databaseType` 分发，并新增 `backend/internal/metadata/mysql_collector.go` 的接口骨架。
- [ ] T006 [P] [Foundation] 在 `backend/internal/sqlguard/validator.go` 中加入 `databaseType` 分发，并新增 `backend/internal/sqlguard/mysql_validator.go` 的接口骨架。
- [ ] T007 [P] [Foundation] 在 `backend/internal/query/executor.go` 中加入 `databaseType` 分发，并新增 `backend/internal/query/mysql_executor.go` 的接口骨架。
- [ ] T008 [P] [Foundation] 在 `backend/internal/dbstore/store.go`、`backend/internal/dbstore/migrations.go` 中确认 `database_type` 读写、列表和详情读取兼容 MySQL 记录。
- [ ] T009 [P] [Foundation] 在 `backend/cmd/server/main.go` 和 `backend/internal/logging/logging.go` 中保留单行 stdout 日志格式，继续输出启动摘要且不暴露 key、密码或完整 URL。
- [ ] T010 [P] [Foundation] 更新 `backend/internal/logging/logging_test.go`，确认 `time="..." level=... msg="..."` 格式、`key_loaded` 摘要和 MySQL 流程日志不会泄露敏感信息。

**Checkpoint**：后端已经具备 MySQL 分发入口、日志约束和存储兼容基础。

---

## Phase 3：User Story 1 - 添加 MySQL 连接并采集 metadata (Priority: P1)

**Goal**：用户可以新增 `interview_db` 这样的 MySQL 连接，系统能采集 metadata 并在前端浏览。

**Independent Test**：添加一个有效的 MySQL 连接后，刷新数据库列表，确认连接状态、metadata 状态和 metadata 浏览树正常显示。

### Tests for User Story 1

- [ ] T011 [P] [US1] 新增 `backend/internal/metadata/mysql_collector_test.go`，覆盖 `information_schema` 到 `MetadataDocument` 的映射、schema/table/view/column、nullable、primary key 和 comment。
- [ ] T012 [P] [US1] 扩展 `backend/internal/api/dbs_handlers_test.go`，覆盖 `PUT /api/v1/dbs/{name}` 的 `databaseType=mysql`、连接失败、脱敏 URL、`GET /api/v1/dbs/{name}` metadata 状态与 `GET /api/v1/dbs` 列表。

### Implementation for User Story 1

- [ ] T013 [US1] 在 `backend/internal/mysqlconn/connector.go` 完成 MySQL URL 解析、driver DSN 生成、连接测试和 `interview_db` 场景支持。
- [ ] T014 [US1] 在 `backend/internal/metadata/mysql_collector.go` 完成基于 `information_schema` 的 metadata 采集，并把结果写回 `backend/internal/dbstore/store.go`。
- [ ] T015 [US1] 更新 `backend/internal/api/dbs_handlers.go`，在添加连接后保存 MySQL metadata、设置 `metadataStatus` / `connectionStatus` 并返回脱敏 summary。
- [ ] T016 [US1] 更新 `frontend/src/components/ConnectionPanel.vue`、`frontend/src/components/MetadataExplorer.vue`、`frontend/src/App.vue`，增加数据库类型选择、MySQL 状态展示和 metadata 树浏览。

**Checkpoint**：MySQL 数据库连接和 metadata 浏览可以独立演示。

---

## Phase 4：User Story 2 - MySQL 手写 SQL 查询 (Priority: P1)

**Goal**：用户可以对 MySQL 数据库执行只读 `SELECT` 查询，并在前端看到结果。

**Independent Test**：对 `interview_db` 执行合法 `SELECT`、无 `LIMIT` `SELECT`、非 `SELECT` 和多语句 SQL，验证结果、默认限制和拒绝行为都正确。

### Tests for User Story 2

- [ ] T017 [P] [US2] 新增 `backend/internal/sqlguard/mysql_validator_test.go`，覆盖 SELECT-only、多语句、DROP/UPDATE/INSERT、语法错误、无 LIMIT 自动补齐 `LIMIT 1000`。
- [ ] T018 [P] [US2] 新增 `backend/internal/query/mysql_executor_test.go`，覆盖 `null`、日期/数字、`[]byte`、空结果、超时和错误脱敏。
- [ ] T019 [P] [US2] 扩展 `backend/internal/api/query_handlers_test.go`，确认 `POST /api/v1/dbs/{name}/query` 按 `databaseType=mysql` 走 MySQL validator/executor，并返回 `limitApplied`、`limit` 和 `empty`。

### Implementation for User Story 2

- [ ] T020 [US2] 在 `backend/internal/sqlguard/mysql_validator.go` 完成 MySQL 只读校验、单语句限制和默认 `LIMIT 1000` 处理，并从 `backend/internal/sqlguard/validator.go` 分发。
- [ ] T021 [US2] 在 `backend/internal/query/mysql_executor.go` 完成 MySQL 执行器、结果扫描、类型转换和 `QueryResult` 组装，并从 `backend/internal/query/executor.go` 分发。
- [ ] T022 [US2] 在 `backend/internal/api/query_handlers.go`、`backend/internal/api/interfaces.go` 中把连接记录的 `databaseType` 传递给 SQL Guard 和 query executor，避免手写 SQL 仍按 PostgreSQL 固定处理。

**Checkpoint**：MySQL 手写 SQL 查询链路已经可独立验证。

---

## Phase 5：User Story 3 - 基于 MySQL metadata 生成 SQL (Priority: P1)

**Goal**：用户可以基于 MySQL metadata 输入自然语言，生成 MySQL 方言 SQL 草稿并继续执行。

**Independent Test**：对 `interview_db` 输入中文 prompt，确认能生成符合 MySQL 方言的 SQL 草稿；再执行该 SQL，确认结果可展示。

### Tests for User Story 3

- [ ] T023 [P] [US3] 更新 `backend/internal/llm/openai_test.go`，并新增 `backend/internal/llm/prompt_builder_test.go`，覆盖 MySQL prompt、缺失 key、非法 JSON 和非法 SQL。
- [ ] T024 [P] [US3] 扩展 `backend/internal/api/natural_query_handlers_test.go`，覆盖 `prompt` / `promt`、metadata 缺失、MySQL 草稿持久化和校验失败返回。

### Implementation for User Story 3

- [ ] T025 [US3] 在 `backend/internal/llm/prompt_builder.go` 中注入 `metadata.databaseType=mysql` 和 MySQL metadata 上下文，要求模型输出 MySQL SELECT SQL。
- [ ] T026 [US3] 在 `backend/internal/api/natural_query_handlers.go`、`backend/internal/llm/sql_draft.go`、`backend/internal/dbstore/store.go` 中完成 MySQL 草稿生成、SQL Guard 校验和持久化。
- [ ] T027 [US3] 更新 `frontend/src/components/NaturalLanguagePanel.vue`、`frontend/src/App.vue`，让自然语言生成结果可预览、可填回 QueryEditor，并在 MySQL 场景下保持同样的交互。

**Checkpoint**：MySQL 自然语言生成 SQL 已经可以独立走通。

---

## Phase 6：User Story 4 - 保留 PostgreSQL 现有能力 (Priority: P2)

**Goal**：新增 MySQL 后，原有 PostgreSQL 连接、查询和自然语言生成能力不能退化。

**Independent Test**：继续对现有 PostgreSQL 数据库执行连接、查询和自然语言生成 SQL，结果与当前版本一致。

### Tests for User Story 4

- [ ] T028 [P] [US4] 为 `backend/internal/api/dbs_handlers_test.go`、`backend/internal/api/query_handlers_test.go`、`backend/internal/api/natural_query_handlers_test.go` 增加 PostgreSQL 回归用例，确认 `databaseType` 缺省仍走 PostgreSQL。
- [ ] T029 [P] [US4] 扩展 `frontend/tests/unit/queryFlow.test.ts`、`frontend/tests/unit/metadataExplorer.test.ts`，确认默认 PostgreSQL、MySQL 选择和现有状态渲染都正常。

### Implementation for User Story 4

- [ ] T030 [US4] 在 `backend/internal/metadata/collector.go`、`backend/internal/query/executor.go`、`backend/internal/sqlguard/validator.go` 中确认缺省 `databaseType` 仍回落到 PostgreSQL，旧连接不受 MySQL 分支影响。
- [ ] T031 [US4] 在 `frontend/src/components/ConnectionPanel.vue` 和 `frontend/src/App.vue` 中保持 PostgreSQL 作为默认选项，并确保 MySQL 只是新增选项，不改变现有连接流程。

**Checkpoint**：PostgreSQL 回归路径和默认行为都没有被 MySQL 影响。

---

## Phase 7：验证与收尾

**Purpose**：完成后端、前端、浏览器和文档验证，确保实现与 spec / plan / contracts 一致。

- [ ] T032 [P] [Verify] 运行 `cd backend && go test ./...`，修复 MySQL、回归和日志相关失败。
- [ ] T033 [P] [Verify] 运行 `cd frontend && npm.cmd test`、`npm.cmd run build:dev`，修复类型、组件和状态渲染问题。
- [ ] T034 [Verify] 启动后端和前端，用 Google Chrome 插件验证 `interview_db` 的新增、metadata 浏览、手写 SQL、自然语言生成 SQL 和结果展示。
- [ ] T035 [Docs] 如实现细节与 `specs/003-mysql-db-query/quickstart.md`、`plan.md` 或 `contracts/openapi.yaml` 有偏差，更新对应文档并保留日志字段、安全约束和 `databaseType` 契约。

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**：无依赖，可直接开始。
- **Phase 2**：依赖 Phase 1，建立 MySQL 分发和日志基础。
- **Phase 3 / Phase 4 / Phase 5 / Phase 6**：依赖 Phase 2。
- **Phase 7**：依赖前述实现完成。

### User Story Dependencies

- **User Story 1 (P1)**：可以在 Foundation 完成后独立实现。
- **User Story 2 (P1)**：依赖 User Story 1 的连接记录和数据库类型分发，但可以独立测试。
- **User Story 3 (P1)**：依赖 User Story 1 的 metadata 和 User Story 2 的 SQL Guard，但可以独立测试。
- **User Story 4 (P2)**：依赖共享分发完成，用于防回归。

### Within Each User Story

- 先写测试，再实现。
- 先打通后端，再收尾前端。
- 任何影响共享分发的改动，都要先确认 PostgreSQL 默认行为不变。

### Parallel Opportunities

- `T001` / `T002` / `T003` 可以并行。
- `T004` / `T005` / `T006` / `T007` / `T008` / `T009` / `T010` 可以并行拆开推进。
- `T011` / `T012` 可以并行。
- `T017` / `T018` / `T019` 可以并行。
- `T023` / `T024` 可以并行。
- `T028` / `T029` 可以并行。

---

## Implementation Strategy

### MVP First

1. 完成 Phase 1 和 Phase 2。
2. 完成 User Story 1。
3. **STOP and VALIDATE**：确认 `interview_db` 能被添加且 metadata 可见。
4. 完成 User Story 2。
5. 再完成 User Story 3。
6. 最后补齐 PostgreSQL 回归和 Phase 7 验证。

### Incremental Delivery

1. 先交付 MySQL 连接与 metadata。
2. 再交付 MySQL 手写 SQL 查询。
3. 再交付 MySQL 自然语言生成 SQL。
4. 最后确认 PostgreSQL 没有退化。

## Notes

- `promt` 仅作为兼容输入别名，正式类型和响应统一使用 `prompt`。
- 前端不得暴露完整 DB URL、密码或 LLM key。
- 所有后端返回给前端的 JSON 字段名必须保持 camelCase。
- MySQL 和 PostgreSQL 的 SQL 校验不能混用错误方言规则。
