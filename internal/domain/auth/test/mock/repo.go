package mock

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
)

type CredentialRepo struct {
	s *MockService
}

func NewCredentialRepo(s *MockService) *CredentialRepo {
	return &CredentialRepo{s}
}

func (m *CredentialRepo) Register(ctx context.Context, user *model.User, session *model.Session) error {
	m.s.Users = append(m.s.Users, *user)
	m.s.Sessions = append(m.s.Sessions, *session)
	return nil
}
