package credentials_test

import (
	"fbt/backend/internal/domain/auth/features/credentials/pb"
	"fbt/backend/internal/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	ctx, ctr, conn := test.NewTestConnection(t, 1234)

	client := pb.NewCredentialsClient(conn)

	username := "test"
	password := "12345678"

	t.Run("Register", func(t *testing.T) {
		res, err := client.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Email:    "test@email.com",
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})

	t.Run("Login", func(t *testing.T) {
		res, err := client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})
		require.NoError(t, err)

		require.Equal(t, false, res.Session.TwoFactorVerified)
	})

	ctr.Restore(ctx)
}
