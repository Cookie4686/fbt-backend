package transaction

import (
	"context"
	"fbt/backend/internal/domain/bookkeeping/model"
	"log"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

type Repo interface {
	GetAll(ctx context.Context, userID string) (*[]model.TransactionEntry, error)
	Create(context.Context, *model.TransactionEntry) (transactionID int64, err error)
	Update(context.Context, *model.TransactionEntry) error
	Delete(ctx context.Context, userID string, transactionID int64) error
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db: db})
}

type GetAllRow struct {
	model.Transaction
	model.Entry
}

func (s repo) GetAll(ctx context.Context, userID string) (*[]model.TransactionEntry, error) {
	query := `
		SELECT transactions.transaction_id, transactions.datetime, entries.account_id, entries.amount
		FROM accounts
		JOIN entries ON entries.account_id = accounts.account_id
		JOIN transactions ON transactions.transaction_id = entries.transaction_id
		WHERE accounts.user_id = @user_id
	`
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var te []model.TransactionEntry

	for rows.Next() {
		var i GetAllRow
		if err := rows.Scan(
			&i.TransactionID,
			&i.Datetime,
			&i.AccountID,
			&i.Amount,
		); err != nil {
			return nil, err
		}

		if idx := slices.IndexFunc(te, func(t model.TransactionEntry) bool {
			return t.TransactionID == i.TransactionID
		}); idx == -1 {
			te = append(te, model.TransactionEntry{
				Transaction: model.Transaction{
					TransactionID: i.TransactionID,
					Datetime:      i.Datetime,
				},
				Entries: []model.Entry{
					{AccountID: i.AccountID, Amount: i.Amount},
				},
			})
		} else {
			te[idx].Entries = append(
				te[idx].Entries,
				model.Entry{AccountID: i.AccountID, Amount: i.Amount},
			)
		}
	}

	return &te, nil
}

func (s repo) Create(ctx context.Context, te *model.TransactionEntry) (transactionID int64, err error) {
	log.Print(te.Datetime)

	query := `
		INSERT INTO transactions(datetime)
		VALUES (@datetime)
		RETURNING transaction_id
	`
	args := pgx.NamedArgs{"datetime": te.Datetime}

	err = s.db.QueryRow(ctx, query, args).Scan(&transactionID)
	if err != nil {
		return transactionID, err
	}

	_, err = s.createEntries(ctx, transactionID, te.Entries)
	if err != nil {
		return transactionID, err
	}

	return transactionID, nil
}

func (s *repo) Update(ctx context.Context, te *model.TransactionEntry) error {
	batch := &pgx.Batch{}
	batch.Queue(`
		UPDATE transactions
		SET datetime = @datetime
		WHERE transaction_id = @transaction_id
	`,
		pgx.NamedArgs{"@transaction_id": te.TransactionID},
	)
	batch.Queue(`
		DELETE entries
		WHERE transaction_id = @transaction_id
	`,
		pgx.NamedArgs{"@transaction_id": te.TransactionID},
	)

	_, err := s.db.SendBatch(ctx, batch).Exec()
	if err != nil {
		return err
	}

	_, err = s.createEntries(ctx, te.TransactionID, te.Entries)
	if err != nil {
		return err
	}

	return nil
}

func (s *repo) Delete(ctx context.Context, userID string, transactionID int64) error {
	query := `
		DELETE FROM transactions
		WHERE transaction_id = @transaction_id
	`
	args := pgx.NamedArgs{"transaction_id": transactionID}

	_, err := s.db.Exec(ctx, query, args)

	return err
}

func (s *repo) createEntries(ctx context.Context, transactionID int64, entries []model.Entry) (count int64, err error) {
	_, err = s.db.CopyFrom(
		ctx,
		pgx.Identifier{"entries"},
		[]string{"transaction_id", "account_id", "amount"},
		pgx.CopyFromSlice(len(entries), func(i int) ([]any, error) {
			return []any{transactionID, entries[i].AccountID, entries[i].Amount}, nil
		}),
	)

	return
}
