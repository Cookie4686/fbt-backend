package test

import (
	"context"

	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/interceptor"
)

func NewAuthContext(ctx context.Context, sessionID string, service service.Service) context.Context {
	auth, err := service.Validate(ctx, sessionID)
	if err != nil {
		return ctx
	}

	return interceptor.NewAuthContext(ctx, auth)
}
