package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Addr          string
	SQLitePath    string
	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string
	OpenAIWireAPI string
}

func Load() (Config, error) {
	loadEnvFiles(".env")
	codex := loadCodexConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	addr := getenv("DB_QUERRY_ADDR", ":8080")
	sqlitePath := getenv("DB_QUERRY_SQLITE_PATH", filepath.Join(home, ".db_querry", "db_querrt.db"))

	return Config{
		Addr:          addr,
		SQLitePath:    sqlitePath,
		OpenAIAPIKey:  openAIAPIKey(),
		OpenAIBaseURL: getenv("LLM_BASE_URL", getenv("OPENAI_BASE_URL", defaultString(codex.BaseURL, "https://api.openai.com/v1"))),
		OpenAIModel:   getenv("LLM_MODEL", getenv("OPENAI_MODEL", defaultString(codex.Model, "gpt-4.1-mini"))),
		OpenAIWireAPI: getenv("LLM_WIRE_API", getenv("OPENAI_WIRE_API", defaultString(codex.WireAPI, "chat"))),
	}, nil
}

func openAIAPIKey() string {
	if key := os.Getenv("LLM_API_KEY"); key != "" {
		return key
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return key
	}
	return os.Getenv("openai_api_key")
}

func loadEnvFiles(name string) {
	loadDotEnv(name)
	loadDotEnv(filepath.Join("backend", name))
	loadDotEnv(filepath.Join("env", name))
	loadDotEnv(filepath.Join("backend", "env", name))
	for _, dir := range parentDirs() {
		loadDotEnv(filepath.Join(dir, name))
		loadDotEnv(filepath.Join(dir, "backend", name))
		loadDotEnv(filepath.Join(dir, "env", name))
		loadDotEnv(filepath.Join(dir, "backend", "env", name))
	}
}

func parentDirs() []string {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	var dirs []string
	for {
		dirs = append(dirs, cwd)
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
	return dirs
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

type codexProviderConfig struct {
	BaseURL string
	Model   string
	WireAPI string
}

func loadCodexConfig() codexProviderConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		return codexProviderConfig{}
	}
	file, err := os.Open(filepath.Join(home, ".codex", "config.toml"))
	if err != nil {
		return codexProviderConfig{}
	}
	defer file.Close()

	var provider string
	var inProvider bool
	var result codexProviderConfig
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inProvider = provider != "" && line == "[model_providers."+provider+"]"
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "model_provider" && provider == "" {
			provider = value
			continue
		}
		if key == "model" && result.Model == "" {
			result.Model = value
			continue
		}
		if inProvider {
			switch key {
			case "base_url":
				result.BaseURL = value
			case "wire_api":
				result.WireAPI = value
			}
		}
	}
	return result
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
