package test

import (
	"context"
	"testing"

	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	authService "fbt/backend/internal/domain/auth/service"

	"github.com/stretchr/testify/require"
)

type TestUtil struct {
	AuthService           authService.Service
	CredentialsController credentials.Controller
}

func NewTestUtil(d *util.Dependency) *TestUtil {
	authService := authService.NewService(d)

	return &TestUtil{
		AuthService:           authService,
		CredentialsController: *credentials.NewController(authService, credentials.NewRepo(d.DB)),
	}
}

func (s *TestUtil) NewAuthContext(ctx context.Context, sessionID string) context.Context {
	auth, err := s.AuthService.Validate(ctx, sessionID)
	if err != nil {
		return ctx
	}

	return interceptor.NewAuthContext(ctx, auth)
}

func (s *TestUtil) SetupUser(t *testing.T) *authv1.Session {
	res, err := s.CredentialsController.Register(t.Context(), &authv1.CredentialServiceRegisterRequest{
		Username: "test",
		Email:    "test@email.com",
		Password: "12345678",
	})
	require.NoError(t, err)

	return res.Session
}
