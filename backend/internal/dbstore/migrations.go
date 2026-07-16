package dbstore

func (s *Store) Migrate() error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS db_connections (
  name TEXT PRIMARY KEY,
  database_type TEXT NOT NULL,
  url TEXT NOT NULL,
  display_dsn TEXT NOT NULL,
  metadata_status TEXT NOT NULL,
  metadata_error TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
		`CREATE TABLE IF NOT EXISTS metadata_snapshots (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  db_name TEXT NOT NULL,
  metadata_json TEXT NOT NULL,
  object_count INTEGER NOT NULL,
  warning_json TEXT NOT NULL DEFAULT '[]',
  created_at TEXT NOT NULL,
  FOREIGN KEY (db_name) REFERENCES db_connections(name) ON DELETE CASCADE
)`,
		`CREATE INDEX IF NOT EXISTS idx_metadata_snapshots_db_name_created_at ON metadata_snapshots(db_name, created_at DESC)`,
		`CREATE TABLE IF NOT EXISTS generated_sql_drafts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  db_name TEXT NOT NULL,
  prompt TEXT NOT NULL,
  sql TEXT NOT NULL,
  explanation TEXT NOT NULL DEFAULT '',
  referenced_objects_json TEXT NOT NULL DEFAULT '[]',
  validation_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (db_name) REFERENCES db_connections(name) ON DELETE CASCADE
)`,
		`CREATE INDEX IF NOT EXISTS idx_generated_sql_drafts_db_name_created_at ON generated_sql_drafts(db_name, created_at DESC)`,
	}
	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

