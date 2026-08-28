// Package oauth for oauth-related authentication
package oauth

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

const OAuthRegistrationMaxAge = 15 * time.Minute

type Server struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewOAuthServiceHandler(NewController(service, repo), opts...)
}

func NewController(service service.Service, repo Repo) *Server {
	return &Server{service, repo}
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
		ID:              util.GenerateBase32UUID(),
		Username:        in.Username,
		Email:           in.Email,
		EmailVerified:   oauthRegistration.EmailVerified,
		Password:        pgtype.Text{String: "", Valid: false},
		PasswordSalt:    pgtype.Text{String: "", Valid: false},
		PasswordEnabled: in.PasswordEnabled,
	}
	if in.PasswordEnabled {
		salt := s.service.GenerateSalt()
		passwordHash := s.service.HashPassword(in.Password, salt)
		user.Password = pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true}
		user.PasswordSalt = pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true}
	}

	session := model.NewSession(user.ID, false)

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

	if err == nil {
		// Already Register OAuth
		session, err := s.service.CreateSession(ctx, userOAuth.UserID, false)
		if err != nil {
			return nil, err
		}

		return &authv1.OAuthServiceLoginResponse{
			RegistrationNeeded: false,
			Session:            session.ToProto(),
		}, nil
	}

	if in.Email != nil {
		user, err := s.service.GetUserByEmail(ctx, *in.Email)
		if err != nil && err != errors.NotFound {
			return nil, err
		}

		if err != errors.NotFound {
			// Link OAuth to existing email
			session := model.NewSession(user.ID, false)

			err = s.repo.LinkOAuth(ctx, in.Provider, user.ID, in.Token, session)
			if err != nil {
				return nil, err
			}

			return &authv1.OAuthServiceLoginResponse{
				RegistrationNeeded: false,
				Session:            session.ToProto(),
			}, nil
		}
	}

	// No OAuth and No Email Registration
	oauthRegistration := &model.OauthRegistration{
		RegistrationID: util.GenerateBase32UUID(),
		IDToken:        in.Token,
		EmailVerified:  in.EmailVerified,
		ExpiresAt:      time.Now().Add(OAuthRegistrationMaxAge),
	}

	err = s.repo.CreateOAuthRegistration(ctx, in.Provider, oauthRegistration)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceLoginResponse{
		RegistrationNeeded: true,
		RegistrationId:     oauthRegistration.RegistrationID,
	}, nil
}

func (s *Server) Status(ctx context.Context, in *authv1.OAuthServiceStatusRequest) (*authv1.OAuthServiceStatusResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	providers, err := s.repo.GetUserProvider(ctx, auth.User.ID)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceStatusResponse{Providers: providers}, nil
}

func (s *Server) Link(ctx context.Context, in *authv1.OAuthServiceLinkRequest) (*authv1.OAuthServiceLinkResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	session := model.NewSession(auth.Session.UserID, true)

	err = s.repo.LinkOAuth(ctx, in.Provider, auth.User.ID, in.Token, session)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceLinkResponse{Session: session.ToProto()}, nil
}

func (s *Server) Unlink(ctx context.Context, in *authv1.OAuthServiceUnlinkRequest) (*authv1.OAuthServiceUnlinkResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.repo.UnLinkOAuth(ctx, in.Provider, auth.User.ID)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, auth.User.ID, true)
	if err != nil {
		return nil, err
	}

	return &authv1.OAuthServiceUnlinkResponse{Session: session.ToProto()}, nil
}
