package main

import (
	"net/http"

	"db-querry/backend/internal/api"
	"db-querry/backend/internal/config"
	"db-querry/backend/internal/dbstore"
	"db-querry/backend/internal/llm"
	"db-querry/backend/internal/logging"
	"db-querry/backend/internal/metadata"
	"db-querry/backend/internal/pgconn"
	"db-querry/backend/internal/query"
	"db-querry/backend/internal/sqlguard"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logging.Fatalf("load config: %v", err)
	}

	store, err := dbstore.Open(cfg.SQLitePath)
	if err != nil {
		logging.Fatalf("open sqlite store: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(); err != nil {
		logging.Fatalf("migrate sqlite store: %v", err)
	}

	logging.Infof("llm config loaded: base_url=%s model=%s wire_api=%s key_loaded=%t", cfg.OpenAIBaseURL, cfg.OpenAIModel, cfg.OpenAIWireAPI, cfg.OpenAIAPIKey != "")

	deps := api.Dependencies{
		Config:    cfg,
		Store:     store,
		Connector: pgconn.NewConnector(),
		Metadata:  metadata.NewCollector(),
		SQLGuard:  sqlguard.NewValidator(),
		Query:     query.NewExecutor(store, sqlguard.NewValidator()),
		LLM:       llm.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAIBaseURL, cfg.OpenAIModel, cfg.OpenAIWireAPI),
	}

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: api.NewRouter(deps),
	}

	logging.Infof("db-querry backend listening on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Fatalf("server failed: %v", err)
	}
}
