package credentials

import (
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

type RegisterPayload struct {
	Username string `json:"username" validate:"required,min=3,max=255,alphanumunicode"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Email    string `json:"email" validate:"required,email"`
}
type RegisterResponse = util.Response[model.Session]

type LoginPayload struct {
	Username string `json:"username" validate:"required,min=3,max=255,alphanumunicode"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}
type LoginResponse = util.Response[model.Session]
