package session

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

type Feature struct {
	Logout   http.HandlerFunc
	Validate http.HandlerFunc
}

type Controller interface {
	Logout(context.Context, *model.Auth) (*LogoutResponse, error)
	Validate(context.Context, *model.Auth) (*ValidateResponse, error)
}

type LogoutPayload struct{}
type LogoutResponse = util.Response[struct{}]

type ValidatePayload struct{}
type ValidateResponse = util.Response[model.Auth]
