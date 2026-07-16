package dbstore

import "testing"

func TestMigrateCreatesTables(t *testing.T) {
	store, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		t.Fatal(err)
	}

	for _, table := range []string{"db_connections", "metadata_snapshots", "generated_sql_drafts"} {
		var name string
		err := store.db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&name)
		if err != nil {
			t.Fatalf("table %s not created: %v", table, err)
		}
	}
}
