package model

import (
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Account struct {
	ID      int32  `json:"id" db:"account_id"`
	Name    string `json:"name" db:"name"`
	IsDebit bool   `json:"is_debit" db:"is_debit"`
	UserId  string `json:"user_id" db:"user_id"`
}

func (a *Account) ToProto() *bookkeepingv1.Account {
	return &bookkeepingv1.Account{
		Id:      a.ID,
		Name:    a.Name,
		IsDebit: a.IsDebit,
		UserId:  a.UserId,
	}
}

type Tag struct {
	TagID  int64  `json:"tag_id" db:"tag_id"`
	Name   string `json:"name" db:"name"`
	UserID string `json:"user_id" db:"user_id"`
}

type AccountTag struct {
	Account
	Tags []Tag `json:"tags"`
}

type Transaction struct {
	TransactionID int64     `json:"transaction_id" db:"transaction_id"`
	Datetime      time.Time `json:"datetime" db:"datetime"`
}

type Entry struct {
	// TransactionID int64   `json:"transaction_id" db:"transaction_id"`
	AccountID int32   `json:"account_id" db:"account_id"`
	Amount    float32 `json:"amount" db:"amount"`
}

type TransactionEntry struct {
	Transaction
	Entries []Entry `json:"entries"`
}

func (te *TransactionEntry) ToProto() *bookkeepingv1.TransactionEntry {
	entries := make([]*bookkeepingv1.Entry, len(te.Entries))

	for idx, e := range te.Entries {
		entries[idx] = &bookkeepingv1.Entry{
			AccountId: e.AccountID,
			Amount:    e.Amount,
		}
	}

	return &bookkeepingv1.TransactionEntry{
		Id:      te.TransactionID,
		Time:    timestamppb.New(te.Datetime),
		Entries: entries,
	}
}
