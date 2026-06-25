package api

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/domain/auth/middleware"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/domain/bookkeeping"
	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewGRPCServer(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config) *grpc.Server {
	d := util.NewDependency(logger, db, cfg)
	m := middleware.NewMiddleware(d, service.NewService(d))
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpc.UnaryServerInterceptor(m.LoggerInterceptor),
		grpc.UnaryServerInterceptor(m.AuthInterceptor),
	))

	auth.RegisterService(s, d)
	bookkeeping.RegisterService(s, d)

	return s
}
