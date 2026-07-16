package dbstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"db-querry/backend/internal/api"
	_ "modernc.org/sqlite"
)
type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &Store{db: db}, nil
}

func OpenMemory() (*Store, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) UpsertConnection(ctx context.Context, record api.DBConnectionRecord) error {
	now := time.Now().Format(time.RFC3339)
	if record.CreatedAt == "" {
		record.CreatedAt = now
	}
	record.UpdatedAt = now
	_, err := s.db.ExecContext(ctx, `
INSERT INTO db_connections (name, database_type, url, display_dsn, metadata_status, metadata_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET
  database_type=excluded.database_type,
  url=excluded.url,
  display_dsn=excluded.display_dsn,
  metadata_status=excluded.metadata_status,
  metadata_error=excluded.metadata_error,
  updated_at=excluded.updated_at
`, record.Name, record.DatabaseType, record.URL, record.DisplayDSN, string(record.MetadataStatus), record.MetadataError, record.CreatedAt, record.UpdatedAt)
	return err
}

func (s *Store) DeleteConnection(ctx context.Context, name string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM db_connections WHERE name = ?`, name)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return api.ErrNotFound
	}
	return nil
}

func (s *Store) ListConnections(ctx context.Context) ([]api.DBConnectionRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT c.name, c.database_type, c.url, c.display_dsn, c.metadata_status, COALESCE(c.metadata_error, ''),
       c.created_at, c.updated_at, MAX(m.created_at) AS metadata_updated_at
FROM db_connections c
LEFT JOIN metadata_snapshots m ON m.db_name = c.name
GROUP BY c.name, c.database_type, c.url, c.display_dsn, c.metadata_status, c.metadata_error, c.created_at, c.updated_at
ORDER BY c.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []api.DBConnectionRecord
	for rows.Next() {
		record, err := scanConnection(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (s *Store) GetConnection(ctx context.Context, name string) (api.DBConnectionRecord, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT c.name, c.database_type, c.url, c.display_dsn, c.metadata_status, COALESCE(c.metadata_error, ''),
       c.created_at, c.updated_at, MAX(m.created_at) AS metadata_updated_at
FROM db_connections c
LEFT JOIN metadata_snapshots m ON m.db_name = c.name
WHERE c.name = ?
GROUP BY c.name, c.database_type, c.url, c.display_dsn, c.metadata_status, c.metadata_error, c.created_at, c.updated_at
`, name)
	record, err := scanConnection(row)
	if err == sql.ErrNoRows {
		return api.DBConnectionRecord{}, api.ErrNotFound
	}
	return record, err
}

func (s *Store) UpdateMetadataStatus(ctx context.Context, name string, status api.DBMetadataStatus, metadataError string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE db_connections SET metadata_status = ?, metadata_error = ?, updated_at = ? WHERE name = ?`, string(status), metadataError, time.Now().Format(time.RFC3339), name)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return api.ErrNotFound
	}
	return nil
}

func (s *Store) InsertMetadataSnapshot(ctx context.Context, dbName string, metadata api.MetadataDocument, objectCount int, warnings []string) error {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	warningsJSON, err := json.Marshal(warnings)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO metadata_snapshots (db_name, metadata_json, object_count, warning_json, created_at) VALUES (?, ?, ?, ?, ?)`, dbName, string(metadataJSON), objectCount, string(warningsJSON), time.Now().Format(time.RFC3339))
	return err
}

func (s *Store) GetLatestMetadataSnapshot(ctx context.Context, dbName string) (api.MetadataDocument, *string, error) {
	var raw string
	var createdAt string
	err := s.db.QueryRowContext(ctx, `SELECT metadata_json, created_at FROM metadata_snapshots WHERE db_name = ? ORDER BY created_at DESC, id DESC LIMIT 1`, dbName).Scan(&raw, &createdAt)
	if err == sql.ErrNoRows {
		return api.MetadataDocument{}, nil, api.ErrNotFound
	}
	if err != nil {
		return api.MetadataDocument{}, nil, err
	}
	var metadata api.MetadataDocument
	if err := json.Unmarshal([]byte(raw), &metadata); err != nil {
		return api.MetadataDocument{}, nil, err
	}
	return metadata, &createdAt, nil
}

func (s *Store) InsertGeneratedSQLDraft(ctx context.Context, dbName string, draft api.GeneratedSQLDraft) error {
	referenced, err := json.Marshal(draft.ReferencedObjects)
	if err != nil {
		return err
	}
	validation, err := json.Marshal(draft.Validation)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO generated_sql_drafts (db_name, prompt, sql, explanation, referenced_objects_json, validation_json, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, dbName, draft.Prompt, draft.SQL, draft.Explanation, string(referenced), string(validation), time.Now().Format(time.RFC3339))
	return err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanConnection(row scanner) (api.DBConnectionRecord, error) {
	var record api.DBConnectionRecord
	var status string
	var metadataUpdatedAt sql.NullString
	err := row.Scan(&record.Name, &record.DatabaseType, &record.URL, &record.DisplayDSN, &status, &record.MetadataError, &record.CreatedAt, &record.UpdatedAt, &metadataUpdatedAt)
	if err != nil {
		return api.DBConnectionRecord{}, err
	}
	record.MetadataStatus = api.DBMetadataStatus(status)
	if metadataUpdatedAt.Valid {
		record.MetadataUpdatedAt = &metadataUpdatedAt.String
	}
	return record, nil
}
