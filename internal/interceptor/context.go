package interceptor

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"

	"connectrpc.com/connect"
)

var sessionIDKey = "session_id"

func NewTokenContext(ctx context.Context, sessionId string) context.Context {
	ctx, callInfo := connect.NewClientContext(ctx)
	callInfo.RequestHeader().Set(sessionIDKey, sessionId)

	return ctx
}

func FromTokenContext(ctx context.Context) (token string, err error) {
	if callInfo, ok := connect.CallInfoForHandlerContext(ctx); !ok {
		return "", errors.MissingMetadata
	} else if token = callInfo.RequestHeader().Get(sessionIDKey); token == "" {
		return "", errors.MissingMetadata
	} else {
		return token, nil
	}
}

type authKeyType int

const authKey authKeyType = iota

func NewAuthContext(ctx context.Context, auth *model.Auth) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

func FromAuthContext(ctx context.Context) (*model.Auth, error) {
	a, ok := ctx.Value(authKey).(*model.Auth)
	if !ok {
		return nil, errors.MissingMetadata
	}

	return a, nil
}
