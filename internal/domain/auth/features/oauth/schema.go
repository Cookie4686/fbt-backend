package oauth

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

type Feature struct {
	Register http.HandlerFunc
	Login    http.HandlerFunc
	Status   http.HandlerFunc
}

type Controller interface {
	Register(context.Context, *RegisterPayload) (*RegisterResponse, error)
	Login(context.Context, *LoginPayload) (*LoginResponse, error)
	Status(context.Context, *model.Auth) (*StatusResponse, error)
}

type RegisterPayload struct {
	RegistrationID  string `json:"registration_id"`
	Provider        string `json:"provider"`
	TokenID         string `json:"id_token"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordEnabled bool   `json:"password_enabled"`
}
type RegisterResponse = util.Response[model.Session]

type LoginPayload struct {
	IDToken  string  `json:"token"`
	Email    *string `json:"email"`
	Provider string  `json:"provider"`
}
type LoginResponsePayload struct {
	Session            *model.Session `json:"session"`
	RegistrationId     *string        `json:"registration_id"`
	RegistrationNeeded bool           `json:"registration_needed"`
}
type LoginResponse = util.Response[LoginResponsePayload]

type StatusPayload struct{}
type StatusResponsePaylaod struct {
	Providers []string `json:"providers"`
}
type StatusResponse = util.Response[StatusResponsePaylaod]
