# Feature Specification: 导出当前查询结果为 CSV/JSON

**Feature Branch**: `002-export-query-results`

**Created**: 2026-07-20

**Status**: Draft

**Input**: User description: "我要为 db-querry 添加一个新功能：导出当前查询结果为 CSV/JSON。Results 面板增加 Export CSV、Export JSON 两个按钮。无结果、空结果、查询中时按钮禁用。JSON 导出当前 QueryResult 或仅导出 rows/columns，需要先定。CSV 正确处理逗号、双引号、换行、null、日期/数字。文件名包含数据库名和时间，例如 local-query-20260720-143000.csv。单元测试覆盖 CSV 转义和按钮可用状态。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 导出当前查询结果为 CSV (Priority: P1)

用户执行 SQL 查询并看到结果表格后，可以点击 Results 面板中的 `Export CSV` 按钮，将当前查询结果下载为 CSV 文件，用于在 Excel、表格工具或其他数据处理工具中继续分析。

**Why this priority**: CSV 是查询结果最常见的外部交换格式，直接提升查询工具的可用性。

**Independent Test**: 执行一个包含多列、多行、特殊字符和空值的查询，点击 `Export CSV`，验证下载文件名、列顺序、单元格转义和内容完整性。

**Acceptance Scenarios**:

1. **Given** 用户已经成功执行查询并看到非空结果，**When** 用户点击 `Export CSV`，**Then** 系统下载一个 `.csv` 文件，文件包含当前结果的表头和所有当前行。
2. **Given** 查询结果中包含逗号、双引号、换行、null、数字和日期值，**When** 用户导出 CSV，**Then** CSV 内容符合标准转义规则，导入表格工具后列和值不串列、不丢失。
3. **Given** 当前数据库名为 `local`，**When** 用户导出 CSV，**Then** 下载文件名包含数据库名和时间戳，例如 `local-query-20260720-143000.csv`。

---

### User Story 2 - 导出当前查询结果为 JSON (Priority: P1)

用户执行 SQL 查询并看到结果表格后，可以点击 Results 面板中的 `Export JSON` 按钮，将当前查询结果下载为 JSON 文件，用于调试、归档或传递给其他程序。

**Why this priority**: JSON 是当前 API 和前端状态的原生结构，适合保留查询结果的结构化上下文。

**Independent Test**: 执行一个查询，点击 `Export JSON`，验证下载文件是合法 JSON，且包含当前查询结果所需结构。

**Acceptance Scenarios**:

1. **Given** 用户已经成功执行查询并看到非空结果，**When** 用户点击 `Export JSON`，**Then** 系统下载一个 `.json` 文件。
2. **Given** 当前查询结果包含 columns、rows、rowCount、durationMs、limitApplied、limit、empty 和 validation，**When** 用户导出 JSON，**Then** JSON 文件导出完整当前 `QueryResult`，而不是只导出 rows。
3. **Given** 当前数据库名为 `local`，**When** 用户导出 JSON，**Then** 下载文件名包含数据库名和时间戳，例如 `local-query-20260720-143000.json`。

---

### User Story 3 - 禁用不可用导出操作 (Priority: P2)

用户在没有可导出结果时不应点击导出按钮，避免下载空文件或产生误解。

**Why this priority**: 明确的按钮状态能减少无效操作和错误提示。

**Independent Test**: 分别在未查询、查询中、查询返回空结果时查看 Results 面板，验证导出按钮禁用。

**Acceptance Scenarios**:

1. **Given** 用户还没有执行任何查询，**When** Results 面板显示初始状态，**Then** `Export CSV` 和 `Export JSON` 按钮禁用。
2. **Given** 查询正在执行中，**When** Results 面板处于 loading 状态，**Then** `Export CSV` 和 `Export JSON` 按钮禁用。
3. **Given** 查询成功但结果为空，**When** Results 面板显示空结果，**Then** `Export CSV` 和 `Export JSON` 按钮禁用。

### Edge Cases

- CSV 字段包含逗号时必须被双引号包裹。
- CSV 字段包含双引号时必须将双引号转义为两个双引号。
- CSV 字段包含换行时必须保留换行并正确包裹字段。
- null 或 undefined 值在 CSV 中导出为空字段，在 JSON 中保留为 `null`。
- 数字、布尔值、日期或时间字符串在 CSV 中按展示值导出，在 JSON 中保持当前 `QueryResult.rows` 中的 JSON 值。
- 列顺序必须与 `QueryResult.columns` 顺序一致，不能依赖对象 key 枚举顺序。
- 文件名中的数据库名必须进行安全化处理，避免空格、斜杠或特殊字符导致下载异常。
- 如果浏览器阻止下载，系统应保留当前结果状态，不影响继续查询。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Results 面板必须显示 `Export CSV` 和 `Export JSON` 两个导出按钮。
- **FR-002**: 当没有查询结果、查询结果为空、或查询正在执行时，系统必须禁用两个导出按钮。
- **FR-003**: `Export CSV` 必须导出当前前端持有的查询结果，不重新执行 SQL。
- **FR-004**: `Export JSON` 必须导出完整当前 `QueryResult`，包括 columns、rows、rowCount、durationMs、limitApplied、limit、empty 和 validation。
- **FR-005**: CSV 导出必须按照 `QueryResult.columns` 的顺序生成表头和每一行。
- **FR-006**: CSV 导出必须正确处理逗号、双引号、换行、null、日期、数字和布尔值。
- **FR-007**: 导出文件名必须包含数据库名、`query` 标识和本地时间戳，格式为 `{dbName}-query-{YYYYMMDD-HHmmss}.{csv|json}`。
- **FR-008**: 导出功能不得向后端发送数据库凭据、模型密钥或额外高权限配置。
- **FR-009**: 导出当前结果不应改变现有 SQL 校验、查询执行、自然语言生成 SQL 或 metadata 展示行为。
- **FR-010**: 前端测试必须覆盖 CSV 转义规则、按钮禁用状态和文件名格式。
- **FR-011**: 所有新增前端类型、props 和 emits 必须保持 TypeScript 严格类型。

### Key Entities *(include if feature involves data)*

- **QueryResult**: 当前查询返回给前端的结构化结果，是 JSON 导出的完整对象，也是 CSV 导出的数据来源。
- **Export Format**: 用户选择的导出格式，支持 `csv` 和 `json`。
- **Export Filename**: 由数据库名、固定 `query` 标识、本地时间戳和扩展名组成的下载文件名。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 用户在查询成功并返回非空结果后，可以在 2 次点击内下载 CSV 或 JSON 文件。
- **SC-002**: CSV 导出的特殊字符测试样例全部通过，包括逗号、双引号、换行和 null。
- **SC-003**: JSON 导出的文件可以被标准 JSON parser 成功解析，并保留完整当前 `QueryResult`。
- **SC-004**: 未查询、查询中、空结果三种状态下，导出按钮均不可点击。
- **SC-005**: 新增单元测试在本地测试命令中通过，不破坏现有查询和结果展示测试。

## Assumptions

- v1 只导出当前前端已经持有的结果，不新增后端导出 API，也不重新执行 SQL。
- 当前查询结果受现有后端 limit 规则约束；导出不会绕过默认 LIMIT 或只读 SQL 安全策略。
- JSON 导出选择完整 `QueryResult`，因为它保留 columns、rows、统计信息和 SQL validation 上下文。
- CSV 导出只表达表格数据，不包含 durationMs、validation 或 limit 等元信息。
- 导出按钮放在 Results 面板标题区域，风格沿用现有 Element Plus 和 MotherDuck-inspired / brutalist dashboard 视觉体系。
- 文件时间戳使用浏览器本地时间。
