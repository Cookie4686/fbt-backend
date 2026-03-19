package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (api *api) useMiddlewareLogger() {
	logger := api.logger
	router := api.router

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("",
				zap.String("Path", r.URL.Path),
				zap.String("Method", r.Method),
			)
			next.ServeHTTP(w, r)
		})
	})
}
