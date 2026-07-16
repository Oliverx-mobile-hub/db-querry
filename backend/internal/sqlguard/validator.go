package sqlguard

import (
	"fmt"
	"regexp"
	"strings"

	"db-querry/backend/internal/api"
)

type Validator struct{}

func NewValidator() Validator { return Validator{} }

var forbiddenPattern = regexp.MustCompile(`(?i)\b(insert|update|delete|drop|alter|truncate|create|replace|grant|revoke|exec|execute|call|merge|copy)\b`)
var topLevelLimitPattern = regexp.MustCompile(`(?i)\blimit\s+\d+\s*;?\s*$`)
var selectPattern = regexp.MustCompile(`(?is)^\s*(with\b.+?\bselect\b|select\b)`)

func (Validator) Validate(sqlText string) api.SQLValidationResult {
	trimmed := strings.TrimSpace(sqlText)
	limitValue := 1000
	result := api.SQLValidationResult{
		Valid:         false,
		Executable:    false,
		StatementType: "unknown",
		NormalizedSQL: trimmed,
		Limit:         nil,
		Errors:        []api.ValidationError{},
	}

	if trimmed == "" {
		return addError(result, "emptySql", "SQL 不能为空")
	}

	withoutStrings := removeStringLiterals(trimmed)
	if hasMultipleStatements(withoutStrings) {
		return addError(result, "multipleStatements", "SQL 只能包含一条 SELECT 查询")
	}
	if !selectPattern.MatchString(withoutStrings) {
		return addError(result, "notSelect", "SQL 只能是 SELECT 查询")
	}
	result.StatementType = "select"

	if forbiddenPattern.MatchString(withoutStrings) {
		return addError(result, "forbiddenStatement", "SQL 包含禁止的语句或关键字")
	}
	if regexp.MustCompile(`(?i)\binto\b`).MatchString(withoutStrings) {
		return addError(result, "selectIntoForbidden", "不允许 SELECT INTO")
	}

	result.Valid = true
	result.Executable = true
	if !topLevelLimitPattern.MatchString(trimmed) {
		result.LimitApplied = true
		result.Limit = &limitValue
		result.NormalizedSQL = strings.TrimSuffix(trimmed, ";") + fmt.Sprintf(" LIMIT %d", limitValue)
	} else {
		result.LimitApplied = false
		result.NormalizedSQL = strings.TrimSuffix(trimmed, ";")
	}
	return result
}

func hasMultipleStatements(input string) bool {
	trimmed := strings.TrimSpace(input)
	trimmed = strings.TrimSuffix(trimmed, ";")
	return strings.Contains(trimmed, ";")
}

func addError(result api.SQLValidationResult, code, message string) api.SQLValidationResult {
	result.Valid = false
	result.Executable = false
	result.Errors = append(result.Errors, api.ValidationError{Code: code, Message: message})
	return result
}

func removeStringLiterals(input string) string {
	var builder strings.Builder
	inSingle := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '\'' {
			inSingle = !inSingle
			builder.WriteByte(' ')
			continue
		}
		if inSingle {
			builder.WriteByte(' ')
			continue
		}
		builder.WriteByte(ch)
	}
	return builder.String()
}
