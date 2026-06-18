package user

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

type Feature struct {
	GetByUsername http.HandlerFunc
}

type Controller interface {
	GetByUsername(ctx context.Context, username string) (*GetByUsernameResponse, error)
}

type GetByUsernamePayload struct{}
type GetByUsernameResponse = util.Response[model.User]
