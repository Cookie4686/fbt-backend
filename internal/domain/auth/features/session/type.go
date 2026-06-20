package session

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"net/http"
)

type Feature struct {
	AUTH_Logout   http.Handler
	AUTH_Validate http.Handler
}

type Controller interface {
	Logout(context.Context, *model.Auth) (*LogoutResponse, error)
	Validate(context.Context, *model.Auth) (*ValidateResponse, error)
}
