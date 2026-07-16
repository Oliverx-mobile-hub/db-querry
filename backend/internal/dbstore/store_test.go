package dbstore

import (
	"context"
	"testing"

	"db-querry/backend/internal/api"
)

func TestConnectionAndMetadataStore(t *testing.T) {
	ctx := context.Background()
	store, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.Migrate(); err != nil {
		t.Fatal(err)
	}

	record := api.DBConnectionRecord{Name: "local", DatabaseType: "postgres", URL: "postgres://user:pass@localhost/db", DisplayDSN: "postgres://user@localhost/db", MetadataStatus: api.MetadataPending}
	if err := store.UpsertConnection(ctx, record); err != nil {
		t.Fatal(err)
	}
	if err := store.InsertMetadataSnapshot(ctx, "local", api.MetadataDocument{DatabaseType: "postgres", Schemas: []api.MetadataSchema{}}, 0, nil); err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateMetadataStatus(ctx, "local", api.MetadataReady, ""); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetConnection(ctx, "local")
	if err != nil {
		t.Fatal(err)
	}
	if got.URL == got.DisplayDSN {
		t.Fatalf("expected sanitized display dsn")
	}
	metadata, _, err := store.GetLatestMetadataSnapshot(ctx, "local")
	if err != nil {
		t.Fatal(err)
	}
	if metadata.DatabaseType != "postgres" {
		t.Fatalf("unexpected metadata: %+v", metadata)
	}
	if err := store.DeleteConnection(ctx, "local"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetConnection(ctx, "local"); err == nil {
		t.Fatalf("expected deleted connection to be missing")
	}
}
