package query

import "testing"

func TestJSONValueConvertsBytes(t *testing.T) {
	got := jsonValue([]byte("abc"))
	if got != "abc" {
		t.Fatalf("unexpected value: %#v", got)
	}
}

