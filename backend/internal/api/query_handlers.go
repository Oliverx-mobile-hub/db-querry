package api

import (
	"errors"
	"net/http"
)

type queryRequest struct {
	SQL string `json:"sql"`
}

func (h *Handler) queryDB(w http.ResponseWriter, r *http.Request, name string) {
	var request queryRequest
	if err := decodeJSON(r, &request); err != nil || request.SQL == "" {
		writeError(w, http.StatusBadRequest, "invalidRequest", "请求体必须包含 sql", nil)
		return
	}
	result, err := h.deps.Query.Execute(r.Context(), name, request.SQL)
	if err != nil {
		if result.Validation != nil && !result.Validation.Executable {
			writeError(w, http.StatusBadRequest, "sqlValidationFailed", "SQL 校验失败", map[string]any{"validation": result.Validation})
			return
		}
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "dbNotFound", "数据库连接不存在", nil)
			return
		}
		writeError(w, http.StatusBadRequest, "queryExecutionFailed", "查询执行失败", nil)
		return
	}
	writeOK(w, result)
}
