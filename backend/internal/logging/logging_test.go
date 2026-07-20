package logging

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestInfofWritesStructuredStdoutFriendlyLine(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetNow(func() time.Time { return time.Date(2026, 7, 20, 16, 5, 6, 0, time.Local) })
	t.Cleanup(Reset)

	Infof("llm config loaded: base_url=%s key_loaded=%t", "https://api2.codexcn.com/v1", true)

	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, `time="2026-07-20 16:05:06"`) {
		t.Fatalf("log timestamp missing: %s", got)
	}
	if !strings.Contains(got, "level=info") {
		t.Fatalf("log level missing: %s", got)
	}
	if !strings.Contains(got, "base_url=https://api2.codexcn.com/v1") {
		t.Fatalf("log message missing: %s", got)
	}
	if strings.Contains(got, "sk-") {
		t.Fatalf("log should not contain key material: %s", got)
	}
}
