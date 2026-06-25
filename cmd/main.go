package main

import (
	"context"
	"fbt/backend/internal/server"
	"fbt/backend/internal/util"
	"net"

	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	d, err := util.NewDependency(ctx, ".env")
	if err != nil {
		log.Fatalf("failed to create dependency: %v", err)
	}

	s := server.NewServer(d)

	if lis, err := net.Listen("tcp", fmt.Sprintf(":%d", d.CFG.API.PORT)); err != nil {
		d.Logger.Fatal("failed to listen", zap.Error(err))
	} else if err := s.Serve(lis); err != nil {
		d.Logger.Fatal("failed to serve", zap.Error(err))
	} else {
		d.Logger.Info(
			"server started",
			zap.String("network", lis.Addr().Network()),
			zap.String("address", lis.Addr().String()),
		)
	}
}
