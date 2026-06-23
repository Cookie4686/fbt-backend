package user

import (
	"context"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/user/pb"
	"fbt/backend/internal/domain/auth/service"

	"google.golang.org/grpc"
)

type Server struct {
	service service.Service

	pb.UnimplementedUserServer
}

func NewServer(service service.Service) *Server {
	return &Server{service, pb.UnimplementedUserServer{}}
}

func RegisterService(service service.Service, s *grpc.Server) {
	pb.RegisterUserServer(s, NewServer(service))
}

func (s *Server) GetByUsername(ctx context.Context, in *pb.GetByUsernameRequest) (*pb.GetByUsernameReply, error) {
	user, err := s.service.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	return &pb.GetByUsernameReply{User: &common.User{
		Id:              user.Id,
		Username:        user.Username,
		Email:           user.Email,
		EmailVerified:   user.EmailVerified,
		PasswordEnabled: user.PasswordEnabled,
	}}, nil
}
