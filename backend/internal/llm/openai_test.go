package llm

import (
	"strings"
	"testing"

	"db-querry/backend/internal/api"
)

func TestParseDraft(t *testing.T) {
	draft, err := parseDraft(`{"sql":"SELECT 1","explanation":"ok","referencedObjects":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	if draft.SQL != "SELECT 1" {
		t.Fatalf("unexpected draft: %+v", draft)
	}
}

func TestBuildPromptIncludesDatabaseType(t *testing.T) {
	prompt := buildPrompt("查询用户", api.MetadataDocument{DatabaseType: api.DatabaseTypeMySQL})
	if !strings.Contains(prompt, "Database type: mysql") {
		t.Fatalf("expected mysql database type in prompt, got %s", prompt)
	}
}
