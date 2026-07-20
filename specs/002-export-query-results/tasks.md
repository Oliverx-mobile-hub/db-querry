# Tasks：导出当前查询结果为 CSV/JSON

**Input**：`/specs/002-export-query-results/` 下的 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`quickstart.md`、`contracts/openapi.yaml`

**Tests**：本功能要求覆盖 CSV 转义、JSON 导出、按钮禁用状态和文件名格式；如调整后端启动日志，也要运行后端测试确认不泄露敏感信息。

**Organization**：按 7 个 phase 拆分。v1 以“前端导出当前结果”为主，不新增后端导出 API。

## Format：`[ID] [P?] [Area] Description`

- **[P]**：可并行执行，前提是不同文件且无直接依赖。
- **[Area]**：任务所属区域，例如 `Setup`、`Frontend`、`Utils`、`Tests`、`Backend`、`Verify`。
- 每个任务必须包含明确文件路径。

## Phase 1：项目对齐

**Purpose**：确认本功能沿用现有前后端结构、结果数据类型和导出约定。

- [X] T001 [Setup] 复查 `specs/002-export-query-results/spec.md`、`plan.md`、`data-model.md`、`contracts/openapi.yaml`，确认 v1 不新增后端导出 API，导出基于当前 `QueryResult`。
- [X] T002 [P] [Setup] 复查前端结果展示入口 `frontend/src/components/ResultTable.vue`、`frontend/src/api/types.ts`、`frontend/src/App.vue`，明确导出按钮挂载位置和所需 props。

**Checkpoint**：导出功能的落点、数据源和边界已确认。

---

## Phase 2：基础能力

**Purpose**：建立前端导出工具、时间戳/文件名生成规则和结果状态判断的基础实现。

- [X] T003 [P] [Utils] 新增 `frontend/src/utils/exportResults.ts`，实现文件名安全化、时间戳格式化和导出文件对象构建。
- [X] T004 [P] [Utils] 在 `frontend/src/utils/exportResults.ts` 中实现 CSV 序列化函数，按 `QueryResult.columns` 顺序输出并处理逗号、双引号、换行、null、数字和布尔值。
- [X] T005 [Frontend] 为 `frontend/src/api/types.ts` 补充导出相关类型（如 `ExportFormat`、`ExportFile`、`ExportContext`），保持严格类型。
- [X] T006 [Frontend] 在 `frontend/src/components/ResultTable.vue` 设计导出按钮所需的状态计算输入：`result`、`loading`、`dbName`。

**Checkpoint**：导出通用工具和类型基础可复用，尚未接入 UI。

---

## Phase 3：User Story 1 - Export CSV (Priority: P1)

**Goal**：用户在查询成功并看到非空结果后，可以导出当前结果为 CSV 文件。

**Independent Test**：执行一个包含逗号、双引号、换行、null、数字和日期值的查询，点击 `Export CSV`，验证下载文件名、表头顺序和转义规则。

### Tests for User Story 1

- [X] T007 [P] [Tests] 为 `frontend/src/utils/exportResults.ts` 编写 CSV 序列化测试 `frontend/tests/unit/exportResults.test.ts`，覆盖逗号、双引号、换行、null、数字、日期和列顺序。
- [X] T008 [P] [Tests] 为 `frontend/src/components/ResultTable.vue` 编写按钮状态测试 `frontend/tests/unit/ResultTable.test.ts`，覆盖非空结果可用、无结果禁用、空结果禁用、查询中禁用。

### Implementation for User Story 1

- [X] T009 [Frontend] 在 `frontend/src/components/ResultTable.vue` 增加 `Export CSV` 按钮和下载触发入口。
- [X] T010 [Frontend] 在 `frontend/src/utils/exportResults.ts` 实现 CSV 下载触发逻辑，使用 Blob 和 `URL.createObjectURL` 生成 `.csv` 文件。
- [X] T011 [Frontend] 确保 `frontend/src/components/ResultTable.vue` 在 `result.empty === false && loading === false && result != null` 时才允许导出 CSV。
- [X] T012 [Frontend] 为 CSV 下载文件名接入数据库名和本地时间戳，格式为 `{safeDbName}-query-{YYYYMMDD-HHmmss}.csv`。

**Checkpoint**：CSV 导出可以独立演示，不依赖 JSON 导出。

---

## Phase 4：User Story 2 - Export JSON (Priority: P1)

**Goal**：用户可以将当前完整 `QueryResult` 导出为 JSON 文件，用于调试和结构化保存。

**Independent Test**：执行一个查询后点击 `Export JSON`，验证下载文件是合法 JSON，并包含 columns、rows、rowCount、durationMs、limitApplied、limit、empty 和 validation。

### Tests for User Story 2

- [X] T013 [P] [Tests] 为 `frontend/src/utils/exportResults.ts` 编写 JSON 导出测试 `frontend/tests/unit/exportResults.test.ts`，验证完整 `QueryResult` 序列化。
- [X] T014 [P] [Tests] 为 `frontend/src/components/ResultTable.vue` 补充 JSON 导出按钮状态测试，覆盖与 CSV 相同的禁用条件。

### Implementation for User Story 2

- [X] T015 [Frontend] 在 `frontend/src/components/ResultTable.vue` 增加 `Export JSON` 按钮。
- [X] T016 [Frontend] 在 `frontend/src/utils/exportResults.ts` 实现 JSON 下载触发逻辑，使用 `JSON.stringify(result, null, 2)` 生成 `.json` 文件。
- [X] T017 [Frontend] 确保 JSON 文件名遵循 `{safeDbName}-query-{YYYYMMDD-HHmmss}.json`。

**Checkpoint**：CSV 和 JSON 两种导出格式都可独立使用。

---

## Phase 5：User Story 3 - Disable States & UX (Priority: P2)

**Goal**：在无结果、空结果、查询中状态下，导出按钮明确禁用，避免无效下载操作。

**Independent Test**：分别在未查询、查询中和空结果三种场景下观察 Results 面板，确认两个导出按钮均为禁用状态。

### Tests for User Story 3

- [X] T018 [P] [Tests] 补充 `frontend/tests/unit/ResultTable.test.ts`，验证不同状态下按钮是否可点击和是否渲染正确的禁用态样式/属性。

### Implementation for User Story 3

- [X] T019 [Frontend] 在 `frontend/src/components/ResultTable.vue` 统一计算导出按钮可用性，避免在组件内出现分散判断。
- [X] T020 [Frontend] 在 `frontend/src/components/ResultTable.vue` 保持现有 Results 面板风格，并将导出按钮放在标题区域，符合当前 MotherDuck-inspired / brutalist dashboard 风格。
- [X] T021 [Frontend] 在 `frontend/src/App.vue` 继续向 `ResultTable` 传递当前数据库名和查询 loading 状态，确保按钮禁用条件完整。

**Checkpoint**：导出按钮状态与查询状态一致，UI 不会在空场景下误导用户。

---

## Phase 6：后端日志对齐

**Purpose**：让后端启动日志更适合本地排查和 K8s/云日志采集，且不泄露敏感信息。

- [X] T022 [Backend] 调整 `backend/cmd/server/main.go` 的启动日志输出，保留监听地址和 llm 配置摘要（base_url、model、wire_api、key_loaded），不输出任何 key 或完整连接 URL。
- [X] T023 [P] [Backend] 如需要，抽出轻量日志封装到 `backend/internal/logging/` 或现有启动路径中，保持 stdout 输出和固定时间戳格式。
- [X] T024 [Backend] 更新或新增 `backend/cmd/server/main.go` 相关测试/断言，确认日志配置不会影响 `go test ./...`。

**Checkpoint**：后端启动日志满足可观察性要求，同时保持敏感信息安全。

---

## Phase 7：验证与收尾

**Purpose**：完成跨文件验证，确保导出和日志改动没有破坏现有查询流程。

- [X] T025 [Verify] 运行 `frontend` 单元测试，确认 `exportResults.test.ts` 和 `ResultTable.test.ts` 全部通过。
- [X] T026 [Verify] 运行 `backend` 测试 `go test ./...`，确认日志调整与现有 API 行为兼容。
- [ ] T027 [Verify] 在浏览器中手动验证一次导出 CSV、导出 JSON、按钮禁用状态和文件名格式。
- [X] T028 [P] [Docs] 如实现细节与 quickstart 有偏差，更新 `specs/002-export-query-results/quickstart.md` 的手动验证步骤。

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**：无依赖，可立即开始。
- **Phase 2**：依赖 Phase 1 完成。
- **Phase 3 / Phase 4 / Phase 5**：依赖 Phase 2 完成。
- **Phase 6**：可与前端导出实现并行，但必须在收尾前完成。
- **Phase 7**：依赖前述实现完成。

### User Story Dependencies

- **User Story 1 (CSV)**：可在 Phase 2 完成后独立实现。
- **User Story 2 (JSON)**：可在 Phase 2 完成后独立实现。
- **User Story 3 (Disable States & UX)**：依赖导出按钮已接入，但可与 Story 1/2 并行收尾。

### Within Each User Story

- 先写测试，再实现。
- CSV utility 先于组件按钮接入。
- 组件状态判断先于样式收尾。

### Parallel Opportunities

- `T003` / `T004` / `T005` / `T006` 可并行。
- `T007` / `T008` 可并行。
- `T013` / `T014` 可并行。
- `T022` / `T023` 可并行。

---

## Implementation Strategy

### MVP First

1. 完成 Phase 1 和 Phase 2。
2. 完成 CSV 导出（Phase 3）。
3. **STOP and VALIDATE**：确认 CSV 下载、转义、文件名和按钮状态。
4. 再完成 JSON 导出（Phase 4）。
5. 最后补齐禁用状态、日志对齐和收尾验证。

### Incremental Delivery

1. 先交付 CSV。
2. 再交付 JSON。
3. 再收紧按钮状态和 UX。
4. 最后检查日志和测试。
