package auth

import (
	"fbt/backend/internal/domain/auth/features"
	"fbt/backend/internal/domain/auth/middleware"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(d *util.Dependency, r *mux.Router) middleware.Middleware {
	f, m := features.NewFeatures(d)

	r.Handle("/credentials/register", f.Credentials.Register).Methods(http.MethodPost)
	r.Handle("/credentials/login", f.Credentials.Login).Methods(http.MethodPost)

	r.Handle("/oauth/register", f.OAuth.Register).Methods(http.MethodPost)
	r.Handle("/oauth/login", f.OAuth.Login).Methods(http.MethodPost)
	r.Handle("/oauth/status", f.OAuth.AUTH_Status).Methods(http.MethodGet)

	r.Handle("/mfa/totp", f.MFA.AUTH_TOTPUpsertKey).Methods(http.MethodPost)
	r.Handle("/mfa/totp/validate", f.MFA.AUTH_TOTPValidate).Methods(http.MethodPost)
	r.Handle("/mfa/status", f.MFA.AUTH_MFAStatus).Methods(http.MethodGet)

	r.Handle("/validate", f.Session.AUTH_Validate).Methods(http.MethodPost)
	r.Handle("/logout", f.Session.AUTH_Logout).Methods(http.MethodPost)

	r.Handle("/users/{username}", f.User.GetByUsername).Methods(http.MethodGet)

	return m
}
