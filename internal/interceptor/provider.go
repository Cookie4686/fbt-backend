// Package interceptor for server interceptor (middleware)
package interceptor

import (
	"context"

	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/gen/proto/go/health/v1/healthv1connect"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"go.uber.org/zap"
)

type InterceptorProvider interface {
	Auth() connect.UnaryInterceptorFunc
	Logging() connect.UnaryInterceptorFunc
	Validator() connect.Interceptor
}

type interceptorProvider struct {
	*util.Dependency

	service service.Service
}

func NewInterceptorProvider(d *util.Dependency, service service.Service) InterceptorProvider {
	return InterceptorProvider(&interceptorProvider{Dependency: d, service: service})
}

func IsPrivateMethod(fullMethod string) bool {
	switch fullMethod {
	case
		authv1connect.CredentialServiceRegisterProcedure,
		authv1connect.CredentialServiceLoginProcedure,
		authv1connect.OAuthServiceRegisterProcedure,
		authv1connect.OAuthServiceLoginProcedure,
		authv1connect.WebAuthnServiceGetUserPasskeyProcedure,
		authv1connect.WebAuthnServiceUpdatePasskeyCounterProcedure,
		healthv1connect.HealthServiceStatusProcedure:
		return false
	default:
		return true
	}
}

func VerifyTokenContext(ctx context.Context, service service.Service) (context.Context, error) {
	token, err := FromTokenContext(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := service.Validate(ctx, token)
	if err == errors.NotFound {
		return nil, errors.Unauthorized
	} else if err != nil {
		return nil, errors.DBError
	}

	return NewAuthContext(ctx, auth), err
}

func (s *interceptorProvider) Logging() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			res, err := next(ctx, req)
			if err != nil {
				s.Logger.Error(req.Spec().Procedure, zap.String("error", err.Error()))

				if errors.IsPublicError(err) {
					return res, err
				}
				// Hide Error Information
				return res, errors.InternalError
			}

			s.Logger.Info(req.Spec().Procedure)

			return res, nil
		}
	}
}

func (s *interceptorProvider) Auth() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			// Intercept Private route
			if IsPrivateMethod(req.Spec().Procedure) {
				authContext, err := VerifyTokenContext(ctx, s.service)
				if err != nil {
					return nil, err
				}

				return next(authContext, req)
			}

			return next(ctx, req)
		}
	}
}

func (s *interceptorProvider) Validator() connect.Interceptor {
	return validate.NewInterceptor()
}
