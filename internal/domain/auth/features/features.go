package features

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/features/mfa"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/features/session"
	"fbt/backend/internal/domain/auth/features/user"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

type Features struct {
	service service.Service
	d       *util.Dependency
}

func NewFeatures(d *util.Dependency) *Features {
	s := service.NewService(d)
	return &Features{service: s, d: d}
}

func (f *Features) RegisterCredentials(s *grpc.Server) {
	credentials.RegisterService(f.service, credentials.NewRepo(f.d.DB), s)
}

func (f *Features) RegisterMFA(s *grpc.Server) {
	mfa.RegisterService(f.service, mfa.NewRepo(f.d.DB), s)
}

func (f *Features) RegisterOAuth(s *grpc.Server) {
	oauth.RegisterService(f.service, oauth.NewRepo(f.d.DB), s)
}

func (f *Features) RegisterSession(s *grpc.Server) {
	session.RegisterService(f.service, s)
}

func (f *Features) RegisterUser(s *grpc.Server) {
	user.RegisterService(f.service, s)
}
