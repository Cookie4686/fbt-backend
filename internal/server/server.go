package server

import (
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/domain/bookkeeping"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

func NewServer(d *util.Dependency) *grpc.Server {
	service := service.NewService(d)

	m := interceptor.NewMiddleware(d, service)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpc.UnaryServerInterceptor(m.Logging),
		grpc.UnaryServerInterceptor(m.Auth),
	))

	auth.RegisterService(server, d)
	bookkeeping.RegisterService(server, d)

	return server
}
