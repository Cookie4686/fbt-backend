// Package transaction for transaction services (General Journal)
package transaction

import (
	"context"
	"net/http"

	"fbt/backend/gen/proto/go/bookkeeping/v1/bookkeepingv1connect"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/interceptor"

	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

type con struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return bookkeepingv1connect.NewTransactionServiceHandler(NewController(service, repo), opts...)
}

func NewController(service service.Service, repo Repo) *con {
	return &con{service, repo}
}

func (c *con) GetAll(ctx context.Context, in *bookkeepingv1.TransactionServiceGetAllRequest) (*bookkeepingv1.TransactionServiceGetAllResponse, error) {
	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	tes, err := c.repo.GetAll(ctx, auth.Session.UserID)
	if err != nil {
		return nil, err
	}

	protoTes := make([]*bookkeepingv1.TransactionEntry, len(*tes))
	for idx, te := range *tes {
		protoTes[idx] = te.ToProto()
	}

	return &bookkeepingv1.TransactionServiceGetAllResponse{TransactionEntry: protoTes}, nil
}

func (c *con) Create(ctx context.Context, in *bookkeepingv1.TransactionServiceCreateRequest) (*bookkeepingv1.TransactionServiceCreateResponse, error) {
	entries := protoEntriesToEntries(in.Entries)

	transactionID, err := c.repo.Create(ctx, model.Transaction{
		Description: pgtype.Text{String: in.Description, Valid: in.Description != ""},
		Note:        pgtype.Text{String: in.Note, Valid: in.Note != ""},
		Datetime:    in.Time.AsTime(),
	}, entries)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.TransactionServiceCreateResponse{TransactionEntry: &bookkeepingv1.TransactionEntry{
		Id:          transactionID,
		Time:        in.Time,
		Description: in.Description,
		Note:        in.Note,
		Entries:     in.Entries,
	}}, nil
}

func (c *con) Update(ctx context.Context, in *bookkeepingv1.TransactionServiceUpdateRequest) (*bookkeepingv1.TransactionServiceUpdateResponse, error) {
	entries := protoEntriesToEntries(in.Entries)

	err := c.repo.Update(ctx, model.Transaction{
		Description: pgtype.Text{String: in.Description, Valid: in.Description != ""},
		Note:        pgtype.Text{String: in.Note, Valid: in.Note != ""},
		Datetime:    in.Time.AsTime(),
	}, entries)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.TransactionServiceUpdateResponse{TransactionEntry: &bookkeepingv1.TransactionEntry{
		Id:          in.Id,
		Time:        in.Time,
		Description: in.Description,
		Note:        in.Note,
		Entries:     in.Entries,
	}}, nil
}

func (c *con) Delete(ctx context.Context, in *bookkeepingv1.TransactionServiceDeleteRequest) (*bookkeepingv1.TransactionServiceDeleteResponse, error) {
	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	err = c.repo.Delete(ctx, auth.Session.UserID, in.Id)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.TransactionServiceDeleteResponse{}, nil
}

func protoEntriesToEntries(protoEntries []*bookkeepingv1.Entry) []model.Entry {
	entries := make([]model.Entry, len(protoEntries))
	for idx, e := range protoEntries {
		entries[idx].AccountID = e.AccountId
		entries[idx].Amount = e.Amount
	}

	return entries
}
