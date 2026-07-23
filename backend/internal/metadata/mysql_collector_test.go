package metadata

import "testing"

func TestMySQLMetadataSQLContainsPrimaryKeyJoin(t *testing.T) {
	if mysqlMetadataSQL == "" {
		t.Fatalf("mysql metadata SQL must not be empty")
	}
}
