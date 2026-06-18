package credentials

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

type con struct {
	service *service.AuthService
	repo    *repo
}

func NewController(service *service.AuthService, db *pgxpool.Pool) Controller {
	return Controller(con{service: service, repo: newRepo(db)})
}

func (s con) Register(ctx context.Context, payload *RegisterPayload) (*RegisterResponse, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	passwordHash := argon2.IDKey([]byte(payload.Password), salt, 2, 19*1024, 1, 32)

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        payload.Username,
		Email:           payload.Email,
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

	return &RegisterResponse{StatusCode: http.StatusOK, Payload: session}, nil
}

func (h con) Login(ctx context.Context, body *LoginPayload) (*LoginResponse, error) {
	// Get User Data From Database
	user, err := h.service.GetUserByUsername(ctx, body.Username)
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
		session, err := h.service.CreateSession(ctx, user.Id)
		if err != nil {
			return nil, err
		}
		return &LoginResponse{StatusCode: http.StatusOK, Payload: session}, nil
	} else {
		return &LoginResponse{StatusCode: http.StatusUnauthorized}, nil
	}
}
