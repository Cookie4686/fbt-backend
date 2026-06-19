package credentials_test

import (
	"bytes"
	"encoding/json"
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialsHandler(t *testing.T) {
	ctx, ctr, svr := test.NewTestAPI(t)
	defer svr.Close()

	client := svr.Client()

	username := "test"
	password := "12345678"

	t.Run("Register", func(t *testing.T) {
		var payload bytes.Buffer
		err := json.NewEncoder(&payload).Encode(credentials.RegisterPayload{
			Username: username,
			Email:    "test@email.com",
			Password: password,
		})
		require.NoError(t, err)

		res, err := client.Post(svr.URL+"/credentials/register", "application/json", &payload)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		var body model.Session
		err = json.NewDecoder(res.Body).Decode(&body)
		require.NoError(t, err)

		require.Equal(t, false, body.TwoFactorVerified)
	})

	t.Run("Login", func(t *testing.T) {
		var payload bytes.Buffer
		err := json.NewEncoder(&payload).Encode(credentials.LoginPayload{
			Username: username,
			Password: password,
		})

		require.NoError(t, err)

		res, err := client.Post(svr.URL+"/credentials/login", "application/json", &payload)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)

		var body model.Session
		err = json.NewDecoder(res.Body).Decode(&body)
		require.NoError(t, err)

		require.Equal(t, false, body.TwoFactorVerified)
	})

	ctr.Restore(ctx)
}
