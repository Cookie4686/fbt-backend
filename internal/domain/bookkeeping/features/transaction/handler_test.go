package transaction_test

import (
	"bytes"
	"encoding/json"
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/test"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepo(t *testing.T) {
	_, _, svr := test.NewTestAPI(t)
	defer svr.Close()

	_, client := test.SetupUser(t, svr)

	accountURL, err := url.ParseRequestURI(svr.URL + "/accounts")
	require.NoError(t, err)

	cash := model.Account{Name: "Cash", IsDebit: true}
	loan := model.Account{Name: "Loan", IsDebit: false}

	for _, a := range []*model.Account{&cash, &loan} {
		var payload bytes.Buffer
		err := json.NewEncoder(&payload).Encode(account.CreatePayload{
			Name:    a.Name,
			IsDebit: a.IsDebit,
		})
		require.NoError(t, err)

		res, err := client.Post(accountURL.String(), "application/json", &payload)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		var ac model.Account
		err = json.NewDecoder(res.Body).Decode(&ac)
		require.NoError(t, err)

		a.ID = ac.ID
	}

	transactionURL, err := url.ParseRequestURI(svr.URL + "/transactions")
	require.NoError(t, err)

	t.Run("GetAll (empty)", func(t *testing.T) {
		res, err := client.Get(transactionURL.String())
		require.NoError(t, err)

		var te []model.TransactionEntry
		err = json.NewDecoder(res.Body).Decode(&te)
		require.NoError(t, err)

		assert.Len(t, te, 0)
	})

	t.Run("Create", func(t *testing.T) {
		var payload bytes.Buffer
		err = json.NewEncoder(&payload).Encode(transaction.CreatePayload{
			Datetime: time.Now().Format(time.RFC3339),
			Entries: []transaction.EntryPayload{
				{AccountID: cash.ID, Amount: 100},
				{AccountID: loan.ID, Amount: 100},
			},
		})
		require.NoError(t, err)

		res, err := client.Post(transactionURL.String(), "application/json", &payload)
		require.NoError(t, err)

		var te model.TransactionEntry
		err = json.NewDecoder(res.Body).Decode(&te)
		require.NoError(t, err)

		assert.Len(t, te.Entries, 2)
	})

	t.Run("GetAll", func(t *testing.T) {
		res, err := client.Get(transactionURL.String())
		require.NoError(t, err)

		var te []model.TransactionEntry
		err = json.NewDecoder(res.Body).Decode(&te)
		require.NoError(t, err)

		assert.Len(t, te, 1)
	})
}
