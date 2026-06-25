package middleware

import (
	"context"
	"fbt/backend/internal/errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func IsPrivateRoute(fullMethod string) bool {
	switch fullMethod {
	case "/credentials.Credentials/Register", "/credentials.Credentials/Login":
		return false
	case "/oauth.OAuth/Register", "/oauth.OAuth/Login":
		return false
	default:
		return true
	}
}

func (s *middleware) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Intercept Private route
	if IsPrivateRoute(info.FullMethod) {
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
		ctx = context.WithValue(ctx, "auth", auth)
	}

	m, err := handler(ctx, req)

	return m, err
}

func (s *middleware) LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	m, err := handler(ctx, req)

	if err != nil {
		s.Logger.Error(info.FullMethod, zap.String("error", err.Error()))
	} else {
		s.Logger.Info(info.FullMethod)
	}
	return m, err
}
