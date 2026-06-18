package credentials

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

type Feature struct {
	Register http.HandlerFunc
	Login    http.HandlerFunc
}

type Controller interface {
	Register(context.Context, *RegisterPayload) (*RegisterResponse, error)
	Login(context.Context, *LoginPayload) (*LoginResponse, error)
}

type RegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterResponse = util.Response[model.Session]

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse = util.Response[model.Session]
