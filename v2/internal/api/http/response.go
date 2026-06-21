package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"

	"github.com/omcrgnt/demo/v2/internal/domain"
)

func WriteServiceError(w nethttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		WriteError(w, nethttp.StatusNotFound, err)
	case errors.Is(err, domain.ErrInvalidInput):
		WriteError(w, nethttp.StatusBadRequest, err)
	default:
		WriteError(w, nethttp.StatusInternalServerError, err)
	}
}

func WriteError(w nethttp.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSON(w nethttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
