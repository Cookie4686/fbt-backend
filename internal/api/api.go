package api

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/domain/bookkeeping"
	"fbt/backend/internal/domain/health"
	"fbt/backend/internal/util"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewAPIHandler(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	d := util.NewDependency(logger, db, cfg)

	useMiddlewareLogger(d, r)

	health.Routes(d, r)
	m := auth.Routes(d, r)

	privateRouter := r.NewRoute().Subrouter()
	privateRouter.Use(m.Auth)

	bookkeeping.Routes(d, privateRouter)

	return r
}
