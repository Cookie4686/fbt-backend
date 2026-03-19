package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

func Routes(logger *zap.Logger, db *pgxpool.Pool, router *mux.Router) {
	auth := NewAuthRepository(db)

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		params := mux.Vars(r)
		id := params["id"]

		user, err := auth.GetUserByUsername(ctx, id)
		if err != nil {
			util.SendError(logger, w, r, err)
			return
		}

		util.SendJson(w, http.StatusOK, &user)
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		body := struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil {
			util.SendError(logger, w, r, err)
			return
		}

		salt := make([]byte, 16)
		rand.Read(salt)
		passwordHash := argon2.IDKey([]byte(body.Password), salt, 2, 19*1024, 1, 32)

		user := &User{
			Id:            util.GenerateBase32UUID(),
			Username:      body.Username,
			Email:         body.Email,
			EmailVerified: false,
			Password:      base64.StdEncoding.EncodeToString(passwordHash),
			PasswordSalt:  base64.StdEncoding.EncodeToString(salt),
		}
		session := &Session{
			Id:                util.GenerateBase64UUID(),
			UserId:            user.Id,
			ExpiresAt:         time.Now().Add(sessionExpiresIn),
			TwoFactorVerified: false,
		}
		err := auth.CredentialsRegister(ctx, user, session)
		if err != nil {
			util.SendError(logger, w, r, err)
			return
		}

		util.SendJson(w, http.StatusOK, &session)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil {
			util.SendError(logger, w, r, err)
			return
		}

		// Get User Data From Database
		user, err := auth.GetUserByUsername(ctx, body.Username)
		if err != nil {
			util.SendError(logger, w, r, err)
			return
		}
		storedHash, err := base64.StdEncoding.DecodeString(user.Password)
		if err != nil {
			util.SendError(logger, w, r, err)
			return
		}
		storedSalt, err := base64.StdEncoding.DecodeString(user.PasswordSalt)
		if err != nil {
			util.SendError(logger, w, r, err)
			return
		}

		passwordHash := argon2.IDKey([]byte(body.Password), storedSalt, 2, 19*1024, 1, 32)
		logger.Debug("/login", zap.ByteString("passwordHash", passwordHash), zap.ByteString("storedHash", storedHash))
		if subtle.ConstantTimeCompare(passwordHash, storedHash) == 1 {
			// Create Session in Database
			session, err := auth.CreateSession(ctx, user.Id)
			if err != nil {
				util.SendError(logger, w, r, err)
				return
			}
			util.SendJson(w, http.StatusOK, &session)
		}

		w.WriteHeader(http.StatusUnauthorized)
	}).Methods("POST")
}
