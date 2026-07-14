package account

import (
	"context"
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

type Repo interface {
	GetAll(ctx context.Context, userID string) (*[]model.Account, error)
	Create(context.Context, *model.Account) (accountID int32, err error)
	Update(context.Context, *model.Account) error
	Delete(ctx context.Context, userID string, accountID int32) error
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo(&repo{db: db})
}

func (s repo) GetAll(ctx context.Context, userID string) (*[]model.Account, error) {
	query := `
		SELECT * FROM accounts
		WHERE user_id = @user_id
	`
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	return util.FetchMany[model.Account](s.db, ctx, query, args)
}

func (s repo) Create(ctx context.Context, account *model.Account) (accountID int32, err error) {
	query := `
		INSERT INTO accounts(name, is_debit, user_id)
		VALUES (@name, @is_debit, @user_id)
		RETURNING account_id
	`
	args := pgx.NamedArgs{
		"name":     account.Name,
		"is_debit": account.IsDebit,
		"user_id":  account.UserId,
	}

	err = s.db.QueryRow(ctx, query, args).Scan(&accountID)

	return
}

func (s *repo) Update(ctx context.Context, account *model.Account) error {
	query := `
		UPDATE accounts SET
			name = @name,
			is_debit = @is_debit
		WHERE account_id = @account_id AND user_id = @user_id
	`
	args := pgx.NamedArgs{
		"account_id": account.ID,
		"name":       account.Name,
		"is_debit":   account.IsDebit,
		"user_id":    account.UserId,
	}

	_, err := s.db.Exec(ctx, query, args)

	return err
}

func (s *repo) Delete(ctx context.Context, userID string, accountID int32) error {
	query := `
		DELETE FROM accounts
		WHERE account_id = @account_id AND user_id = @user_id
	`
	args := pgx.NamedArgs{
		"account_id": accountID,
		"user_id":    userID,
	}

	_, err := s.db.Exec(ctx, query, args)

	return err
}
