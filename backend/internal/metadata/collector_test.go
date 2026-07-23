package metadata

import (
	"strings"
	"testing"
)

func TestCollectorConstructor(t *testing.T) {
	_ = NewCollector()
}

func TestMySQLMetadataSQLTargetsCurrentDatabase(t *testing.T) {
	if !strings.Contains(mysqlMetadataSQL, "information_schema.columns") {
		t.Fatalf("expected MySQL metadata SQL to use information_schema.columns")
	}
	if !strings.Contains(mysqlMetadataSQL, "c.table_schema = DATABASE()") {
		t.Fatalf("expected MySQL metadata SQL to target current database")
	}
}
