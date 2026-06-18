package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
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

func (s con) Register(ctx context.Context, body *RegisterPayload) (*RegisterResponse, error) {
	oauthRegistration, err := s.repo.GetOAuthRegistration(ctx, body.RegistrationID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(oauthRegistration.ExpiresAt) {
		if err := s.repo.DeleteOAuthRegistration(ctx, body.Provider, body.TokenID); err != nil {
			return nil, err
		} else {
			return nil, errors.RegistrationSessionExpire
		}
	}

	if (oauthRegistration.RegistrationID != body.RegistrationID) ||
		(oauthRegistration.IDToken != body.TokenID) {
		return nil, errors.BadRequest
	}

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        body.Username,
		Email:           body.Email,
		EmailVerified:   false,
		Password:        pgtype.Text{String: "", Valid: false},
		PasswordSalt:    pgtype.Text{String: "", Valid: false},
		PasswordEnabled: body.PasswordEnabled,
	}
	if body.PasswordEnabled {
		salt := make([]byte, 16)
		rand.Read(salt)
		passwordHash := argon2.IDKey([]byte(body.Password), salt, 2, 19*1024, 1, 32)
		user.Password = pgtype.Text{String: base64.StdEncoding.EncodeToString(passwordHash), Valid: true}
		user.PasswordSalt = pgtype.Text{String: base64.StdEncoding.EncodeToString(salt), Valid: true}
	}

	session := &model.Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            user.Id,
		ExpiresAt:         time.Now().Add(model.SessionExpiresIn),
		TwoFactorVerified: false,
	}
	err = s.repo.OAuthRegister(ctx, body.RegistrationID, user, session)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{StatusCode: http.StatusOK, Payload: session}, nil
}

func (s con) Login(ctx context.Context, body *LoginPayload) (*LoginResponse, error) {
	userOAuth, err := s.repo.GetUserOAuth(ctx, body.Provider, body.IDToken)
	if err != nil && err != errors.NotFound {
		return nil, err
	}

	var userId string = ""
	if err == nil {
		// Already Register OAuth
		userId = userOAuth.UserID
	} else if body.Email != nil {
		user, err := s.service.GetUserByEmail(ctx, *body.Email)
		if err == nil {
			// Link OAuth to existing email
			err := s.repo.LinkOAuth(ctx, body.Provider, user.Id, body.IDToken)
			if err != nil {
				return nil, err
			}
			userId = user.Id
		} else if err != errors.NotFound {
			return nil, err
		}
	}

	if userId != "" {
		session, err := s.service.CreateSession(ctx, userId)
		if err != nil {
			return nil, err
		}
		return &LoginResponse{StatusCode: http.StatusOK, Payload: &LoginResponsePayload{
			Session:            session,
			RegistrationNeeded: false,
		}}, nil
	} else {
		// No OAuth and No Email Registration
		oauthRegistration := &model.OauthRegistration{
			RegistrationID: util.GenerateBase32UUID(),
			IDToken:        body.IDToken,
			ExpiresAt:      time.Now().Add(model.SessionExpiresIn),
		}

		err := s.repo.CreateOAuthRegistration(ctx, body.Provider, oauthRegistration)
		if err != nil {
			return nil, err
		}
		return &LoginResponse{StatusCode: http.StatusOK, Payload: &LoginResponsePayload{
			RegistrationNeeded: true,
			RegistrationId:     &oauthRegistration.RegistrationID,
		}}, nil
	}
}

func (s con) Status(ctx context.Context, auth *model.Auth) (*StatusResponse, error) {
	providers, err := s.repo.GetUserProvider(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	return &StatusResponse{StatusCode: http.StatusOK, Payload: &StatusResponsePaylaod{
		Providers: providers,
	}}, nil
}
