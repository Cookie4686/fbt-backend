package transaction

import (
	"context"
	"fbt/backend/internal/domain/bookkeeping/features/transaction/pb"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service service.Service
	repo    Repo

	pb.UnimplementedTransactionServer
}

func NewServer(service service.Service, repo Repo) *Server {
	return &Server{service, repo, pb.UnimplementedTransactionServer{}}
}

func RegisterService(service service.Service, repo Repo, s *grpc.Server) {
	pb.RegisterTransactionServer(s, NewServer(service, repo))
}

func (c *Server) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllReply, error) {
	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	tes, err := c.repo.GetAll(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}

	commonTes := make([]*pb.TransactionEntry, len(*tes))
	for idx, te := range *tes {
		commonTes[idx] = toCommonTransactionEntry(&te)
	}

	return &pb.GetAllReply{TransactionEntry: commonTes}, nil
}

func (c *Server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	entries := make([]model.Entry, len(in.Entries))
	for idx, e := range in.Entries {
		entries[idx].AccountID = e.AccountID
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

	return &pb.CreateReply{TransactionEntry: toCommonTransactionEntry(te)}, nil
}

func (c *Server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	err := c.repo.Update(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateReply{TransactionEntry: in.TransactionEntry}, nil
}

func (c *Server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	err = c.repo.Delete(ctx, auth.Session.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteReply{}, nil
}

func toCommonTransactionEntry(te *model.TransactionEntry) *pb.TransactionEntry {
	entries := make([]*pb.Entry, len(te.Entries))

	for idx, e := range te.Entries {
		entries[idx] = &pb.Entry{
			AccountID: e.AccountID,
			Amount:    e.Amount,
		}
	}

	return &pb.TransactionEntry{
		Id:      te.TransactionID,
		Time:    timestamppb.New(te.Datetime),
		Entries: entries,
	}
}
