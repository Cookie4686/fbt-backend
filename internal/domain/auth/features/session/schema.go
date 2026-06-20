package session

import (
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

type LogoutPayload struct{}
type LogoutResponse = util.Response[struct{}]

type ValidatePayload struct{}
type ValidateResponse = util.Response[model.Auth]
