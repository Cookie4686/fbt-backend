package server

import (
	"fbt/backend/internal/util"
	"net"
	"net/http"

	"go.uber.org/zap"
)

func StartListening(server *http.Server, d *util.Dependency) {
	if lis, err := net.Listen("tcp", server.Addr); err != nil {
		d.Logger.Fatal("failed to listen", zap.Error(err))
	} else if err := server.Serve(lis); err != nil {
		d.Logger.Fatal("failed to serve", zap.Error(err))
	} else {
		d.Logger.Info(
			"server started",
			zap.String("network", lis.Addr().Network()),
			zap.String("address", lis.Addr().String()),
		)
	}
}
