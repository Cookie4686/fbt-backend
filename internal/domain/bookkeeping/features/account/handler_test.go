package account_test

import (
	"bytes"
	"encoding/json"
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/test"
	"net/http"
	"net/url"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	_, _, svr := test.NewTestAPI(t)
	defer svr.Close()

	session, client := test.SetupUser(t, svr)

	url, err := url.ParseRequestURI(svr.URL + "/accounts")
	require.NoError(t, err)

	accounts := []model.Account{
		{Name: "Cash", IsDebit: true, UserId: session.UserId},
		{Name: "Bank-1", IsDebit: true, UserId: session.UserId},
		{Name: "Loan-1", IsDebit: false, UserId: session.UserId},
	}

	t.Run("Get All (Empty)", func(t *testing.T) {
		res, err := client.Get(url.String())
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		var accounts []model.Account
		err = json.NewDecoder(res.Body).Decode(&accounts)
		require.NoError(t, err)

		assert.Len(t, accounts, 0)
	})

	t.Run("Create", func(t *testing.T) {
		for idx, a := range accounts {
			var payload bytes.Buffer
			err := json.NewEncoder(&payload).Encode(account.CreatePayload{
				Name:    a.Name,
				IsDebit: a.IsDebit,
			})
			require.NoError(t, err)
			res, err := client.Post(url.String(), "application/json", &payload)
			require.NoError(t, err)

			var ac model.Account
			err = json.NewDecoder(res.Body).Decode(&ac)
			require.NoError(t, err)

			accounts[idx].ID = ac.ID
		}
	})

	t.Run("Get All", func(t *testing.T) {
		res, err := client.Get(url.String())
		require.NoError(t, err)

		var accs []model.Account
		err = json.NewDecoder(res.Body).Decode(&accs)
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, accs)
	})

	updated := &accounts[slices.IndexFunc(accounts, func(a model.Account) bool {
		return a.Name == "Bank-1"
	})]
	updated.Name = "Loan-2"
	updated.IsDebit = false

	t.Run("Update", func(t *testing.T) {
		var payload bytes.Buffer
		err := json.NewEncoder(&payload).Encode(account.UpdatePayload{
			ID:      updated.ID,
			Name:    updated.Name,
			IsDebit: updated.IsDebit,
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, url.String(), &payload)
		require.NoError(t, err)
		res, err := client.Do(req)
		require.NoError(t, err)

		res, err = client.Get(url.String())
		require.NoError(t, err)

		var accsInDB []model.Account
		err = json.NewDecoder(res.Body).Decode(&accsInDB)
		require.NoError(t, err)

		updatedInDB := (accsInDB)[slices.IndexFunc(accsInDB, func(a model.Account) bool {
			return a.ID == updated.ID
		})]

		assert.Equal(t, updated.Name, updatedInDB.Name, "Account Name should be changed")
		assert.Equal(t, updated.IsDebit, updatedInDB.IsDebit, "Account Is Debit should be changed")
		assert.ElementsMatch(t, accounts, accsInDB)
	})

	t.Run("Delete", func(t *testing.T) {
		var payload bytes.Buffer
		err := json.NewEncoder(&payload).Encode(account.DeletePayload{
			ID: updated.ID,
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, url.String(), &payload)
		require.NoError(t, err)
		res, err := client.Do(req)
		require.NoError(t, err)

		accounts = slices.DeleteFunc(accounts, func(a model.Account) bool {
			return a.ID == updated.ID
		})

		res, err = client.Get(url.String())
		require.NoError(t, err)

		var accsInDB []model.Account
		err = json.NewDecoder(res.Body).Decode(&accsInDB)
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, accsInDB)
	})
}
