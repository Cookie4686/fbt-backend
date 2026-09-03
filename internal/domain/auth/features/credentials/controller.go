// Package credentials for password-related authentication services
package credentials

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"net/http"

	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

type Controller struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewCredentialServiceHandler(NewController(service, repo), opts...)
}

func NewController(service service.Service, repo Repo) *Controller {
	return &Controller{service, repo}
}

func (s *Controller) Register(ctx context.Context, in *authv1.CredentialServiceRegisterRequest) (*authv1.CredentialServiceRegisterResponse, error) {
	salt := s.service.GenerateSalt()
	passwordHash := s.service.HashPassword(in.Password, salt)

	user := &model.User{
		ID:              util.GenerateBase32UUID(),
		Username:        in.Username,
		Email:           in.Email,
		EmailVerified:   false,
		Password:        pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true},
		PasswordSalt:    pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true},
		PasswordEnabled: true,
	}
	session := model.NewSession(user.ID, false)

	err := s.repo.Register(ctx, user, session)
	if err != nil {
		return nil, err
	}

	return &authv1.CredentialServiceRegisterResponse{Session: session.ToProto()}, nil
}

func (s *Controller) Login(ctx context.Context, in *authv1.CredentialServiceLoginRequest) (*authv1.CredentialServiceLoginResponse, error) {
	// Get User Data From Database
	user, err := s.service.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	storedHash, err := base64.StdEncoding.DecodeString(user.Password.String)
	if err != nil {
		return nil, err
	}

	storedSalt, err := base64.StdEncoding.DecodeString(user.PasswordSalt.String)
	if err != nil {
		return nil, err
	}

	// Compare Password Hash
	passwordHash := s.service.HashPassword(in.Password, storedSalt)
	if subtle.ConstantTimeCompare(passwordHash, storedHash) == 1 && user.PasswordEnabled {
		// Create Session in Database
		session, err := s.service.CreateSession(ctx, user.ID, false)
		if err != nil {
			return nil, err
		}

		return &authv1.CredentialServiceLoginResponse{Session: session.ToProto()}, nil
	} else {
		return nil, errors.Unauthorized
	}
}
