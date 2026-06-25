package interceptor

import (
	"context"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor interface {
	Auth(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (any, error)
	Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}

type interceptor struct {
	*util.Dependency

	service service.Service
}

func NewMiddleware(d *util.Dependency, service service.Service) Interceptor {
	return Interceptor(&interceptor{Dependency: d, service: service})
}

func IsPrivateMethod(fullMethod string) bool {
	switch fullMethod {
	case "/credentials.Credentials/Register", "/credentials.Credentials/Login":
		return false
	case "/oauth.OAuth/Register", "/oauth.OAuth/Login":
		return false
	default:
		return true
	}
}

func (s *interceptor) Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m, err := handler(ctx, req)

	if err != nil {
		s.Logger.Error(info.FullMethod, zap.String("error", err.Error()))
	} else {
		s.Logger.Info(info.FullMethod)
	}

	return m, err
}

func (s *interceptor) Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Intercept Private route
	if IsPrivateMethod(info.FullMethod) {
		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.MissingMetadata
		}

		authorization := md["session_id"]
		if len(authorization) < 1 {
			return nil, errors.Unauthorized
		}
		token := authorization[0]

		auth, err := s.service.Validate(ctx, token)
		if err == errors.NotFound {
			return nil, errors.Unauthorized
		} else if err != nil {
			return nil, errors.DBError
		}

		ctx = NewAuthContext(ctx, auth)
	}

	m, err := handler(ctx, req)

	return m, err
}
