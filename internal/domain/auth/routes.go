package auth

import (
	"fbt/backend/internal/domain/auth/handler"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Routes(logger *zap.Logger, db *pgxpool.Pool, router *mux.Router) {
	handler := handler.NewAuthHandler(logger, db)

	router.HandleFunc("/users/{username}", handler.GetUser).Methods("GET")

	router.HandleFunc("/credentials/register", handler.CredentialsRegister).Methods("POST")

	router.HandleFunc("/credentials/login", handler.CredentialsLogin).Methods("POST")

	router.HandleFunc("/logout", handler.Logout).Methods("POST")

	router.HandleFunc("/validate", handler.Validate).Methods("POST")

	router.HandleFunc("/oauth/login", handler.OAuthLogin).Methods("POST")

	router.HandleFunc("/oauth/register", handler.OAuthRegister).Methods("POSt")
}
