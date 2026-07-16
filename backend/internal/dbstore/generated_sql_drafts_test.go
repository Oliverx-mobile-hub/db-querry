package dbstore

import (
	"context"
	"testing"

	"db-querry/backend/internal/api"
)

func TestInsertGeneratedSQLDraft(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.Migrate(); err != nil {
		t.Fatal(err)
	}
	if err := store.UpsertConnection(ctx, api.DBConnectionRecord{Name: "local", DatabaseType: "postgres", URL: "postgres://u:p@h/db", DisplayDSN: "postgres://u@h/db", MetadataStatus: api.MetadataReady}); err != nil {
		t.Fatal(err)
	}
	if err := store.InsertGeneratedSQLDraft(ctx, "local", api.GeneratedSQLDraft{Prompt: "p", SQL: "SELECT 1", Validation: api.SQLValidationResult{Valid: true}}); err != nil {
		t.Fatal(err)
	}
}

