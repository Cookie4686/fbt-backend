package auth

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/features/user"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
)

type handler struct {
	credentials *credentials.Feature
	mfa         *mfa.Feature
	oauth       *oauth.Feature
	session     *session.Feature
	user        *user.Feature
}

func newHandler(d *util.Dependency, service service.Service) *handler {
	return &handler{
		credentials: credentials.NewFeature(d, service),
		mfa:         mfa.NewFeature(d, service),
		oauth:       oauth.NewFeature(d, service),
		session:     session.NewFeature(d, service),
		user:        user.NewFeature(d, service),
	}
}
