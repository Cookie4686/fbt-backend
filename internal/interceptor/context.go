package interceptor

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"
)

func NewAuthContext(ctx context.Context, auth *model.Auth) context.Context {
	return context.WithValue(ctx, "auth", auth)
}

func FromAuthContext(ctx context.Context) (*model.Auth, error) {
	a, ok := ctx.Value("auth").(*model.Auth)
	if !ok {
		return nil, errors.MissingMetadata
	}
	return a, nil
}
