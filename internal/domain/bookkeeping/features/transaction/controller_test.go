package transaction_test

import (
	"testing"

	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/test"

	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransaction(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)
	session := testUtil.SetupUser(t)

	service := service.NewService(d)
	accController := account.NewController(service, account.NewRepo(d.DB))
	controller := transaction.NewController(service, transaction.NewRepo(d.DB))

	cash := model.Account{Name: "Cash", IsDebit: true}
	loan := model.Account{Name: "Loan", IsDebit: false}

	for _, a := range []*model.Account{&cash, &loan} {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := accController.Create(ctx, &bookkeepingv1.AccountServiceCreateRequest{
			Name:    a.Name,
			IsDebit: a.IsDebit,
		})
		require.NoError(t, err)

		a.ID = res.Account.Id
	}

	t.Run("GetAll (empty)", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.Create(ctx, &bookkeepingv1.TransactionServiceCreateRequest{
			Time: timestamppb.Now(),
			Entries: []*bookkeepingv1.Entry{
				{AccountId: cash.ID, Amount: 100},
				{AccountId: loan.ID, Amount: 100},
			},
		})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry.Entries, 2)
	})

	t.Run("GetAll", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 1)
	})
}
