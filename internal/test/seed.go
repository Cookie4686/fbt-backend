package test

import (
	"bytes"
	"encoding/json"
	"fbt/backend/internal/domain/auth/features/credentials"
	"fbt/backend/internal/domain/auth/model"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func SetupUser(t *testing.T, svr *httptest.Server) (*model.Session, *http.Client) {
	client := svr.Client()

	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(credentials.RegisterPayload{
		Username: "test",
		Email:    "test@email.com",
		Password: "12345678",
	})
	require.NoError(t, err)

	res, err := client.Post(svr.URL+"/credentials/register", "application/json", &payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var session model.Session
	err = json.NewDecoder(res.Body).Decode(&session)
	require.NoError(t, err)

	authCookie := []*http.Cookie{{Name: "session_id", Value: session.Id}}

	baseURL, err := url.Parse(svr.URL)
	require.NoError(t, err)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	jar.SetCookies(baseURL, authCookie)

	client.Jar = jar

	return &session, client
}
