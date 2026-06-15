package handler

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/argon2"
)

func (s *AuthHandler) CredentialsRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	salt := make([]byte, 16)
	rand.Read(salt)
	passwordHash := argon2.IDKey([]byte(body.Password), salt, 2, 19*1024, 1, 32)

	user := &model.User{
		Id:              util.GenerateBase32UUID(),
		Username:        body.Username,
		Email:           body.Email,
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
	err := s.repo.CredentialsRegister(ctx, user, session)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	util.SendJson(w, http.StatusOK, &session)
}

func (s *AuthHandler) CredentialsLogin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	// Get User Data From Database
	user, err := s.repo.GetUserByUsername(ctx, body.Username)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}
	// TODO: Handle Non-Credentials User
	storedHash, err := base64.StdEncoding.DecodeString(user.Password.String)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}
	storedSalt, err := base64.StdEncoding.DecodeString(user.PasswordSalt.String)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	// Compare Password Hash
	passwordHash := argon2.IDKey([]byte(body.Password), storedSalt, 2, 19*1024, 1, 32)
	if subtle.ConstantTimeCompare(passwordHash, storedHash) == 1 {
		// Create Session in Database
		session, err := s.repo.CreateSession(ctx, user.Id)
		if err != nil {
			util.SendError(s.logger, w, r, err)
			return
		}
		util.SendJson(w, http.StatusOK, &session)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
