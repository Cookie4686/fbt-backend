package model

import (
	"time"
)

type Account struct {
	ID      int32  `json:"id" db:"account_id"`
	Name    string `json:"name" db:"name"`
	IsDebit bool   `json:"is_debit" db:"is_debit"`
	UserId  string `json:"user_id" db:"user_id"`
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
