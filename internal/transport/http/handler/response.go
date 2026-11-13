package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/recco-app/recco/internal/domain"
)

// APIResponse standardizes JSON response payloads.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// WriteJSON sends a formatted JSON response.
func WriteJSON(w http.ResponseWriter, status int, data interface{}, meta ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	}
	if len(meta) > 0 {
		resp.Meta = meta[0]
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// WriteError inspects error type and renders appropriate HTTP status and message.
func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		msg = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		msg = err.Error()
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		msg = err.Error()
	case errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
		msg = err.Error()
	default:
		msg = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   msg,
	})
}
