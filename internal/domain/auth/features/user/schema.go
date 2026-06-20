package user

import (
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

type GetByUsernamePayload struct{}
type GetByUsernameResponse = util.Response[model.User]
