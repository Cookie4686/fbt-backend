package bookkeeping

import (
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"
	"net/http"

	"connectrpc.com/connect"
)

func RegisterService(mux *http.ServeMux, d *util.Dependency, opts ...connect.HandlerOption) *http.ServeMux {
	s := service.NewService(d)

	mux.Handle(account.NewServiceHandler(s, account.NewRepo(d.DB), opts...))
	mux.Handle(transaction.NewServiceHandler(s, transaction.NewRepo(d.DB), opts...))

	return mux
}
