// Package account for account management service (Chart of Accounts)
package account

import (
	"context"
	"net/http"

	"fbt/backend/gen/proto/go/bookkeeping/v1/bookkeepingv1connect"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/interceptor"

	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	"connectrpc.com/connect"
)

type con struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return bookkeepingv1connect.NewAccountServiceHandler(NewController(service, repo), opts...)
}

func NewController(service service.Service, repo Repo) *con {
	return &con{service, repo}
}

func (c *con) GetAll(ctx context.Context, in *bookkeepingv1.AccountServiceGetAllRequest) (*bookkeepingv1.AccountServiceGetAllResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	accs, err := c.repo.GetAll(ctx, auth.Session.UserID)
	if err != nil {
		return nil, err
	}

	accounts := make([]*bookkeepingv1.Account, len(*accs))
	for idx, a := range *accs {
		accounts[idx] = a.ToProto()
	}

	return &bookkeepingv1.AccountServiceGetAllResponse{Account: accounts}, nil
}

func (c *con) Create(ctx context.Context, in *bookkeepingv1.AccountServiceCreateRequest) (*bookkeepingv1.AccountServiceCreateResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	account := &model.Account{
		Name:    in.Name,
		IsDebit: in.IsDebit,
		UserID:  auth.Session.UserID,
	}

	accountID, err := c.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	account.ID = accountID

	return &bookkeepingv1.AccountServiceCreateResponse{Account: account.ToProto()}, nil
}

func (c *con) Update(ctx context.Context, in *bookkeepingv1.AccountServiceUpdateRequest) (*bookkeepingv1.AccountServiceUpdateResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	account := &model.Account{
		ID:      in.Id,
		Name:    in.Name,
		IsDebit: in.IsDebit,
		UserID:  auth.Session.UserID,
	}

	err = c.repo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.AccountServiceUpdateResponse{Account: account.ToProto()}, nil
}

func (c *con) Delete(ctx context.Context, in *bookkeepingv1.AccountServiceDeleteRequest) (*bookkeepingv1.AccountServiceDeleteResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	err = c.repo.Delete(ctx, auth.Session.UserID, in.Id)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.AccountServiceDeleteResponse{}, nil
}
