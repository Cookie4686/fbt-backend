package credentials

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"
)

type con struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewCredentialServiceHandler(&con{service, repo}, opts...)
}

func (s *con) Register(ctx context.Context, req *authv1.CredentialServiceRegisterRequest) (*authv1.CredentialServiceRegisterResponse, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	passwordHash := argon2.IDKey([]byte(req.Password), salt, 2, 19*1024, 1, 32)

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        req.Username,
		Email:           req.Email,
		EmailVerified:   false,
		Password:        pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true},
		PasswordSalt:    pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true},
		PasswordEnabled: true,
	}
	session := model.NewSession(user.Id, false)
	err := s.repo.Register(ctx, user, session)
	if err != nil {
		return nil, err
	}

	return &authv1.CredentialServiceRegisterResponse{Session: session.ToProto()}, nil
}

func (s *con) Login(ctx context.Context, req *authv1.CredentialServiceLoginRequest) (*authv1.CredentialServiceLoginResponse, error) {
	// Get User Data From Database
	user, err := s.service.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !user.PasswordEnabled {
		return nil, errors.Unauthorized
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
	passwordHash := argon2.IDKey([]byte(req.Password), storedSalt, 2, 19*1024, 1, 32)
	if subtle.ConstantTimeCompare(passwordHash, storedHash) == 1 {
		// Create Session in Database
		session, err := s.service.CreateSession(ctx, user.Id, false)
		if err != nil {
			return nil, err
		}
		return &authv1.CredentialServiceLoginResponse{Session: session.ToProto()}, nil
	} else {
		return nil, errors.Unauthorized
	}
}
