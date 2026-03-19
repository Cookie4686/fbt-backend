package util

import (
	"encoding/json"
	"fbt/backend/internal/errors"
	"net/http"

	"go.uber.org/zap"
)

func SendError(logger *zap.Logger, w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case errors.NotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Error(r.URL.Path, zap.String("Method", r.Method), zap.Error(err))
}

func SendJson(w http.ResponseWriter, statusCode int, v any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
