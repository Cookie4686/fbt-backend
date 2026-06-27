package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"
)

type Server struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewOAuthServiceHandler(&Server{service, repo}, opts...)
}

func (s *Server) Register(ctx context.Context, in *authv1.OAuthServiceRegisterRequest) (*authv1.OAuthServiceRegisterResponse, error) {
	oauthRegistration, err := s.repo.GetOAuthRegistration(ctx, in.RegistrationId)
	if err != nil {
		return nil, err
	}

	if time.Now().After(oauthRegistration.ExpiresAt) {
		if err := s.repo.DeleteOAuthRegistration(ctx, in.Provider, in.TokenId); err != nil {
			return nil, err
		} else {
			return nil, errors.RegistrationSessionExpire
		}
	}

	if (oauthRegistration.RegistrationID != in.RegistrationId) ||
		(oauthRegistration.IDToken != in.TokenId) {
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
	err = s.repo.OAuthRegister(ctx, in.RegistrationId, user, session)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceRegisterResponse{Session: session.ToProto()}, nil
}

func (s *Server) Login(ctx context.Context, in *authv1.OAuthServiceLoginRequest) (*authv1.OAuthServiceLoginResponse, error) {
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
		return &authv1.OAuthServiceLoginResponse{
			RegistrationNeeded: false,
			Session:            session.ToProto(),
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
		return &authv1.OAuthServiceLoginResponse{
			RegistrationNeeded: true,
			RegistrationId:     oauthRegistration.RegistrationID,
		}, nil
	}
}

func (s *Server) Status(ctx context.Context, in *authv1.OAuthServiceStatusRequest) (*authv1.OAuthServiceStatusResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	providers, err := s.repo.GetUserProvider(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceStatusResponse{Providers: providers}, nil
}
