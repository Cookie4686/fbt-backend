package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/mfa/pb"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"

	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service service.Service
	repo    Repo

	pb.UnimplementedMFAServer
}

func NewServer(service service.Service, repo Repo) *Server {
	return &Server{service, repo, pb.UnimplementedMFAServer{}}
}

func RegisterService(service service.Service, repo Repo, s *grpc.Server) {
	pb.RegisterMFAServer(s, NewServer(service, repo))
}

func (s *Server) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	userMfaList, err := s.repo.GetMFAList(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}
	return &pb.StatusReply{TotpEnabled: userMfaList.Totp}, nil
}

func (s *Server) TOTPValidate(ctx context.Context, in *pb.TOTPValidateRequest) (*pb.TOTPValidateReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	userTotp, err := s.repo.GetTOTP(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	secret, err := s.service.Decrypt(userTotp.Key)
	if err != nil {
		return nil, err
	}

	isValid := totp.Validate(in.Code, *secret)

	return &pb.TOTPValidateReply{
		IsValid: isValid,
	}, nil
}

func (s *Server) TOTPUpsertKey(ctx context.Context, in *pb.TOTPUpsertRequest) (*pb.TOTPUpsertReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	encryptedKey, err := s.service.Encrypt(in.Key)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpsertTOTP(ctx, *encryptedKey, auth.User.Id)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}
	err = s.service.InvalidateSession(ctx, &auth.Session)
	if err != nil {
		return nil, err
	}

	return &pb.TOTPUpsertReply{Session: &common.Session{
		Id:                session.Id,
		UserID:            session.UserId,
		TwoFactorVerified: session.TwoFactorVerified,
		ExpiresAt:         timestamppb.New(session.ExpiresAt),
	}}, nil
}
