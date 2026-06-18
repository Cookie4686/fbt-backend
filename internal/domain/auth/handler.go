package auth

import (
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/features/user"
	"fbt/backend/internal/domain/auth/service"
)

type handler struct {
	credentials *credentials.Feature
	mfa         *mfa.Feature
	oauth       *oauth.Feature
	session     *session.Feature
	user        *user.Feature
}

func newHandler(d *dependency.Dependency, service *service.AuthService) *handler {
	return &handler{
		credentials: credentials.NewFeature(d, service),
		mfa:         mfa.NewFeature(d, service),
		oauth:       oauth.NewFeature(d, service),
		session:     session.NewFeature(d, service),
		user:        user.NewFeature(d, service),
	}
}
