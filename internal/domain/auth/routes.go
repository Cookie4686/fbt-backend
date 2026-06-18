package auth

import (
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth/service"

	"github.com/gorilla/mux"
)

func Routes(d *dependency.Dependency, r *mux.Router) {
	service := service.NewAuthService(d)
	middleware := newMiddleware(d, service)
	AUTH := middleware.Auth

	h := newHandler(d, service)

	r.Handle("/validate", AUTH(h.session.Validate)).Methods("POST")
	r.Handle("/logout", AUTH(h.session.Logout)).Methods("POST")

	r.Handle("/credentials/register", h.credentials.Register).Methods("POST")
	r.Handle("/credentials/login", h.credentials.Login).Methods("POST")

	r.Handle("/oauth/register", h.oauth.Register).Methods("POST")
	r.Handle("/oauth/login", h.oauth.Login).Methods("POST")
	r.Handle("/oauth/status", AUTH(h.oauth.Status)).Methods("GET")

	r.Handle("/mfa/totp", AUTH(h.mfa.TOTPUpsertKey)).Methods("POST")
	r.Handle("/mfa/totp/validate", AUTH(h.mfa.TOTPValidate)).Methods("POST")
	r.Handle("/mfa/status", AUTH(h.mfa.MFAStatus)).Methods("GET")

	r.Handle("/users/{username}", h.user.GetByUsername).Methods("GET")
}
