package test

import (
	"testing"

	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/util"

	authv1 "fbt/backend/gen/proto/go/auth/v1"

	authService "fbt/backend/internal/domain/auth/service"

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
