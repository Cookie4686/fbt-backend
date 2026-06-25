package session

import (
	"context"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/session/pb"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service service.Service

	pb.UnimplementedSessionServer
}

func NewServer(service service.Service) *Server {
	return &Server{service, pb.UnimplementedSessionServer{}}
}

func RegisterService(service service.Service, s *grpc.Server) {
	pb.RegisterSessionServer(s, NewServer(service))
}

func (s *Server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateReply{
		User: &common.User{
			Id:              auth.User.Id,
			Username:        auth.User.Username,
			Email:           auth.User.Email,
			EmailVerified:   auth.User.EmailVerified,
			PasswordEnabled: auth.User.PasswordEnabled,
		},
		Session: &common.Session{
			Id:                auth.Session.Id,
			UserID:            auth.Session.UserId,
			TwoFactorVerified: auth.Session.TwoFactorVerified,
			ExpiresAt:         timestamppb.New(auth.Session.ExpiresAt),
		},
	}, nil
}

func (s *Server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := util.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	err = s.service.InvalidateSession(ctx, &model.Session{Id: auth.Session.Id})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
