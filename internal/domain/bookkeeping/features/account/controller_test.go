package account_test

import (
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"
	"fbt/backend/gen/proto/go/bookkeeping/v1/bookkeepingv1connect"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/test"
	"net/http"
	"slices"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	client := bookkeepingv1connect.NewAccountServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	session := test.SetupUser(t, ctx, baseURL)

	accounts := []model.Account{
		{Name: "Cash", IsDebit: true, UserId: session.UserId},
		{Name: "Bank-1", IsDebit: true, UserId: session.UserId},
		{Name: "Loan-1", IsDebit: false, UserId: session.UserId},
	}

	t.Run("Get All (Empty)", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Account, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		for idx, a := range accounts {
			res, err := client.Create(ctx, &bookkeepingv1.AccountServiceCreateRequest{
				Name:    a.Name,
				IsDebit: a.IsDebit,
			})
			require.NoError(t, err)

			accounts[idx].ID = res.Account.Id
		}
	})

	t.Run("Get All", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, protoToModel(res.Account))
	})

	updated := &accounts[slices.IndexFunc(accounts, func(a model.Account) bool {
		return a.Name == "Bank-1"
	})]
	updated.Name = "Loan-2"
	updated.IsDebit = false

	t.Run("Update", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		_, err := client.Update(ctx, &bookkeepingv1.AccountServiceUpdateRequest{
			Id:      updated.ID,
			Name:    updated.Name,
			IsDebit: updated.IsDebit,
		})
		require.NoError(t, err)

		res, err := client.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		updatedInDB := (res.Account)[slices.IndexFunc(res.Account, func(a *bookkeepingv1.Account) bool {
			return a.Id == updated.ID
		})]

		assert.Equal(t, updated.Name, updatedInDB.Name, "Account Name should be changed")
		assert.Equal(t, updated.IsDebit, updatedInDB.IsDebit, "Account Is Debit should be changed")
		assert.ElementsMatch(t, accounts, protoToModel(res.Account))
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		_, err := client.Delete(ctx, &bookkeepingv1.AccountServiceDeleteRequest{
			Id: updated.ID,
		})
		require.NoError(t, err)

		accounts = slices.DeleteFunc(accounts, func(a model.Account) bool {
			return a.ID == updated.ID
		})

		res, err := client.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, protoToModel(res.Account))
	})
}

func protoToModel(accs []*bookkeepingv1.Account) []model.Account {
	accounts := make([]model.Account, len(accs))
	for idx, a := range accs {
		accounts[idx] = model.Account{
			ID:      a.Id,
			Name:    a.Name,
			IsDebit: a.IsDebit,
			UserId:  a.UserId,
		}
	}

	return accounts
}
