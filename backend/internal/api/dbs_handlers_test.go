package api

import (
	"context"
	"testing"
)

func TestSanitizeDSNRemovesPassword(t *testing.T) {
	got := sanitizeDSN("postgres://postgres:secret@localhost:5432/postgres")
	if got != "postgres://postgres@localhost:5432/postgres" {
		t.Fatalf("unexpected sanitized dsn: %s", got)
	}
}

func TestSanitizeMySQLDSNRemovesPassword(t *testing.T) {
	got := sanitizeDSN("mysql://root:secret@localhost:3306/interview_db")
	if got != "mysql://root@localhost:3306/interview_db" {
		t.Fatalf("unexpected sanitized dsn: %s", got)
	}
}

func TestResolveDatabaseType(t *testing.T) {
	mysqlType, err := resolveDatabaseType("", "mysql://root:secret@localhost:3306/interview_db")
	if err != nil {
		t.Fatal(err)
	}
	if mysqlType != DatabaseTypeMySQL {
		t.Fatalf("expected mysql, got %s", mysqlType)
	}

	postgresType, err := resolveDatabaseType("", "postgres://postgres:secret@localhost:5432/postgres")
	if err != nil {
		t.Fatal(err)
	}
	if postgresType != DatabaseTypePostgres {
		t.Fatalf("expected postgres, got %s", postgresType)
	}

	if _, err := resolveDatabaseType(DatabaseTypeMySQL, "postgres://postgres:secret@localhost:5432/postgres"); err == nil {
		t.Fatalf("expected database type mismatch error")
	}
}

type fakeStore struct{}

func (fakeStore) UpsertConnection(context.Context, DBConnectionRecord) error    { return nil }
func (fakeStore) DeleteConnection(context.Context, string) error                { return nil }
func (fakeStore) ListConnections(context.Context) ([]DBConnectionRecord, error) { return nil, nil }
func (fakeStore) GetConnection(context.Context, string) (DBConnectionRecord, error) {
	return DBConnectionRecord{}, nil
}
func (fakeStore) UpdateMetadataStatus(context.Context, string, DBMetadataStatus, string) error {
	return nil
}
func (fakeStore) InsertMetadataSnapshot(context.Context, string, MetadataDocument, int, []string) error {
	return nil
}
func (fakeStore) GetLatestMetadataSnapshot(context.Context, string) (MetadataDocument, *string, error) {
	return MetadataDocument{}, nil, nil
}
func (fakeStore) InsertGeneratedSQLDraft(context.Context, string, GeneratedSQLDraft) error {
	return nil
}
func (fakeStore) Close() error { return nil }
