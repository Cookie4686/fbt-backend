package util

import (
	"encoding/json"
	"fbt/backend/internal/errors"
	"net/http"

	"github.com/go-playground/validator/v10"
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
		if err := json.NewEncoder(w).Encode(p.Payload); err != nil {
			return errors.BadRequest
		}
	}
	return nil
}

type Response[T any] struct {
	StatusCode int
	Payload    *T
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func ExtractPayload[T any](r *http.Request) (*T, error) {
	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	if err := validate.Struct(body); err != nil {
		return nil, err
	}
	return &body, nil
}
