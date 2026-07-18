package api

import "testing"

func TestValidDBName(t *testing.T) {
	if !validDBName("local_1") {
		t.Fatalf("expected valid db name")
	}
	if validDBName("../secret") {
		t.Fatalf("expected invalid db name")
	}
}
