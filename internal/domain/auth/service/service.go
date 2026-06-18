package service

import (
	"fbt/backend/internal/dependency"
)

type AuthService struct {
	*dependency.Dependency
}

func NewAuthService(d *dependency.Dependency) *AuthService {
	return &AuthService{d}
}
