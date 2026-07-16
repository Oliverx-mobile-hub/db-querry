package api

import (
	"encoding/json"
	"net/http"
)

type responseEnvelope struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *apiError `json:"error"`
}

type apiError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

func writeOK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, responseEnvelope{Success: true, Data: data, Error: nil})
}

func writeError(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	writeJSON(w, status, responseEnvelope{
		Success: false,
		Data:    nil,
		Error: &apiError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body responseEnvelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}
