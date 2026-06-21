package account

import (
	"fbt/backend/internal/domain/bookkeeping/model"
	"fbt/backend/internal/util"
)

type GetAllPayload struct{}
type GetAllResponse = util.Response[*[]model.Account]

type CreatePayload struct {
	Name    string `json:"name"`
	IsDebit bool   `json:"is_debit"`
}
type CreateResponse = util.Response[model.Account]

type UpdatePayload struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	IsDebit bool   `json:"is_debit"`
}
type UpdateResponse = util.Response[model.Account]

type DeletePayload struct {
	ID int32 `json:"id"`
}
type DeleteResponse = util.Response[struct{}]
