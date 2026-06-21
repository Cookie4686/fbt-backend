package account

import (
	"context"
	auth "fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"net/http"
)

type con struct {
	service service.Service
	repo    Repo
}

func NewController(service service.Service, repo Repo) Controller {
	return Controller(&con{service, repo})
}

func (c *con) GetAll(ctx context.Context, auth *auth.Auth) (*GetAllResponse, error) {
	accounts, err := c.repo.GetAll(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}
	return &GetAllResponse{StatusCode: http.StatusOK, Payload: &accounts}, nil
}

func (c *con) Create(ctx context.Context, auth *auth.Auth, payload *CreatePayload) (*CreateResponse, error) {
	account := model.Account{
		Name:    payload.Name,
		IsDebit: payload.IsDebit,
		UserId:  auth.Session.UserId,
	}
	accountID, err := c.repo.Create(ctx, &account)
	if err != nil {
		return nil, err
	}

	account.ID = accountID
	return &CreateResponse{StatusCode: http.StatusOK, Payload: &account}, nil
}

func (c *con) Update(ctx context.Context, auth *auth.Auth, payload *UpdatePayload) (*UpdateResponse, error) {
	account := &model.Account{
		ID:      payload.ID,
		Name:    payload.Name,
		IsDebit: payload.IsDebit,
		UserId:  auth.Session.UserId,
	}
	err := c.repo.Update(ctx, account)
	if err != nil {
		return nil, err
	}
	return &UpdateResponse{StatusCode: http.StatusOK, Payload: account}, nil
}

func (c *con) Delete(ctx context.Context, auth *auth.Auth, payload *DeletePayload) (*DeleteResponse, error) {
	err := c.repo.Delete(ctx, auth.Session.UserId, payload.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{StatusCode: http.StatusOK, Payload: nil}, nil
}
