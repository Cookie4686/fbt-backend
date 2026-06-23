package auth

import (
	"fbt/backend/internal/domain/auth/features"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

func RegisterService(server *grpc.Server, d *util.Dependency) *grpc.Server {
	f := features.NewFeatures(d)

	f.RegisterCredentials(server)
	f.RegisterMFA(server)
	f.RegisterOAuth(server)
	f.RegisterSession(server)
	// f.RegisterUser(server)

	return server
}
