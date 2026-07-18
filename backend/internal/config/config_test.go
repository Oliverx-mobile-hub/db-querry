package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsSharedBackendEnvFromRepoRoot(t *testing.T) {
	root := t.TempDir()
	writeTestEnv(t, filepath.Join(root, "backend", "env"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("LLM_API_KEY", "")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")
	t.Setenv("LLM_WIRE_API", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_MODEL", "")
	t.Setenv("OPENAI_WIRE_API", "")
	t.Setenv("DB_QUERRY_ADDR", "")
	t.Setenv("DB_QUERRY_SQLITE_PATH", "")
	chdir(t, root)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	assertLoadedTestEnv(t, cfg)
}

func TestLoadReadsSharedBackendEnvFromBackendDir(t *testing.T) {
	root := t.TempDir()
	backendDir := filepath.Join(root, "backend")
	writeTestEnv(t, filepath.Join(backendDir, "env"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("LLM_API_KEY", "")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")
	t.Setenv("LLM_WIRE_API", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_MODEL", "")
	t.Setenv("OPENAI_WIRE_API", "")
	t.Setenv("DB_QUERRY_ADDR", "")
	t.Setenv("DB_QUERRY_SQLITE_PATH", "")
	chdir(t, backendDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	assertLoadedTestEnv(t, cfg)
}

func TestLoadReadsSharedBackendEnvFromNestedDir(t *testing.T) {
	root := t.TempDir()
	writeTestEnv(t, filepath.Join(root, "backend", "env"))
	nestedDir := filepath.Join(root, "frontend", "src")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("LLM_API_KEY", "")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")
	t.Setenv("LLM_WIRE_API", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_MODEL", "")
	t.Setenv("OPENAI_WIRE_API", "")
	t.Setenv("DB_QUERRY_ADDR", "")
	t.Setenv("DB_QUERRY_SQLITE_PATH", "")
	chdir(t, nestedDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	assertLoadedTestEnv(t, cfg)
}

func writeTestEnv(t *testing.T, envDir string) {
	t.Helper()
	if err := os.MkdirAll(envDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	shared := []byte("DB_QUERRY_ADDR=:18080\nLLM_BASE_URL=https://moacode.org/team/v1\nLLM_MODEL=gpt-test\nLLM_WIRE_API=responses\nLLM_API_KEY=team-test-key\n")
	if err := os.WriteFile(filepath.Join(envDir, ".env"), shared, 0o600); err != nil {
		t.Fatalf("WriteFile(.env) error = %v", err)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(original)
	})
}

func assertLoadedTestEnv(t *testing.T, cfg Config) {
	t.Helper()
	if cfg.OpenAIAPIKey == "" {
		t.Fatal("OpenAIAPIKey is empty")
	}
	if cfg.OpenAIBaseURL != "https://moacode.org/team/v1" {
		t.Fatalf("OpenAIBaseURL = %q", cfg.OpenAIBaseURL)
	}
	if cfg.OpenAIModel != "gpt-test" {
		t.Fatalf("OpenAIModel = %q", cfg.OpenAIModel)
	}
	if cfg.OpenAIWireAPI != "responses" {
		t.Fatalf("OpenAIWireAPI = %q", cfg.OpenAIWireAPI)
	}
	if cfg.Addr != ":18080" {
		t.Fatalf("Addr = %q", cfg.Addr)
	}
}
