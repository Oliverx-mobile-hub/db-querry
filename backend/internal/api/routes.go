package api

import (
	"net/http"
	"strings"

	"db-querry/backend/internal/config"
)

type Dependencies struct {
	Config    config.Config
	Store     Store
	Connector Connector
	Metadata  MetadataCollector
	SQLGuard  SQLValidator
	Query     QueryExecutor
	LLM       SQLGenerator
}

type Handler struct {
	deps Dependencies
}

func NewRouter(deps Dependencies) http.Handler {
	h := &Handler{deps: deps}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/dbs", h.handleDBs)
	mux.HandleFunc("/api/v1/dbs/", h.handleDBByName)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func splitDBPath(path string) (name string, suffix string, ok bool) {
	trimmed := strings.TrimPrefix(path, "/api/v1/dbs/")
	if trimmed == "" || strings.HasPrefix(trimmed, "/") {
		return "", "", false
	}
	parts := strings.SplitN(trimmed, "/", 2)
	name = parts[0]
	if len(parts) == 2 {
		suffix = "/" + parts[1]
	}
	return name, suffix, true
}
