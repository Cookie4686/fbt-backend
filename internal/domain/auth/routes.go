package auth

import (
	"fbt/backend/internal/config"
	"fbt/backend/internal/domain/auth/handler"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Routes(logger *zap.Logger, db *pgxpool.Pool, cfg *config.Config, router *mux.Router) {
	handler := handler.NewAuthHandler(logger, db, cfg)

	router.HandleFunc("/validate", handler.Validate).Methods("POST")
	router.HandleFunc("/logout", handler.Logout).Methods("POST")

	router.HandleFunc("/credentials/register", handler.CredentialsRegister).Methods("POST")
	router.HandleFunc("/credentials/login", handler.CredentialsLogin).Methods("POST")

	router.HandleFunc("/oauth/login", handler.OAuthLogin).Methods("POST")
	router.HandleFunc("/oauth/register", handler.OAuthRegister).Methods("POST")
	router.HandleFunc("/oauth/user/{id}", handler.OAuthUserProviders).Methods("GET")

	router.HandleFunc("/mfa/totp", handler.TOTPUpsertKey).Methods("POST")
	router.HandleFunc("/mfa/totp/validate", handler.TOTPValidate).Methods("POST")
	router.HandleFunc("/mfa/users/{id}", handler.GetUserMFAList).Methods("GET")

	router.HandleFunc("/users/{username}", handler.GetUser).Methods("GET")
}
