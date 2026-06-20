package auth

import (
	"fbt/backend/internal/domain/auth/features"
	"fbt/backend/internal/util"

	"github.com/gorilla/mux"
)

func Routes(d *util.Dependency, r *mux.Router) {
	f := features.NewFeatures(d)

	r.Handle("/credentials/register", f.Credentials.Register).Methods("POST")
	r.Handle("/credentials/login", f.Credentials.Login).Methods("POST")

	r.Handle("/oauth/register", f.OAuth.Register).Methods("POST")
	r.Handle("/oauth/login", f.OAuth.Login).Methods("POST")
	r.Handle("/oauth/status", f.OAuth.AUTH_Status).Methods("GET")

	r.Handle("/mfa/totp", f.MFA.AUTH_TOTPUpsertKey).Methods("POST")
	r.Handle("/mfa/totp/validate", f.MFA.AUTH_TOTPValidate).Methods("POST")
	r.Handle("/mfa/status", f.MFA.AUTH_MFAStatus).Methods("GET")

	r.Handle("/validate", f.Session.AUTH_Validate).Methods("POST")
	r.Handle("/logout", f.Session.AUTH_Logout).Methods("POST")

	r.Handle("/users/{username}", f.User.GetByUsername).Methods("GET")
}
