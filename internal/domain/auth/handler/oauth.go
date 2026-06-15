package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"
)

func (s *AuthHandler) OAuthRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		RegistrationID  string `json:"registration_id"`
		Provider        string `json:"provider"`
		TokenID         string `json:"id_token"`
		Username        string `json:"username"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordEnabled bool   `json:"password_enabled"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	oauthRegistration, err := s.repo.GetOAuthRegistration(ctx, body.RegistrationID)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	if time.Now().After(oauthRegistration.ExpiresAt) {
		if err := s.repo.DeleteOAuthRegistration(ctx, body.Provider, body.TokenID); err != nil {
			util.SendError(s.logger, w, r, err)
		} else {
			util.SendError(s.logger, w, r, errors.RegistrationSessionExpire)
		}
		return
	}

	if (oauthRegistration.RegistrationID != body.RegistrationID) ||
		(oauthRegistration.IDToken != body.TokenID) {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		util.SendError(s.logger, w, r, err)
		return
	}

	util.SendJson(w, http.StatusOK, &session)
}

func (s *AuthHandler) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		IDToken  string  `json:"token"`
		Email    *string `json:"email"`
		Provider string  `json:"provider"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	response := struct {
		Session            *model.Session `json:"session"`
		RegistrationId     string         `json:"registration_id"`
		RegistrationNeeded bool           `json:"registration_needed"`
	}{}

	userOAuth, err := s.repo.GetUserOAuth(ctx, body.Provider, body.IDToken)
	if err != nil && err != errors.NotFound {
		util.SendError(s.logger, w, r, err)
		return
	}

	var userId string = ""
	if err == nil {
		// Already Register OAuth
		userId = userOAuth.UserID
	} else if body.Email != nil {
		user, err := s.repo.GetUserByEmail(ctx, *body.Email)
		if err == nil {
			// Link OAuth to existing email
			err := s.repo.LinkOAuth(ctx, body.Provider, user.Id, body.IDToken)
			if err != nil {
				util.SendError(s.logger, w, r, err)
				return
			}
			userId = user.Id
		} else if err != errors.NotFound {
			util.SendError(s.logger, w, r, err)
			return
		}
	}

	if userId != "" {
		session, err := s.repo.CreateSession(ctx, userId)
		if err != nil {
			util.SendError(s.logger, w, r, err)
			return
		}
		response.Session = session
		response.RegistrationNeeded = false
	} else {
		// No OAuth and No Email Registration
		oauthRegistration := &model.OauthRegistration{
			RegistrationID: util.GenerateBase32UUID(),
			IDToken:        body.IDToken,
			ExpiresAt:      time.Now().Add(model.SessionExpiresIn),
		}

		err := s.repo.CreateOAuthRegistration(ctx, body.Provider, oauthRegistration)
		if err != nil {
			util.SendError(s.logger, w, r, err)
			return
		}
		response.RegistrationId = oauthRegistration.RegistrationID
		response.RegistrationNeeded = true
	}
	util.SendJson(w, http.StatusOK, &response)

}
