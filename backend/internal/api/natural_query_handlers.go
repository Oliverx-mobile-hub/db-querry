package api

import (
	"errors"
	"net/http"
)

type naturalQueryRequest struct {
	Prompt string `json:"prompt"`
	Promt  string `json:"promt"`
}

func (h *Handler) naturalQuery(w http.ResponseWriter, r *http.Request, name string) {
	var request naturalQueryRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalidRequest", "请求体必须包含 prompt", nil)
		return
	}
	prompt := request.Prompt
	if prompt == "" {
		prompt = request.Promt
	}
	if prompt == "" {
		writeError(w, http.StatusBadRequest, "invalidRequest", "请求体必须包含 prompt", nil)
		return
	}
	metadataDoc, _, err := h.deps.Store.GetLatestMetadataSnapshot(r.Context(), name)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusBadRequest, "metadataCollectionFailed", "请先采集 metadata", nil)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internalError", "读取 metadata 失败", nil)
		return
	}
	draft, err := h.deps.LLM.GenerateSQL(r.Context(), prompt, metadataDoc)
	if err != nil {
		writeError(w, http.StatusBadRequest, "llmUnavailable", "生成 SQL 失败", nil)
		return
	}
	draft.Prompt = prompt
	draft.Validation = h.deps.SQLGuard.Validate(draft.SQL)
	_ = h.deps.Store.InsertGeneratedSQLDraft(r.Context(), name, draft)
	writeOK(w, draft)
}
