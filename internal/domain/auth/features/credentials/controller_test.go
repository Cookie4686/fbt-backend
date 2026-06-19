package credentials_test

import (
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/test/mock"
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	service    = mock.NewMockService()
	repo       = mock.NewCredentialRepo(service)
	controller = credentials.NewController(service, repo)
)

func TestRegister(t *testing.T) {
	payload := credentials.RegisterPayload{
		Username: mock.ValidUser.Username,
		Email:    "test@email.com",
		Password: mock.ValidUser.Password,
	}

	res, err := controller.Register(t.Context(), &payload)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)

	session := service.Sessions[slices.IndexFunc(service.Sessions, func(s model.Session) bool { return s.Id == res.Payload.Id })]
	assert.Equal(t, session, *res.Payload, "Session should match session in the database")
	assert.Equal(t, false, res.Payload.TwoFactorVerified, "Session should not be 2FA verified")

	user := service.Users[slices.IndexFunc(service.Users, func(user model.User) bool { return user.Username == payload.Username })]
	assert.Equal(t, payload.Username, user.Username, "User Username should match Input")
	assert.Equal(t, payload.Email, user.Email, "User Email should match Input")
	assert.Equal(t, false, user.EmailVerified, "User Email should not be verified")
	assert.NotEqual(t, payload.Password, user.Password, "User Password should not match Input")
}

func TestLogin(t *testing.T) {
	t.Run("Credential Login Controller", func(t *testing.T) {
		payload := credentials.LoginPayload{
			Username: mock.ValidUser.Username,
			Password: mock.ValidUser.Password,
		}

		res, err := controller.Login(t.Context(), &payload)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		session := service.Sessions[slices.IndexFunc(service.Sessions, func(s model.Session) bool { return s.Id == res.Payload.Id })]
		assert.Equal(t, session, *res.Payload, "Session should match session in the database")
		assert.Equal(t, false, res.Payload.TwoFactorVerified, "Session should not be 2FA verified")
	})
}
