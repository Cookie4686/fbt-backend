package api

import (
	"fbt/backend/internal/dependency"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func useMiddlewareLogger(d *dependency.Dependency, router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := NewLoggingResponseWriter(w)
			next.ServeHTTP(w, r)

			statusCode := wrappedWriter.statusCode

			d.Logger.Info("",
				zap.String("Path", r.URL.Path),
				zap.String("Method", r.Method),
				zap.Int("Status", statusCode),
			)
		})
	})
}
