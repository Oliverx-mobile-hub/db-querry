# Feature Specification: MySQL 数据库查询支持

**Feature Branch**: `003-mysql-db-query`

**Created**: 2026-07-22

**Status**: Draft

**Input**: User description: "我要为 db-querry 添加一个新功能：添加MySQL db查询。支持MySQL db支持的测试用例，然后运行测试，如果后端测试OK，那么打开后端和前端，使用插件Google Chrome测试前端，确保MySQL db的基本功能能实现；添加数据库interview_db，生成sql，查询数据库interview_db，并显示结果，支持自然语言生成sql，查询数据库interview_db，并显示结果。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 添加 MySQL 数据库连接并采集 metadata (Priority: P1)

用户可以在现有数据库列表中新增一个 MySQL 数据库连接，例如 `interview_db`，系统会测试连接、采集 MySQL 的 schema/table/view/column metadata，并把结果展示在前端。

**Why this priority**: 没有连接和 metadata，就无法执行查询，也无法做自然语言生成 SQL。

**Independent Test**: 添加一个有效的 MySQL 连接后，刷新数据库列表，确认连接状态、metadata 状态和 metadata 浏览树都正常显示。

**Acceptance Scenarios**:

1. **Given** 用户输入一个有效的 MySQL 连接和名称 `interview_db`，**When** 保存连接，**Then** 系统应连接数据库、采集 metadata、持久化到本地 SQLite，并在前端显示可浏览的表和视图信息。
2. **Given** MySQL 连接无效或数据库不可达，**When** 保存连接，**Then** 系统应返回可理解的错误，并将连接状态标记为失败，不显示伪造的 metadata。

---

### User Story 2 - 执行 MySQL 手写 SQL 并显示结果 (Priority: P1)

用户可以选中 `interview_db`，输入 MySQL 的只读 `SELECT` 语句并执行，前端以表格显示结果。

**Why this priority**: 手写 SQL 是查询工具的核心能力，也是验证 MySQL 支持是否真正可用的关键路径。

**Independent Test**: 针对 `interview_db` 执行一个合法 `SELECT`，确认结果表格正确显示；再执行非 `SELECT` 或多语句 SQL，确认被拒绝。

**Acceptance Scenarios**:

1. **Given** 用户已选择 `interview_db`，**When** 输入合法的 MySQL `SELECT` 语句并执行，**Then** 系统应返回查询结果、行数、耗时和 LIMIT 状态，并在前端表格中展示。
2. **Given** 用户输入 `INSERT`、`UPDATE`、`DELETE`、`DROP`、多语句或其他非只读 SQL，**When** 执行，**Then** 后端应拒绝并返回清晰错误，不执行数据库写操作。
3. **Given** 用户输入的 `SELECT` 没有显式 `LIMIT`，**When** 执行，**Then** 后端应默认应用 `LIMIT 1000`，并在响应中说明已应用限制。

---

### User Story 3 - 基于 MySQL metadata 生成 SQL 并执行 (Priority: P1)

用户可以针对 `interview_db` 输入自然语言描述，系统基于 MySQL metadata 生成 SQL 草稿，用户确认后执行并查看结果。

**Why this priority**: 自然语言生成 SQL 是现有产品的重要能力，MySQL 支持必须覆盖同一条路径。

**Independent Test**: 对 `interview_db` 输入一个中文 prompt，确认能生成符合 MySQL 方言的 SQL 草稿；再执行该 SQL，确认结果可展示。

**Acceptance Scenarios**:

1. **Given** `interview_db` 已成功采集 metadata，**When** 用户输入自然语言提示并请求生成 SQL，**Then** 系统应返回针对 MySQL 的 SQL 草稿和引用对象信息。
2. **Given** 用户将生成的 SQL 发送执行，**When** SQL 校验通过，**Then** 系统应返回查询结果并在前端显示。
3. **Given** metadata 缺失或数据库离线，**When** 用户请求自然语言生成 SQL，**Then** 系统应明确提示先完成连接或 metadata 采集。

---

### User Story 4 - 保留 PostgreSQL 现有能力 (Priority: P2)

新增 MySQL 后，原有 PostgreSQL 连接、查询和自然语言生成 SQL 能力不能退化。

**Why this priority**: 这是回归保护，避免新增 MySQL 支持时破坏现有可用功能。

**Independent Test**: 继续对现有 PostgreSQL 数据库执行连接、查询和自然语言生成 SQL，结果与当前版本一致。

**Acceptance Scenarios**:

1. **Given** 系统已有 PostgreSQL 数据库连接，**When** 新增 MySQL 功能上线，**Then** PostgreSQL 的连接、metadata、查询和 NL2SQL 行为仍然可用。
2. **Given** 用户在 MySQL 和 PostgreSQL 间切换，**When** 执行查询或生成 SQL，**Then** 系统应按当前选中的数据库类型使用对应的 metadata 和 SQL 方言规则。

### Edge Cases

- MySQL 表名、列名或 schema 名称包含保留关键字时，应能正确显示并在生成 SQL 时保留合法引用方式。
- MySQL 字段类型包含 `decimal`、`datetime`、`timestamp`、`json`、`blob`、`text`、`enum` 等类型时，应正确展示 metadata 和查询结果。
- 数据库离线、超时或认证失败时，不应显示缓存的伪装为最新的 metadata。
- 空结果和 loading 状态下，前端的查询与自然语言按钮状态应保持明确。
- MySQL 和 PostgreSQL 的 SQL 校验不能混用错误方言规则。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 支持新增 MySQL 数据库连接，并把连接信息持久化到本地 SQLite。
- **FR-002**: 系统 MUST 能够连接 MySQL 数据库、测试可达性，并采集 schema/table/view/column metadata。
- **FR-003**: 系统 MUST 在前端展示 MySQL 数据库的连接状态和 metadata 浏览结果。
- **FR-004**: 系统 MUST 支持对 MySQL 数据库执行只读 `SELECT` 查询，并返回结构化查询结果。
- **FR-005**: 系统 MUST 对 MySQL 查询执行 SQL 校验，拒绝非 `SELECT`、多语句和其他破坏性语句，并在缺少 `LIMIT` 时默认应用 `LIMIT 1000`。
- **FR-006**: 系统 MUST 支持基于 MySQL metadata 的自然语言生成 SQL，并返回可供用户确认和执行的 SQL 草稿。
- **FR-007**: 系统 MUST 保持现有 PostgreSQL 功能不退化，原有数据库连接、metadata、查询和自然语言流程仍可使用。
- **FR-008**: 系统 MUST 继续使用统一 JSON envelope、camelCase 字段名和现有错误展示风格。
- **FR-009**: 系统 MUST 不在前端暴露数据库凭据、密钥或完整连接字符串。
- **FR-010**: 系统 MUST 保持前端 UI 风格与当前 MotherDuck-inspired / brutalist dashboard 风格一致。
- **FR-011**: 实现完成后 MUST 通过后端测试；后端测试通过后，才能用 Chrome 验证前端 MySQL 基本流程可用。

### Key Entities *(include if feature involves data)*

- **DatabaseConnection**: 一个可连接的数据库实例，包含名称、类型、连接信息、连接状态和 metadata 状态。
- **MetadataDocument**: 某个数据库采集到的结构化元数据，包含 schema、table、view、column 和类型信息。
- **QueryResult**: 用户执行 SQL 后得到的结果，包含 columns、rows、rowCount、durationMs、limitApplied、limit、empty。
- **GeneratedSqlDraft**: 自然语言生成的 SQL 草稿，包含 prompt、sql、explanation、referencedObjects 和 validation。
- **DatabaseType**: 数据库类型枚举，至少包含 `postgres` 和 `mysql`。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 用户可以成功添加名为 `interview_db` 的 MySQL 数据库连接，并在前端看到 metadata。
- **SC-002**: 用户可以对 `interview_db` 执行至少一个合法的 MySQL `SELECT` 查询，并在前端看到结果表格。
- **SC-003**: 用户可以对 `interview_db` 输入自然语言并生成 MySQL SQL 草稿，再执行并看到结果。
- **SC-004**: 非 `SELECT` 或多语句 SQL 在 MySQL 场景下必须被拒绝，不能执行。
- **SC-005**: PostgreSQL 原有功能在新增 MySQL 后仍能通过回归测试。
- **SC-006**: 后端测试通过后，前端 MySQL 基本流程可以在 Chrome 中完成一次完整验证。

## Assumptions

- 本功能 v1 支持 MySQL 与现有 PostgreSQL 并存，不移除现有数据库类型。
- MySQL 连接格式、鉴权方式和数据库名由后续 plan 统一细化，但都必须走后端。
- 自然语言生成 SQL 仍然复用现有的 LLM 流程，只是 metadata 和 SQL 方言切换到 MySQL 场景。
- 前端导出、验证和可视化风格保持现有实现，不重新设计整套 UI。
