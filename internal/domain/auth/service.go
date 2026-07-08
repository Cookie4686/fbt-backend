package auth

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/features/email"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/features/webauthn"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"

	"connectrpc.com/connect"
)

func RegisterService(mux *http.ServeMux, d *util.Dependency, opts ...connect.HandlerOption) *http.ServeMux {
	s := service.NewService(d)

	mux.Handle(credentials.NewServiceHandler(s, credentials.NewRepo(d.DB), opts...))
	mux.Handle(email.NewServiceHandler(s, email.NewRepo(d.DB), opts...))
	mux.Handle(mfa.NewServiceHandler(s, mfa.NewRepo(d.DB), opts...))
	mux.Handle(oauth.NewServiceHandler(s, oauth.NewRepo(d.DB), opts...))
	mux.Handle(session.NewServiceHandler(s, opts...))
	mux.Handle(webauthn.NewServiceHandler(s, webauthn.NewRepo(d.DB), opts...))

	return mux
}
