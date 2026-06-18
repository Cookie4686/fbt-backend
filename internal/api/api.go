package api

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/domain/health"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewAPIHandler(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()
	dependency := dependency.NewDependency(logger, db, cfg)

	useMiddlewareLogger(dependency, router)

	health.Routes(dependency, router)
	auth.Routes(dependency, router)

	return router
}
