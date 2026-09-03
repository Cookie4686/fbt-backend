package account_test

import (
	"slices"
	"testing"

	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/test"

	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)

	controller := account.NewController(service.NewService(d), account.NewRepo(d.DB))

	session := testUtil.SetupUser(t)

	accounts := []model.Account{
		{Name: "Cash", IsDebit: true, UserID: session.UserId},
		{Name: "Bank-1", IsDebit: true, UserID: session.UserId},
		{Name: "Loan-1", IsDebit: false, UserID: session.UserId},
	}

	t.Run("Get All (Empty)", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Account, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		for idx, a := range accounts {
			res, err := controller.Create(ctx, &bookkeepingv1.AccountServiceCreateRequest{
				Name:    a.Name,
				IsDebit: a.IsDebit,
			})
			require.NoError(t, err)

			accounts[idx].ID = res.Account.Id
		}
	})

	t.Run("Get All", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, protoToModel(res.Account))
	})

	updated := &accounts[slices.IndexFunc(accounts, func(a model.Account) bool {
		return a.Name == "Bank-1"
	})]
	updated.Name = "Loan-2"
	updated.IsDebit = false

	t.Run("Update", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		_, err := controller.Update(ctx, &bookkeepingv1.AccountServiceUpdateRequest{
			Id:      updated.ID,
			Name:    updated.Name,
			IsDebit: updated.IsDebit,
		})
		require.NoError(t, err)

		res, err := controller.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
		require.NoError(t, err)

		updatedInDB := res.Account[slices.IndexFunc(res.Account, func(a *bookkeepingv1.Account) bool {
			return a.Id == updated.ID
		})]

		assert.Equal(t, updated.Name, updatedInDB.Name, "Account Name should be changed")
		assert.Equal(t, updated.IsDebit, updatedInDB.IsDebit, "Account Is Debit should be changed")
		assert.ElementsMatch(t, accounts, protoToModel(res.Account))
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		_, err := controller.Delete(ctx, &bookkeepingv1.AccountServiceDeleteRequest{
			Id: updated.ID,
		})
		require.NoError(t, err)

		accounts = slices.DeleteFunc(accounts, func(a model.Account) bool {
			return a.ID == updated.ID
		})

		res, err := controller.GetAll(ctx, &bookkeepingv1.AccountServiceGetAllRequest{})
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
			UserID:  a.UserId,
		}
	}

	return accounts
}
