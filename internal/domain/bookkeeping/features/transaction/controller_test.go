package transaction_test

import (
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"
	"fbt/backend/gen/proto/go/bookkeeping/v1/bookkeepingv1connect"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/test"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransaction(t *testing.T) {
	ctx, baseURL := test.NewTestLocalAPI(t)

	accClient := bookkeepingv1connect.NewAccountServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())
	client := bookkeepingv1connect.NewTransactionServiceClient(http.DefaultClient, baseURL, connect.WithGRPC())

	session := test.SetupUser(t, ctx, baseURL)

	cash := model.Account{Name: "Cash", IsDebit: true}
	loan := model.Account{Name: "Loan", IsDebit: false}

	for _, a := range []*model.Account{&cash, &loan} {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := accClient.Create(ctx, &bookkeepingv1.AccountServiceCreateRequest{
			Name:    a.Name,
			IsDebit: a.IsDebit,
		})
		require.NoError(t, err)

		a.ID = res.Account.Id
	}

	t.Run("GetAll (empty)", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.Create(ctx, &bookkeepingv1.TransactionServiceCreateRequest{
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
		ctx := interceptor.NewTokenContext(t.Context(), session.Id)

		res, err := client.GetAll(ctx, &bookkeepingv1.TransactionServiceGetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 1)
	})
}
