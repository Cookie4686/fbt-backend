package transaction

import (
	"context"
	auth "fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"net/http"
	"time"
)

type con struct {
	service service.Service
	repo    Repo
}

func NewController(service service.Service, repo Repo) Controller {
	return Controller(&con{service, repo})
}

func (c *con) GetAll(ctx context.Context, auth *auth.Auth) (*GetAllResponse, error) {
	te, err := c.repo.GetAll(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}
	return &GetAllResponse{StatusCode: http.StatusOK, Payload: te}, nil
}

func (c *con) Create(ctx context.Context, auth *auth.Auth, payload *CreatePayload) (*CreateResponse, error) {
	datetime, err := time.Parse(time.RFC3339, payload.Datetime)
	if err != nil {
		return nil, err
	}
	entries := make([]model.Entry, len(payload.Entries))
	for idx, e := range payload.Entries {
		entries[idx].AccountID = e.AccountID
		entries[idx].Amount = e.Amount
	}
	te := &model.TransactionEntry{
		Transaction: model.Transaction{Datetime: datetime},
		Entries:     entries,
	}
	transactionID, err := c.repo.Create(ctx, te)
	if err != nil {
		return nil, err
	}
	te.TransactionID = transactionID

	return &CreateResponse{StatusCode: http.StatusOK, Payload: te}, nil
}

func (c *con) Update(ctx context.Context, auth *auth.Auth, payload *UpdatePayload) (*UpdateResponse, error) {
	err := c.repo.Update(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &UpdateResponse{StatusCode: http.StatusOK, Payload: nil}, nil
}

func (c *con) Delete(ctx context.Context, auth *auth.Auth, payload *DeletePayload) (*DeleteResponse, error) {
	err := c.repo.Delete(ctx, auth.Session.UserId, payload.TransactionID)
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{StatusCode: http.StatusOK, Payload: nil}, nil
}
