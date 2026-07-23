package metadata

import (
	"context"
	"sort"
	"time"

	"db-querry/backend/internal/api"
	"db-querry/backend/internal/mysqlconn"
)

func (Collector) collectMySQL(ctx context.Context, rawURL string) (api.MetadataDocument, []string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	db, err := mysqlconn.Open(rawURL)
	if err != nil {
		return api.MetadataDocument{}, nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, mysqlMetadataSQL)
	if err != nil {
		return api.MetadataDocument{}, nil, err
	}
	defer rows.Close()

	type key struct{ schema, name, objectType string }
	objects := map[key]*api.MetadataObject{}
	for rows.Next() {
		var schemaName, objectName, objectType, objectComment string
		var columnName, dataType, nullableText, columnComment string
		var primaryKey, ordinal int
		if err := rows.Scan(&schemaName, &objectName, &objectType, &objectComment, &columnName, &dataType, &nullableText, &primaryKey, &ordinal, &columnComment); err != nil {
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
			Nullable:   nullableText == "YES",
			PrimaryKey: primaryKey == 1,
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

	return api.MetadataDocument{DatabaseType: api.DatabaseTypeMySQL, Schemas: schemas}, nil, nil
}

const mysqlMetadataSQL = `
SELECT
  c.table_schema AS schema_name,
  c.table_name AS object_name,
  CASE WHEN t.table_type = 'VIEW' THEN 'view' ELSE 'table' END AS object_type,
  COALESCE(t.table_comment, '') AS object_comment,
  c.column_name,
  c.column_type AS data_type,
  c.is_nullable,
  CASE WHEN pk.column_name IS NULL THEN 0 ELSE 1 END AS primary_key,
  c.ordinal_position,
  COALESCE(c.column_comment, '') AS column_comment
FROM information_schema.columns c
JOIN information_schema.tables t
  ON t.table_schema = c.table_schema
 AND t.table_name = c.table_name
LEFT JOIN (
  SELECT kcu.table_schema, kcu.table_name, kcu.column_name
  FROM information_schema.table_constraints tc
  JOIN information_schema.key_column_usage kcu
    ON kcu.constraint_schema = tc.constraint_schema
   AND kcu.constraint_name = tc.constraint_name
   AND kcu.table_schema = tc.table_schema
   AND kcu.table_name = tc.table_name
  WHERE tc.constraint_type = 'PRIMARY KEY'
) pk
  ON pk.table_schema = c.table_schema
 AND pk.table_name = c.table_name
 AND pk.column_name = c.column_name
WHERE c.table_schema = DATABASE()
  AND t.table_type IN ('BASE TABLE', 'VIEW')
ORDER BY c.table_schema, c.table_name, c.ordinal_position
`
