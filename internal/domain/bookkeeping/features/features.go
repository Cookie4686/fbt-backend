package features

import (
	"fbt/backend/internal/domain/bookkeeping/features/account"
	"fbt/backend/internal/domain/bookkeeping/features/transaction"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"

	"google.golang.org/grpc"
)

type Features struct {
	service service.Service
	d       *util.Dependency
}

func NewFeatures(d *util.Dependency) *Features {
	s := service.NewService(d)

	return &Features{service: s, d: d}
}

func (f *Features) RegisterAccount(s *grpc.Server) {
	account.RegisterService(f.service, account.NewRepo(f.d.DB), s)
}

func (f *Features) RegisterTransaction(s *grpc.Server) {
	transaction.RegisterService(f.service, transaction.NewRepo(f.d.DB), s)
}
