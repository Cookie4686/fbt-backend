package bookkeeping

import (
	"fbt/backend/internal/domain/bookkeeping/features"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

func RegisterService(server *grpc.Server, d *util.Dependency) *grpc.Server {
	f := features.NewFeatures(d)

	f.RegisterAccount(server)
	f.RegisterTransaction(server)

	return server
}
