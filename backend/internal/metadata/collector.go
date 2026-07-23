package metadata

import (
	"context"
	"fmt"
	"sort"
	"time"

	"db-querry/backend/internal/api"
	"github.com/jackc/pgx/v5"
)

type Collector struct{}

func NewCollector() Collector { return Collector{} }

func (c Collector) Collect(ctx context.Context, databaseType api.DatabaseType, url string) (api.MetadataDocument, []string, error) {
	switch api.NormalizeDatabaseType(databaseType) {
	case api.DatabaseTypeMySQL:
		return c.collectMySQL(ctx, url)
	case api.DatabaseTypePostgres:
		return c.collectPostgres(ctx, url)
	default:
		return api.MetadataDocument{}, nil, fmt.Errorf("unsupported database type: %s", databaseType)
	}
}

func (Collector) collectPostgres(ctx context.Context, url string) (api.MetadataDocument, []string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return api.MetadataDocument{}, nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, metadataSQL)
	if err != nil {
		return api.MetadataDocument{}, nil, err
	}
	defer rows.Close()

	type key struct{ schema, name, objectType string }
	objects := map[key]*api.MetadataObject{}
	for rows.Next() {
		var schemaName, objectName, objectType, objectComment string
		var columnName, dataType, columnComment string
		var nullable, primaryKey bool
		var ordinal int
		if err := rows.Scan(&schemaName, &objectName, &objectType, &objectComment, &columnName, &dataType, &nullable, &primaryKey, &ordinal, &columnComment); err != nil {
			return api.MetadataDocument{}, nil, err
		}
		k := key{schema: schemaName, name: objectName, objectType: objectType}
		object := objects[k]
		if object == nil {
			object = &api.MetadataObject{Schema: schemaName, Name: objectName, Type: objectType, Comment: objectComment}
			objects[k] = object
		}
		object.Columns = append(object.Columns, api.MetadataColumn{
			Name:       columnName,
			DataType:   dataType,
			Nullable:   nullable,
			PrimaryKey: primaryKey,
			Ordinal:    ordinal,
			Comment:    columnComment,
		})
	}
	if err := rows.Err(); err != nil {
		return api.MetadataDocument{}, nil, err
	}

	schemaMap := map[string][]api.MetadataObject{}
	for _, object := range objects {
		sort.Slice(object.Columns, func(i, j int) bool { return object.Columns[i].Ordinal < object.Columns[j].Ordinal })
		schemaMap[object.Schema] = append(schemaMap[object.Schema], *object)
	}

	schemas := make([]api.MetadataSchema, 0, len(schemaMap))
	for schemaName, schemaObjects := range schemaMap {
		sort.Slice(schemaObjects, func(i, j int) bool {
			if schemaObjects[i].Type == schemaObjects[j].Type {
				return schemaObjects[i].Name < schemaObjects[j].Name
			}
			return schemaObjects[i].Type < schemaObjects[j].Type
		})
		schemas = append(schemas, api.MetadataSchema{Name: schemaName, Objects: schemaObjects})
	}
	sort.Slice(schemas, func(i, j int) bool { return schemas[i].Name < schemas[j].Name })

	return api.MetadataDocument{DatabaseType: api.DatabaseTypePostgres, Schemas: schemas}, nil, nil
}

const metadataSQL = `
WITH pk_columns AS (
  SELECT
    ns.nspname AS schema_name,
    cls.relname AS table_name,
    att.attname AS column_name
  FROM pg_index idx
  JOIN pg_class cls ON cls.oid = idx.indrelid
  JOIN pg_namespace ns ON ns.oid = cls.relnamespace
  JOIN pg_attribute att ON att.attrelid = cls.oid AND att.attnum = ANY(idx.indkey)
  WHERE idx.indisprimary
)
SELECT
  n.nspname AS schema_name,
  c.relname AS object_name,
  CASE c.relkind WHEN 'v' THEN 'view' WHEN 'm' THEN 'view' ELSE 'table' END AS object_type,
  COALESCE(obj_description(c.oid), '') AS object_comment,
  a.attname AS column_name,
  format_type(a.atttypid, a.atttypmod) AS data_type,
  NOT a.attnotnull AS nullable,
  pk.column_name IS NOT NULL AS primary_key,
  a.attnum AS ordinal,
  COALESCE(col_description(c.oid, a.attnum), '') AS column_comment
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
JOIN pg_attribute a ON a.attrelid = c.oid
LEFT JOIN pk_columns pk ON pk.schema_name = n.nspname AND pk.table_name = c.relname AND pk.column_name = a.attname
WHERE c.relkind IN ('r', 'p', 'v', 'm')
  AND a.attnum > 0
  AND NOT a.attisdropped
  AND n.nspname NOT IN ('pg_catalog', 'information_schema')
ORDER BY n.nspname, c.relname, a.attnum
`
