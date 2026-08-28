// Package service for bookkeeping-related service
package service

import (
	"fbt/backend/internal/util"
)

type (
	Service any
	service struct {
		*util.Dependency
	}
)

func NewService(d *util.Dependency) Service {
	return Service(&service{
		Dependency: d,
	})
}
