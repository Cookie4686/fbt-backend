package features

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/features/user"
	"fbt/backend/internal/domain/auth/middleware"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
)

type features struct {
	Credentials *credentials.Feature
	MFA         *mfa.Feature
	OAuth       *oauth.Feature
	Session     *session.Feature
	User        *user.Feature
}

func NewFeatures(d *util.Dependency) (*features, middleware.Middleware) {
	s := service.NewService(d)
	m := middleware.NewMiddleware(d, s)

	return &features{
		Credentials: credentials.NewFeature(d, s),
		MFA:         mfa.NewFeature(d, s, m),
		OAuth:       oauth.NewFeature(d, s, m),
		Session:     session.NewFeature(d, s, m),
		User:        user.NewFeature(d, s),
	}, m
}
