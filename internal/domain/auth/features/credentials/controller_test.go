package credentials_test

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/test/mock"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	s := mock.NewMockService()
	r := mock.NewCredentialRepo(s)
	controller := credentials.NewController(s, r)

	username := "test"
	password := "12345678"

	t.Run("Credential Register Controller", func(t *testing.T) {
		payload := credentials.RegisterPayload{
			Username: username,
			Email:    "test@email.com",
			Password: password,
		}

		res, err := controller.Register(t.Context(), &payload)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		session := s.Sessions[slices.IndexFunc(s.Sessions, func(s model.Session) bool { return s.Id == res.Payload.Id })]
		require.Equal(t, session, *res.Payload)
		require.Equal(t, false, res.Payload.TwoFactorVerified)

		user := s.Users[slices.IndexFunc(s.Users, func(user model.User) bool { return user.Username == payload.Username })]
		require.Equal(t, payload.Username, user.Username)
		require.Equal(t, payload.Email, user.Email)
		require.Equal(t, false, user.EmailVerified)
		require.NotEqual(t, payload.Password, user.Password)
	})

	t.Run("Credential Login Controller", func(t *testing.T) {
		payload := credentials.LoginPayload{
			Username: username,
			Password: password,
		}

		res, err := controller.Login(t.Context(), &payload)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		session := s.Sessions[slices.IndexFunc(s.Sessions, func(s model.Session) bool { return s.Id == res.Payload.Id })]
		require.Equal(t, session, *res.Payload)
		require.Equal(t, false, res.Payload.TwoFactorVerified)
	})
}
