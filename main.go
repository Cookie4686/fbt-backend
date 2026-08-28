package main

import (
	"context"
	"log"

	"fbt/backend/internal/server"
	"fbt/backend/internal/util"
)

func main() {
	ctx := context.Background()

	d, err := util.NewDependency(ctx)
	if err != nil {
		log.Fatalf("failed to create dependency: %v", err)
	}

	svr := server.NewServer(d)

	server.StartListening(svr, d)
}
