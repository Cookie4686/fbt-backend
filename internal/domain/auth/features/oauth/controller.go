package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fbt/backend/internal/domain/auth/common"
	"fbt/backend/internal/domain/auth/features/oauth/pb"
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

	pb.UnimplementedOAuthServer
}

func NewServer(service service.Service, repo Repo) *Server {
	return &Server{service, repo, pb.UnimplementedOAuthServer{}}
}

func RegisterService(service service.Service, repo Repo, s *grpc.Server) {
	pb.RegisterOAuthServer(s, NewServer(service, repo))
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	oauthRegistration, err := s.repo.GetOAuthRegistration(ctx, in.RegistrationID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(oauthRegistration.ExpiresAt) {
		if err := s.repo.DeleteOAuthRegistration(ctx, in.Provider, in.TokenID); err != nil {
			return nil, err
		} else {
			return nil, errors.RegistrationSessionExpire
		}
	}

	if (oauthRegistration.RegistrationID != in.RegistrationID) ||
		(oauthRegistration.IDToken != in.TokenID) {
		return nil, errors.BadRequest
	}

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        in.Username,
		Email:           in.Email,
		EmailVerified:   false,
		Password:        pgtype.Text{String: "", Valid: false},
		PasswordSalt:    pgtype.Text{String: "", Valid: false},
		PasswordEnabled: in.PasswordEnabled,
	}
	if in.PasswordEnabled {
		salt := make([]byte, 16)
		rand.Read(salt)
		passwordHash := argon2.IDKey([]byte(in.Password), salt, 2, 19*1024, 1, 32)
		user.Password = pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true}
		user.PasswordSalt = pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true}
	}

	session := &model.Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            user.Id,
		ExpiresAt:         time.Now().Add(model.SessionExpiresIn),
		TwoFactorVerified: false,
	}
	err = s.repo.OAuthRegister(ctx, in.RegistrationID, user, session)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{
		Session: &common.Session{
			Id:                session.Id,
			UserID:            session.UserId,
			TwoFactorVerified: session.TwoFactorVerified,
			ExpiresAt:         timestamppb.New(session.ExpiresAt),
		}}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	userOAuth, err := s.repo.GetUserOAuth(ctx, in.Provider, in.Token)
	if err != nil && err != errors.NotFound {
		return nil, err
	}

	var userId string = ""
	if err == nil {
		// Already Register OAuth
		userId = userOAuth.UserID
	} else if in.Email != "" {
		user, err := s.service.GetUserByEmail(ctx, in.Email)
		if err == nil {
			// Link OAuth to existing email
			err := s.repo.LinkOAuth(ctx, in.Provider, user.Id, in.Token)
			if err != nil {
				return nil, err
			}
			userId = user.Id
		} else if err != errors.NotFound {
			return nil, err
		}
	}

	if userId != "" {
		session, err := s.service.CreateSession(ctx, userId, false)
		if err != nil {
			return nil, err
		}
		return &pb.LoginReply{
			RegistrationNeeded: false,
			Session: &common.Session{
				Id:                session.Id,
				UserID:            session.UserId,
				TwoFactorVerified: session.TwoFactorVerified,
				ExpiresAt:         timestamppb.New(session.ExpiresAt),
			},
		}, nil
	} else {
		// No OAuth and No Email Registration
		oauthRegistration := &model.OauthRegistration{
			RegistrationID: util.GenerateBase32UUID(),
			IDToken:        in.Token,
			ExpiresAt:      time.Now().Add(model.SessionExpiresIn),
		}

		err := s.repo.CreateOAuthRegistration(ctx, in.Provider, oauthRegistration)
		if err != nil {
			return nil, err
		}
		return &pb.LoginReply{
			RegistrationNeeded: true,
			RegistrationID:     oauthRegistration.RegistrationID,
		}, nil
	}
}

func (s *Server) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	providers, err := s.repo.GetUserProvider(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	return &pb.StatusReply{Providers: providers}, nil
}
