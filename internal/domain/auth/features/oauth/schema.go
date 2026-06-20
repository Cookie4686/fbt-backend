package oauth

import (
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

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
