package credentials_test

import (
	"testing"

	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/test"

	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	d := test.Setup(t)

	controller := credentials.NewController(service.NewService(d), credentials.NewRepo(d.DB))

	username := "test"
	password := "12345678"

	t.Run("Register", func(t *testing.T) {
		ctx := t.Context()

		res, err := controller.Register(ctx, &authv1.CredentialServiceRegisterRequest{
			Username: username,
			Email:    "test@email.com",
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})

	t.Run("Login", func(t *testing.T) {
		ctx := t.Context()

		res, err := controller.Login(ctx, &authv1.CredentialServiceLoginRequest{
			Username: username,
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})
}
