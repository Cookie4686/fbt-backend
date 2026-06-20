package credentials

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"net/http"
)

type Feature struct {
	Register http.Handler
	Login    http.Handler
}

type Controller interface {
	Register(context.Context, *RegisterPayload) (*RegisterResponse, error)
	Login(context.Context, *LoginPayload) (*LoginResponse, error)
}

type Repo interface {
	Register(context.Context, *model.User, *model.Session) error
}
