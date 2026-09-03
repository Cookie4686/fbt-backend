package transaction_test

import (
	"testing"

	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/test"

	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransactionController(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)
	session := testUtil.SetupUser(t)
	accounts := testUtil.SetupAccount(t, session.Id)

	controller := transaction.NewController(testUtil.BookkeepingService, transaction.NewRepo(d.DB))

	t.Run("GetAll (empty)", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.Create(ctx, &bookkeepingv1.TransactionServiceCreateRequest{
			Time:        timestamppb.Now(),
			Description: "Loan Cash",
			Note:        "Interest Rate 5%",
			Entries: []*bookkeepingv1.Entry{
				{AccountId: accounts.Cash.ID, Amount: 100},
				{AccountId: accounts.Loan.ID, Amount: 100},
			},
		})
		require.NoError(t, err)

		assert.Equal(t, "Loan Cash", res.TransactionEntry.Description)
		assert.Equal(t, "Interest Rate 5%", res.TransactionEntry.Note)
		assert.Len(t, res.TransactionEntry.Entries, 2)
	})

	t.Run("GetAll", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := controller.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 1)
		assert.Equal(t, "Loan Cash", res.TransactionEntry[0].Description)
		assert.Equal(t, "Interest Rate 5%", res.TransactionEntry[0].Note)
	})
}
