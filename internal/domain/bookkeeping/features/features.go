package features

import (
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"
)

type features struct {
	Account     *account.Feature
	Transaction *transaction.Feature
}

func NewFeatures(d *util.Dependency) *features {
	s := service.NewService(d)

	return &features{
		Account:     account.NewFeature(d, s),
		Transaction: transaction.NewFeature(d, s),
	}
}
