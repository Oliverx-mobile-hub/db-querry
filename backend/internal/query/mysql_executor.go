package query

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"db-querry/backend/internal/api"
	"db-querry/backend/internal/mysqlconn"
)

func (e Executor) executeMySQL(ctx context.Context, rawURL string, validation api.SQLValidationResult) (api.QueryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	db, err := mysqlconn.Open(rawURL)
	if err != nil {
		return api.QueryResult{}, err
	}
	defer db.Close()

	started := time.Now()
	rows, err := db.QueryContext(ctx, validation.NormalizedSQL)
	if err != nil {
		return api.QueryResult{}, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return api.QueryResult{}, err
	}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return api.QueryResult{}, err
	}
	columns := make([]api.QueryColumn, 0, len(columnNames))
	for i, name := range columnNames {
		dataType := ""
		if i < len(columnTypes) {
			dataType = columnTypes[i].DatabaseTypeName()
		}
		columns = append(columns, api.QueryColumn{Name: name, DataType: dataType})
	}

	resultRows := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columnNames))
		dest := make([]any, len(columnNames))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return api.QueryResult{}, err
		}
		row := map[string]any{}
		for i, value := range values {
			row[columns[i].Name] = mysqlJSONValue(value, columns[i].DataType)
		}
		resultRows = append(resultRows, row)
	}
	if err := rows.Err(); err != nil {
		return api.QueryResult{}, err
	}

	return api.QueryResult{
		Columns:      columns,
		Rows:         resultRows,
		RowCount:     len(resultRows),
		DurationMs:   time.Since(started).Milliseconds(),
		LimitApplied: validation.LimitApplied,
		Limit:        validation.Limit,
		Empty:        len(resultRows) == 0,
		Validation:   &validation,
	}, nil
}

func mysqlJSONValue(value any, dataType string) any {
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		if isBinaryType(dataType) {
			return base64.StdEncoding.EncodeToString(v)
		}
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return v
	}
}

func isBinaryType(dataType string) bool {
	normalized := strings.ToLower(dataType)
	return strings.Contains(normalized, "blob") ||
		strings.Contains(normalized, "binary") ||
		normalized == "bit"
}
