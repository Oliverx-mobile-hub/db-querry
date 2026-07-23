package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"db-querry/backend/internal/api"
	"github.com/jackc/pgx/v5"
)

type Store interface {
	GetConnection(ctx context.Context, name string) (api.DBConnectionRecord, error)
}

type Validator interface {
	Validate(databaseType api.DatabaseType, sql string) api.SQLValidationResult
}

type Executor struct {
	store     Store
	validator Validator
}

func NewExecutor(store Store, validator Validator) Executor {
	return Executor{store: store, validator: validator}
}

func (e Executor) Execute(ctx context.Context, dbName string, sqlText string) (api.QueryResult, error) {
	record, err := e.store.GetConnection(ctx, dbName)
	if err != nil {
		return api.QueryResult{}, err
	}
	databaseType := api.NormalizeDatabaseType(record.DatabaseType)
	validation := e.validator.Validate(databaseType, sqlText)
	if !validation.Executable {
		return api.QueryResult{Validation: &validation}, fmt.Errorf("sql validation failed")
	}

	switch databaseType {
	case api.DatabaseTypeMySQL:
		return e.executeMySQL(ctx, record.URL, validation)
	default:
		return e.executePostgres(ctx, record.URL, validation)
	}
}

func (e Executor) executePostgres(ctx context.Context, rawURL string, validation api.SQLValidationResult) (api.QueryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, rawURL)
	if err != nil {
		return api.QueryResult{}, err
	}
	defer conn.Close(ctx)

	started := time.Now()
	rows, err := conn.Query(ctx, validation.NormalizedSQL)
	if err != nil {
		return api.QueryResult{}, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	columns := make([]api.QueryColumn, 0, len(fields))
	for _, field := range fields {
		columns = append(columns, api.QueryColumn{Name: field.Name, DataType: fmt.Sprintf("oid:%d", field.DataTypeOID)})
	}

	resultRows := []map[string]any{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return api.QueryResult{}, err
		}
		row := map[string]any{}
		for i, value := range values {
			row[columns[i].Name] = jsonValue(value)
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

func jsonValue(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case sql.NullString:
		if v.Valid {
			return v.String
		}
		return nil
	default:
		return v
	}
}
