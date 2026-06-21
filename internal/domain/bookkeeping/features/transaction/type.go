package transaction

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
	GetAll(ctx context.Context, userID string) (*[]model.TransactionEntry, error)
	Create(context.Context, *model.TransactionEntry) (transactionID int64, err error)
	Update(context.Context, *model.TransactionEntry) error
	Delete(ctx context.Context, userID string, transactionID int64) error
}
