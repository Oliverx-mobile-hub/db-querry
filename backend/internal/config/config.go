package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Addr         string
	SQLitePath   string
	OpenAIAPIKey string
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	addr := getenv("DB_QUERRY_ADDR", ":8080")
	sqlitePath := getenv("DB_QUERRY_SQLITE_PATH", filepath.Join(home, ".db_querry", "db_querrt.db"))

	return Config{
		Addr:         addr,
		SQLitePath:   sqlitePath,
		OpenAIAPIKey: os.Getenv("openai_api_key"),
	}, nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
