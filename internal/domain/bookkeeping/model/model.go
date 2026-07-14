package model

import (
	bookkeepingv1 "fbt/backend/gen/proto/go/bookkeeping/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Account struct {
	ID      int32  `db:"account_id" json:"id"`
	Name    string `db:"name"       json:"name"`
	IsDebit bool   `db:"is_debit"   json:"is_debit"`
	UserId  string `db:"user_id"    json:"user_id"`
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
	TagID  int64  `db:"tag_id"  json:"tag_id"`
	Name   string `db:"name"    json:"name"`
	UserID string `db:"user_id" json:"user_id"`
}

type AccountTag struct {
	Account

	Tags []Tag `json:"tags"`
}

type Transaction struct {
	TransactionID int64     `db:"transaction_id" json:"transaction_id"`
	Datetime      time.Time `db:"datetime"       json:"datetime"`
}

type Entry struct {
	// TransactionID int64   `json:"transaction_id" db:"transaction_id"`
	AccountID int32   `db:"account_id" json:"account_id"`
	Amount    float32 `db:"amount"     json:"amount"`
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
