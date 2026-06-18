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
	case errors.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Error(r.URL.Path, zap.String("Method", r.Method), zap.Error(err))
}

func SendJson[T any](w http.ResponseWriter, p *Response[T]) error {
	if p != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(p.StatusCode)
		json.NewEncoder(w).Encode(p.Payload)
	}
	return nil
}

type Response[T any] struct {
	StatusCode int
	Payload    *T
}

func ExtractPayload[T any](r *http.Request) (*T, error) {
	var body T
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		return &body, err
	}
	return &body, nil
}
