package credentials

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/credentials/pb"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service service.Service
	repo    Repo

	pb.UnimplementedCredentialsServer
}

func NewServer(service service.Service, repo Repo) *Server {
	return &Server{service, repo, pb.UnimplementedCredentialsServer{}}
}

func RegisterService(service service.Service, repo Repo, s *grpc.Server) {
	pb.RegisterCredentialsServer(s, NewServer(service, repo))
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	passwordHash := argon2.IDKey([]byte(in.Password), salt, 2, 19*1024, 1, 32)

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        in.Username,
		Email:           in.Email,
		EmailVerified:   false,
		Password:        pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true},
		PasswordSalt:    pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true},
		PasswordEnabled: true,
	}
	session := &model.Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            user.Id,
		ExpiresAt:         time.Now().Add(model.SessionExpiresIn),
		TwoFactorVerified: false,
	}
	err := s.repo.Register(ctx, user, session)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{Session: &common.Session{
		Id:                session.Id,
		UserID:            session.UserId,
		ExpiresAt:         timestamppb.New(session.ExpiresAt),
		TwoFactorVerified: session.TwoFactorVerified,
	}}, nil
}

func (s *Server) Login(ctx context.Context, body *pb.LoginRequest) (*pb.LoginReply, error) {
	// Get User Data From Database
	user, err := s.service.GetUserByUsername(ctx, body.Username)
	if err != nil {
		return nil, err
	}

	// TODO: Handle Non-Credentials User
	storedHash, err := base64.StdEncoding.DecodeString(user.Password.String)
	if err != nil {
		return nil, err
	}
	storedSalt, err := base64.StdEncoding.DecodeString(user.PasswordSalt.String)
	if err != nil {
		return nil, err
	}

	// Compare Password Hash
	passwordHash := argon2.IDKey([]byte(body.Password), storedSalt, 2, 19*1024, 1, 32)
	if subtle.ConstantTimeCompare(passwordHash, storedHash) == 1 {
		// Create Session in Database
		session, err := s.service.CreateSession(ctx, user.Id)
		if err != nil {
			return nil, err
		}
		return &pb.LoginReply{Session: &common.Session{
			Id:                session.Id,
			UserID:            session.UserId,
			ExpiresAt:         timestamppb.New(session.ExpiresAt),
			TwoFactorVerified: session.TwoFactorVerified,
		}}, nil
	} else {
		return nil, errors.Unauthorized
	}
}
