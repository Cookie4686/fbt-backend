package oauth_test

import (
	"testing"

	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/internal/domain/auth/features/oauth"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth(t *testing.T) {
	d := test.Setup(t)

	service := service.NewService(d)
	controller := oauth.NewController(service, oauth.NewRepo(d.DB))

	var registrationID string

	email := "test@email.com"

	t.Run("Login", func(t *testing.T) {
		res, err := controller.Login(t.Context(), &authv1.OAuthServiceLoginRequest{
			Token:    "token",
			Provider: "google",
			Email:    &email,
		})
		require.NoError(t, err)

		require.Equal(t, true, res.RegistrationNeeded)

		registrationID = res.RegistrationId
	})

	var session *authv1.Session

	t.Run("Register", func(t *testing.T) {
		res, err := controller.Register(t.Context(), &authv1.OAuthServiceRegisterRequest{
			Username:        "test",
			Email:           "test@email.com",
			Password:        "12345678",
			PasswordEnabled: true,

			RegistrationId: registrationID,
			TokenId:        "token",
			Provider:       "google",
		})
		require.NoError(t, err)

		assert.Equal(t, false, res.Session.TwoFactorVerified)

		session = res.Session
	})

	t.Run("Status", func(t *testing.T) {
		ctx := test.NewAuthContext(t.Context(), session.Id, service)

		res, err := controller.Status(ctx, &authv1.OAuthServiceStatusRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Providers, 1)
		assert.Contains(t, res.Providers, "google")
	})
}
