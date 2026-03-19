package main

import (
	"context"
	"fbt/backend/internal/api"
	"fbt/backend/internal/config"
	"fbt/backend/internal/util"

	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := util.NewLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := util.NewDatabasePool(ctx, cfg)
	if err != nil {
		logger.Fatal("DB Init", zap.Error(err))
	}

	apiHandler := api.NewAPIHandler(logger, db)
	logger.Info("Server Started", zap.String("URL", fmt.Sprintf("localhost:%v", cfg.API.PORT)))
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.API.PORT), apiHandler)
}
