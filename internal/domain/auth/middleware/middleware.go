package middleware

import (
	"context"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"net/http"
)

type Middleware interface {
	Auth(next http.Handler) http.Handler
}

type middleware struct {
	*util.Dependency

	service service.Service
}

func NewMiddleware(d *util.Dependency, service service.Service) Middleware {
	return Middleware(&middleware{Dependency: d, service: service})
}

func (s *middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		sessionID, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			util.SendError(s.Logger, w, r, err)
			return
		}

		auth, err := s.service.Validate(ctx, sessionID.Value)
		if err == errors.NotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			util.SendError(s.Logger, w, r, err)
			return
		}

		requestWithAuth := r.WithContext(context.WithValue(ctx, "auth", auth))

		next.ServeHTTP(w, requestWithAuth)
	})
}
