package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var dbNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

type putDBRequest struct {
	URL          string       `json:"url"`
	DatabaseType DatabaseType `json:"databaseType"`
}

func (h *Handler) handleDBs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/dbs" {
		writeError(w, http.StatusNotFound, "notFound", "接口不存在", nil)
		return
	}
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "methodNotAllowed", "请求方法不支持", nil)
		return
	}
	records, err := h.deps.Store.ListConnections(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "读取数据库连接失败", nil)
		return
	}
	dbs := make([]DBSummary, 0, len(records))
	for _, record := range records {
		dbs = append(dbs, h.summaryWithHealth(r.Context(), record))
	}
	writeOK(w, map[string]any{"dbs": dbs})
}

func (h *Handler) handleDBByName(w http.ResponseWriter, r *http.Request) {
	name, suffix, ok := splitDBPath(r.URL.Path)
	if !ok || !validDBName(name) {
		writeError(w, http.StatusBadRequest, "invalidRequest", "数据库名称不合法", nil)
		return
	}

	switch {
	case suffix == "" && r.Method == http.MethodPut:
		h.putDB(w, r, name)
	case suffix == "" && r.Method == http.MethodGet:
		h.getDBMetadata(w, r, name)
	case suffix == "" && r.Method == http.MethodDelete:
		h.deleteDB(w, r, name)
	case suffix == "/query" && r.Method == http.MethodPost:
		h.queryDB(w, r, name)
	case suffix == "/query/natural" && r.Method == http.MethodPost:
		h.naturalQuery(w, r, name)
	default:
		writeError(w, http.StatusNotFound, "notFound", "接口不存在", nil)
	}
}

func (h *Handler) putDB(w http.ResponseWriter, r *http.Request, name string) {
	var request putDBRequest
	if err := decodeJSON(r, &request); err != nil || strings.TrimSpace(request.URL) == "" {
		writeError(w, http.StatusBadRequest, "invalidRequest", "请求体必须包含 url", nil)
		return
	}
	request.URL = strings.TrimSpace(request.URL)
	if _, err := url.ParseRequestURI(request.URL); err != nil {
		writeError(w, http.StatusBadRequest, "invalidRequest", "数据库 URL 格式不正确", nil)
		return
	}
	databaseType, err := resolveDatabaseType(request.DatabaseType, request.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalidRequest", err.Error(), nil)
		return
	}

	displayDSN := sanitizeDSN(request.URL)
	now := time.Now().Format(time.RFC3339)
	record := DBConnectionRecord{
		Name:           name,
		DatabaseType:   databaseType,
		URL:            request.URL,
		DisplayDSN:     displayDSN,
		MetadataStatus: MetadataPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := h.deps.Store.UpsertConnection(r.Context(), record); err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "保存数据库连接失败", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	if err := h.deps.Connector.Test(ctx, databaseType, request.URL); err != nil {
		_ = h.deps.Store.UpdateMetadataStatus(r.Context(), name, MetadataFailed, "数据库连接失败")
		writeError(w, http.StatusBadRequest, "dbConnectionFailed", "数据库连接失败，请检查 URL、网络和凭据", nil)
		return
	}

	metadataDoc, warnings, err := h.deps.Metadata.Collect(ctx, databaseType, request.URL)
	if err != nil {
		_ = h.deps.Store.UpdateMetadataStatus(r.Context(), name, MetadataFailed, "metadata 采集失败")
		writeError(w, http.StatusBadRequest, "metadataCollectionFailed", "metadata 采集失败", nil)
		return
	}
	if err := h.deps.Store.InsertMetadataSnapshot(r.Context(), name, metadataDoc, countObjects(metadataDoc), warnings); err != nil {
		_ = h.deps.Store.UpdateMetadataStatus(r.Context(), name, MetadataFailed, "metadata 保存失败")
		writeError(w, http.StatusInternalServerError, "internalError", "metadata 保存失败", nil)
		return
	}
	if err := h.deps.Store.UpdateMetadataStatus(r.Context(), name, MetadataReady, ""); err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "更新 metadata 状态失败", nil)
		return
	}

	updated, _ := h.deps.Store.GetConnection(r.Context(), name)
	writeOK(w, map[string]any{"db": h.summaryWithHealth(r.Context(), updated)})
}

func (h *Handler) deleteDB(w http.ResponseWriter, r *http.Request, name string) {
	if err := h.deps.Store.DeleteConnection(r.Context(), name); errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "dbNotFound", "数据库连接不存在", nil)
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "删除数据库连接失败", nil)
		return
	}
	writeOK(w, map[string]any{"deleted": true, "name": name})
}

func (h *Handler) getDBMetadata(w http.ResponseWriter, r *http.Request, name string) {
	record, err := h.deps.Store.GetConnection(r.Context(), name)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "dbNotFound", "数据库连接不存在", nil)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "读取数据库连接失败", nil)
		return
	}
	connectionStatus := h.connectionStatus(r.Context(), record)
	if connectionStatus == "offline" {
		writeOK(w, map[string]any{
			"name":              name,
			"metadataStatus":    MetadataFailed,
			"connectionStatus":  connectionStatus,
			"metadataUpdatedAt": record.MetadataUpdatedAt,
			"metadata":          MetadataDocument{DatabaseType: NormalizeDatabaseType(record.DatabaseType), Schemas: []MetadataSchema{}},
		})
		return
	}
	metadataDoc, updatedAt, err := h.deps.Store.GetLatestMetadataSnapshot(r.Context(), name)
	if errors.Is(err, ErrNotFound) {
		metadataDoc = MetadataDocument{DatabaseType: NormalizeDatabaseType(record.DatabaseType), Schemas: []MetadataSchema{}}
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "读取 metadata 失败", nil)
		return
	}
	if updatedAt != nil {
		record.MetadataUpdatedAt = updatedAt
	}
	if metadataDoc.DatabaseType == "" {
		metadataDoc.DatabaseType = NormalizeDatabaseType(record.DatabaseType)
	}
	writeOK(w, map[string]any{
		"name":              name,
		"metadataStatus":    record.MetadataStatus,
		"connectionStatus":  connectionStatus,
		"metadataUpdatedAt": record.MetadataUpdatedAt,
		"metadata":          metadataDoc,
	})
}

func validDBName(name string) bool { return dbNamePattern.MatchString(name) }

func sanitizeDSN(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "<invalid>"
	}
	if parsed.User != nil {
		username := parsed.User.Username()
		parsed.User = url.User(username)
	}
	return parsed.String()
}

func countObjects(metadataDoc MetadataDocument) int {
	count := 0
	for _, schema := range metadataDoc.Schemas {
		count += len(schema.Objects)
	}
	return count
}

func (h *Handler) summaryWithHealth(ctx context.Context, record DBConnectionRecord) DBSummary {
	summary := summaryFromRecord(record)
	summary.ConnectionStatus = h.connectionStatus(ctx, record)
	if summary.ConnectionStatus == "offline" {
		summary.MetadataStatus = MetadataFailed
	}
	return summary
}

func (h *Handler) connectionStatus(ctx context.Context, record DBConnectionRecord) string {
	if record.URL == "" {
		return "offline"
	}
	checkCtx, cancel := context.WithTimeout(ctx, 1200*time.Millisecond)
	defer cancel()
	if err := h.deps.Connector.Test(checkCtx, NormalizeDatabaseType(record.DatabaseType), record.URL); err != nil {
		return "offline"
	}
	return "online"
}

func resolveDatabaseType(requested DatabaseType, rawURL string) (DatabaseType, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	inferred := DatabaseTypePostgres
	switch strings.ToLower(parsed.Scheme) {
	case "postgres", "postgresql":
		inferred = DatabaseTypePostgres
	case "mysql", "mysql2":
		inferred = DatabaseTypeMySQL
	default:
		return "", errors.New("数据库 URL scheme 只支持 postgres、postgresql、mysql")
	}
	if requested == "" {
		return inferred, nil
	}
	requested = DatabaseType(strings.ToLower(string(requested)))
	if !SupportedDatabaseType(requested) {
		return "", errors.New("databaseType 只支持 postgres 或 mysql")
	}
	if parsed.Scheme != "" && (strings.HasPrefix(strings.ToLower(parsed.Scheme), "postgres") || strings.HasPrefix(strings.ToLower(parsed.Scheme), "mysql")) && requested != inferred {
		return "", errors.New("databaseType 与数据库 URL scheme 不匹配")
	}
	return requested, nil
}
