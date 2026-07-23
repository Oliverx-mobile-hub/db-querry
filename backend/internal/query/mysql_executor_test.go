package query

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestMySQLJSONValueConvertsTextBytes(t *testing.T) {
	got := mysqlJSONValue([]byte("abc"), "varchar")
	if got != "abc" {
		t.Fatalf("unexpected value: %#v", got)
	}
}

func TestMySQLJSONValueConvertsBinaryBytes(t *testing.T) {
	got := mysqlJSONValue([]byte{0x01, 0x02}, "blob")
	if got != base64.StdEncoding.EncodeToString([]byte{0x01, 0x02}) {
		t.Fatalf("unexpected value: %#v", got)
	}
}

func TestMySQLJSONValueFormatsTime(t *testing.T) {
	when := time.Date(2026, 7, 22, 10, 0, 0, 0, time.UTC)
	got := mysqlJSONValue(when, "datetime")
	if got != "2026-07-22T10:00:00Z" {
		t.Fatalf("unexpected value: %#v", got)
	}
}
