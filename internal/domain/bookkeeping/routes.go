package bookkeeping

import (
	"fbt/backend/internal/domain/bookkeeping/features"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(d *util.Dependency, r *mux.Router) {
	f := features.NewFeatures(d)

	r.Handle("/accounts", f.Account.GetAll).Methods(http.MethodGet)
	r.Handle("/accounts", f.Account.Create).Methods(http.MethodPost)
	r.Handle("/accounts", f.Account.Update).Methods(http.MethodPut)
	r.Handle("/accounts", f.Account.Delete).Methods(http.MethodDelete)

	r.Handle("/transactions", f.Transaction.GetAll).Methods(http.MethodGet)
	r.Handle("/transactions", f.Transaction.Create).Methods(http.MethodPost)
}
