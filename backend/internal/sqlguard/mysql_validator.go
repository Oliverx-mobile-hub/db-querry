package sqlguard

import (
	"fmt"
	"regexp"
	"strings"

	"db-querry/backend/internal/api"
	"github.com/xwb1989/sqlparser"
)

var mysqlSelectIntoPattern = regexp.MustCompile(`(?i)\binto\b`)

func (Validator) validateMySQL(sqlText string) api.SQLValidationResult {
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
	if forbiddenPattern.MatchString(withoutStrings) {
		return addError(result, "forbiddenStatement", "SQL 包含禁止的语句或关键字")
	}
	if mysqlSelectIntoPattern.MatchString(withoutStrings) {
		return addError(result, "selectIntoForbidden", "不允许 SELECT INTO")
	}

	statementSQL := strings.TrimSuffix(trimmed, ";")
	statement, err := sqlparser.Parse(statementSQL)
	if err != nil {
		return addError(result, "syntaxError", "SQL 语法不正确")
	}

	selectStmt, ok := statement.(*sqlparser.Select)
	if !ok {
		return addError(result, "notSelect", "SQL 只能是 SELECT 查询")
	}
	result.StatementType = "select"
	if selectStmt.Lock != "" {
		return addError(result, "lockingSelectForbidden", "不允许 SELECT 锁定语句")
	}

	result.Valid = true
	result.Executable = true
	if selectStmt.Limit == nil {
		selectStmt.SetLimit(&sqlparser.Limit{Rowcount: sqlparser.NewIntVal([]byte(fmt.Sprintf("%d", limitValue)))})
		result.LimitApplied = true
		result.Limit = &limitValue
	} else {
		result.LimitApplied = false
	}
	result.NormalizedSQL = sqlparser.String(selectStmt)
	return result
}
