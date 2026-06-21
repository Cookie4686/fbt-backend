package transaction

import (
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/util"
)

type EntryPayload struct {
	AccountID int32   `json:"account_id"`
	Amount    float32 `json:"amount"`
}

type GetAllPayload struct{}
type GetAllResponse = util.Response[[]model.TransactionEntry]

type CreatePayload struct {
	Datetime string `json:"datetime"`

	Entries []EntryPayload
}
type CreateResponse = util.Response[model.TransactionEntry]

type UpdatePayload struct {
	TransactionID int64  `json:"transaction_id"`
	Datetime      string `json:"datetime"`

	entries []EntryPayload
}
type UpdateResponse = util.Response[model.TransactionEntry]

type DeletePayload struct {
	TransactionID int64 `json:"transaction_id"`
}
type DeleteResponse = util.Response[struct{}]
