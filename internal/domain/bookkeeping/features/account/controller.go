package account

import (
	"context"
	"fbt/backend/internal/domain/bookkeeping/features/account/pb"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

type Server struct {
	service service.Service
	repo    Repo

	pb.UnimplementedAccountServiceServer
}

func NewServer(service service.Service, repo Repo) *Server {
	return &Server{service, repo, pb.UnimplementedAccountServiceServer{}}
}

func RegisterService(service service.Service, repo Repo, s *grpc.Server) {
	pb.RegisterAccountServiceServer(s, NewServer(service, repo))
}

func (c *Server) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	accs, err := c.repo.GetAll(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}

	accounts := make([]*pb.Account, len(*accs))
	for idx, a := range *accs {
		accounts[idx] = toCommonAccount(&a)
	}

	return &pb.GetAllReply{Account: accounts}, nil
}

func (c *Server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	account := &model.Account{
		Name:    in.Name,
		IsDebit: in.IsDebit,
		UserId:  auth.Session.UserId,
	}

	accountID, err := c.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	account.ID = accountID
	return &pb.CreateReply{Account: toCommonAccount(account)}, nil
}

func (c *Server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	account := &model.Account{
		ID:      in.Id,
		Name:    in.Name,
		IsDebit: in.IsDebit,
		UserId:  auth.Session.UserId,
	}

	err = c.repo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateReply{Account: toCommonAccount(account)}, nil
}

func (c *Server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

func toCommonAccount(account *model.Account) *pb.Account {
	return &pb.Account{
		Id:      account.ID,
		Name:    account.Name,
		IsDebit: account.IsDebit,
		UserID:  account.UserId,
	}
}
