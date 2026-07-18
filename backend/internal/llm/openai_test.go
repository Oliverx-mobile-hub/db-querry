package llm

import "testing"

func TestParseDraft(t *testing.T) {
	draft, err := parseDraft(`{"sql":"SELECT 1","explanation":"ok","referencedObjects":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	if draft.SQL != "SELECT 1" {
		t.Fatalf("unexpected draft: %+v", draft)
	}
}
