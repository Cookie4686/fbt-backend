package transaction_test

import (
	"testing"
	"time"

	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/test"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepo(t *testing.T) {
	d := test.Setup(t)
	testUtil := test.NewTestUtil(d)
	session := testUtil.SetupUser(t)
	accounts := testUtil.SetupAccount(t, session.Id)

	repo := transaction.NewRepo(d.DB)

	t.Run("GetAll (empty)", func(t *testing.T) {
		ctx := testUtil.NewAuthContext(t.Context(), session.Id)

		res, err := repo.GetAll(ctx, session.UserId)
		require.NoError(t, err)

		assert.Len(t, *res, 0)
	})

	t.Run("Create", func(t *testing.T) {
		_, err := repo.Create(t.Context(), model.Transaction{
			Datetime:    time.Now(),
			Description: pgtype.Text{String: "Loan Cash", Valid: true},
			Note:        pgtype.Text{String: "Interest Rate 5%", Valid: true},
		},
			[]model.Entry{
				{AccountID: accounts.Cash.ID, Amount: 100},
				{AccountID: accounts.Loan.ID, Amount: 100},
			},
		)
		require.NoError(t, err)
	})

	t.Run("GetAll", func(t *testing.T) {
		res, err := repo.GetAll(t.Context(), session.UserId)
		require.NoError(t, err)

		assert.Len(t, *res, 1)
		assert.Equal(t, "Loan Cash", (*res)[0].Description.String)
		assert.Equal(t, "Interest Rate 5%", (*res)[0].Note.String)
	})
}
