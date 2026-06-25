package account_test

import (
	"fbt/backend/internal/domain/bookkeeping/features/account/pb"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/test"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestAccount(t *testing.T) {
	ctx, conn := test.NewTestLocalAPI(t)

	client := pb.NewAccountServiceClient(conn)

	session := test.SetupUser(t, ctx, conn)

	accounts := []model.Account{
		{Name: "Cash", IsDebit: true, UserId: session.UserID},
		{Name: "Bank-1", IsDebit: true, UserId: session.UserID},
		{Name: "Loan-1", IsDebit: false, UserId: session.UserID},
	}

	t.Run("Get All (Empty)", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		assert.Len(t, res.Account, 0)
	})

	t.Run("Create", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		for idx, a := range accounts {
			res, err := client.Create(ctx, &pb.CreateRequest{
				Name:    a.Name,
				IsDebit: a.IsDebit,
			})
			require.NoError(t, err)

			accounts[idx].ID = res.Account.Id
		}
	})

	t.Run("Get All", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, commonToAccount(res.Account))
	})

	updated := &accounts[slices.IndexFunc(accounts, func(a model.Account) bool {
		return a.Name == "Bank-1"
	})]
	updated.Name = "Loan-2"
	updated.IsDebit = false

	t.Run("Update", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		_, err := client.Update(ctx, &pb.UpdateRequest{
			Id:      updated.ID,
			Name:    updated.Name,
			IsDebit: updated.IsDebit,
		})
		require.NoError(t, err)

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		updatedInDB := (res.Account)[slices.IndexFunc(res.Account, func(a *pb.Account) bool {
			return a.Id == updated.ID
		})]

		assert.Equal(t, updated.Name, updatedInDB.Name, "Account Name should be changed")
		assert.Equal(t, updated.IsDebit, updatedInDB.IsDebit, "Account Is Debit should be changed")
		assert.ElementsMatch(t, accounts, commonToAccount(res.Account))
	})

	t.Run("Delete", func(t *testing.T) {
		ctx = metadata.AppendToOutgoingContext(t.Context(), "session_id", session.Id)

		_, err := client.Delete(ctx, &pb.DeleteRequest{
			Id: updated.ID,
		})
		require.NoError(t, err)

		accounts = slices.DeleteFunc(accounts, func(a model.Account) bool {
			return a.ID == updated.ID
		})

		res, err := client.GetAll(ctx, &pb.GetAllRequest{})
		require.NoError(t, err)

		assert.ElementsMatch(t, accounts, commonToAccount(res.Account))
	})
}

func commonToAccount(accs []*pb.Account) []model.Account {
	accounts := make([]model.Account, len(accs))
	for idx, a := range accs {
		accounts[idx] = model.Account{
			ID:      a.Id,
			Name:    a.Name,
			IsDebit: a.IsDebit,
			UserId:  a.UserID,
		}
	}
	return accounts
}
