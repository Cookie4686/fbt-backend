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
	case errors.SessionExpire:
		w.WriteHeader(http.StatusUnauthorized)
	case errors.RegistrationSessionExpire:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Error(r.URL.Path, zap.String("Method", r.Method), zap.Error(err))
}

func SendJson(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
