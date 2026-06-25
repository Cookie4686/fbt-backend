package transaction_test

import (
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/test"
	"testing"

	accountPB "fbt/backend/internal/domain/bookkeeping/features/account/pb"
	"fbt/backend/internal/domain/bookkeeping/features/transaction/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransaction(t *testing.T) {
	ctx, conn := test.NewTestLocalAPI(t)

	accClient := accountPB.NewAccountServiceClient(conn)
	client := pb.NewTransactionClient(conn)

	session := test.SetupUser(t, ctx, conn)

	cash := model.Account{Name: "Cash", IsDebit: true}
	loan := model.Account{Name: "Loan", IsDebit: false}

	for _, a := range []*model.Account{&cash, &loan} {
		ctx = metadata.AppendToOutgoingContext(ctx, "session_id", session.Id)

		res, err := accClient.Create(ctx, &accountPB.CreateRequest{
			Name:    a.Name,
			IsDebit: a.IsDebit,
		})
		require.NoError(t, err)

		a.ID = res.Account.Id
	}

	t.Run("GetAll (empty)", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.Create(ctx, &pb.CreateRequest{
			Time: timestamppb.Now(),
			Entries: []*pb.Entry{
				{AccountID: cash.ID, Amount: 100},
				{AccountID: loan.ID, Amount: 100},
			},
		})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry.Entries, 2)
	})

	t.Run("GetAll", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.TransactionEntry, 1)
	})
}
