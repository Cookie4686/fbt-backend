package service

import (
	"fbt/backend/internal/util"
)

type Service interface {
}

type service struct {
	*util.Dependency
}

func NewService(d *util.Dependency) Service {
	return Service(&service{
		Dependency: d,
	})
}
