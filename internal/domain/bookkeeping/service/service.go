package service

import (
	"fbt/backend/internal/util"
)

type Service any
type service struct {
	*util.Dependency
}

func NewService(d *util.Dependency) Service {
	return Service(&service{
		Dependency: d,
	})
}
