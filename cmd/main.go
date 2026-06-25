package main

import (
	"context"
	"fbt/backend/internal/api"
	"fbt/backend/internal/config"
	"fbt/backend/internal/util"
	"net"

	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig(".env")
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

	s := api.NewGRPCServer(logger, db, cfg)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.API.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
