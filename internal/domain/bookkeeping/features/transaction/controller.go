package transaction

import (
	"context"
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"
	"fbt/backend/gen/proto/go/bookkeeping/v1/bookkeepingv1connect"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/interceptor"
	"net/http"

	"connectrpc.com/connect"
)

type Server struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return bookkeepingv1connect.NewTransactionServiceHandler(&Server{service, repo}, opts...)
}

func (c *Server) GetAll(ctx context.Context, in *bookkeepingv1.TransactionServiceGetAllRequest) (*bookkeepingv1.TransactionServiceGetAllResponse, error) {
	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	tes, err := c.repo.GetAll(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}

	protoTes := make([]*bookkeepingv1.TransactionEntry, len(*tes))
	for idx, te := range *tes {
		protoTes[idx] = te.ToProto()
	}

	return &bookkeepingv1.TransactionServiceGetAllResponse{TransactionEntry: protoTes}, nil
}

func (c *Server) Create(ctx context.Context, in *bookkeepingv1.TransactionServiceCreateRequest) (*bookkeepingv1.TransactionServiceCreateResponse, error) {
	entries := make([]model.Entry, len(in.Entries))
	for idx, e := range in.Entries {
		entries[idx].AccountID = e.AccountId
		entries[idx].Amount = e.Amount
	}

	te := &model.TransactionEntry{
		Transaction: model.Transaction{Datetime: in.Time.AsTime()},
		Entries:     entries,
	}

	transactionID, err := c.repo.Create(ctx, te)
	if err != nil {
		return nil, err
	}

	te.TransactionID = transactionID

	return &bookkeepingv1.TransactionServiceCreateResponse{TransactionEntry: te.ToProto()}, nil
}

func (c *Server) Update(ctx context.Context, in *bookkeepingv1.TransactionServiceUpdateRequest) (*bookkeepingv1.TransactionServiceUpdateResponse, error) {
	err := c.repo.Update(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.TransactionServiceUpdateResponse{TransactionEntry: in.TransactionEntry}, nil
}

func (c *Server) Delete(ctx context.Context, in *bookkeepingv1.TransactionServiceDeleteRequest) (*bookkeepingv1.TransactionServiceDeleteResponse, error) {
	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	err = c.repo.Delete(ctx, auth.Session.UserId, in.Id)
	if err != nil {
		return nil, err
	}

	return &bookkeepingv1.TransactionServiceDeleteResponse{}, nil
}
