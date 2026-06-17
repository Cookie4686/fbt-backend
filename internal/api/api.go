package api

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger
	db     *pgxpool.Pool
	router *mux.Router
}

func NewAPIHandler(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()
	api := &api{logger, db, router}

	api.useMiddlewareLogger()

	api.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.SendJson(w, http.StatusOK, struct {
			Status string `json:"status"`
		}{Status: "ok"})
	})

	auth.Routes(logger, db, cfg, router)

	return api.router
}
