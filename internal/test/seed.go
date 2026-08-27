package test

import (
	"testing"

	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/internal/domain/auth/features/credentials"
	authService "fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"

	"github.com/stretchr/testify/require"
)

func SetupUser(t *testing.T, d *util.Dependency, service authService.Service) *authv1.Session {
	if service == nil {
		service = authService.NewService(d)
	}

	controller := credentials.NewController(service, credentials.NewRepo(d.DB))

	res, err := controller.Register(t.Context(), &authv1.CredentialServiceRegisterRequest{
		Username: "test",
		Email:    "test@email.com",
		Password: "12345678",
	})
	require.NoError(t, err)

	return res.Session
}
