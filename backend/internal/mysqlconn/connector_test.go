package mysqlconn

import (
	"strings"
	"testing"
)

func TestDriverDSNFromURL(t *testing.T) {
	dsn, err := DriverDSN("mysql://root:secret@localhost:3306/interview_db?charset=utf8mb4")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"root:secret@", "tcp(localhost:3306)", "/interview_db?", "parseTime=true", "charset=utf8mb4"} {
		if !strings.Contains(dsn, want) {
			t.Fatalf("expected dsn to contain %q, got %s", want, dsn)
		}
	}
	if strings.Contains(strings.ToLower(dsn), "multistatements=true") {
		t.Fatalf("multi statements must not be enabled: %s", dsn)
	}
}

func TestDriverDSNRejectsMissingDatabase(t *testing.T) {
	if _, err := DriverDSN("mysql://root:secret@localhost:3306"); err == nil {
		t.Fatalf("expected missing database error")
	}
}
