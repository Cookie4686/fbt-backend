package account

import (
	"context"
	auth "fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/bookkeeping/model"
	"net/http"
)

type Feature struct {
	GetAll http.Handler
	Create http.Handler
	Update http.Handler
	Delete http.Handler
}

type Controller interface {
	GetAll(context.Context, *auth.Auth) (*GetAllResponse, error)
	Create(context.Context, *auth.Auth, *CreatePayload) (*CreateResponse, error)
	Update(context.Context, *auth.Auth, *UpdatePayload) (*UpdateResponse, error)
	Delete(context.Context, *auth.Auth, *DeletePayload) (*DeleteResponse, error)
}

type Repo interface {
	GetAll(ctx context.Context, userID string) (*[]model.Account, error)
	Create(context.Context, *model.Account) (accountID int32, err error)
	Update(context.Context, *model.Account) error
	Delete(ctx context.Context, userID string, accountID int32) error
}
