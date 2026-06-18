package mock

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/util"
	"slices"
	"time"
)

type MockService struct {
	Users    []model.User
	Sessions []model.Session
}

func NewMockService() *MockService {
	mock := &MockService{
		Users:    make([]model.User, 0),
		Sessions: make([]model.Session, 0),
	}
	return mock
}

func (m *MockService) CreateSession(ctx context.Context, userId string) (*model.Session, error) {
	session := &model.Session{
		Id:                util.GenerateBase64UUID(),
		UserId:            userId,
		ExpiresAt:         time.Now().Add(model.SessionExpiresIn),
		TwoFactorVerified: false,
	}

	// TODO: Mock PK
	m.Sessions = append(m.Sessions, *session)

	// TODO: Mock FK
	return session, nil
}

func (m *MockService) Validate(ctx context.Context, sessionId string) (*model.Auth, error) {
	session, err := m.findSession(sessionId)
	if err != nil {
		return nil, err
	}
	user, _ := m.findUser(session.UserId)

	auth := model.Auth{Session: *session, User: *user}

	if time.Now().After(auth.Session.ExpiresAt) {
		if err := m.InvalidateSession(ctx, &auth.Session); err != nil {
			return nil, errors.DBError
		}
		return nil, errors.SessionExpire
	}
	if time.Now().After(auth.Session.ExpiresAt.Add(-model.SessionExpiresIn / 2)) {
		auth.Session.ExpiresAt = time.Now().Add(model.SessionExpiresIn)
		if err := m.UpdateSessionExpiration(ctx, &auth.Session); err != nil {
			return nil, errors.DBError
		}
	}

	return &auth, nil
}

func (m *MockService) UpdateSessionExpiration(ctx context.Context, session *model.Session) error {
	s, err := m.findSession(session.Id)
	if err != nil {
		return err
	}
	s.ExpiresAt = session.ExpiresAt
	return nil
}

func (m *MockService) InvalidateSession(ctx context.Context, session *model.Session) error {
	session, err := m.findSession(session.Id)
	if err != nil {
		return err
	}
	m.Sessions = slices.DeleteFunc(m.Sessions, func(s model.Session) bool { return s.Id == session.Id })
	return nil
}

func (m *MockService) findSession(sessionID string) (*model.Session, error) {
	var session *model.Session
	for _, s := range m.Sessions {
		if s.Id == sessionID {
			session = &s
			break
		}
	}
	return session, errors.NotFound
}

func (m *MockService) findUser(userID string) (*model.User, error) {
	var user *model.User
	for _, u := range m.Users {
		if u.Id == userID {
			user = &u
			break
		}
	}
	return user, errors.NotFound
}

func (m *MockService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	idx := slices.IndexFunc(m.Users, func(u model.User) bool { return u.Username == username })
	if idx == -1 {
		return nil, errors.NotFound
	}
	return &m.Users[idx], nil
}

func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	idx := slices.IndexFunc(m.Users, func(u model.User) bool { return u.Email == email })
	if idx == -1 {
		return nil, errors.NotFound
	}
	return &m.Users[idx], nil
}

func (m *MockService) Decrypt(encryptedValue string) (*string, error) {
	return util.Decrypt(encryptedValue, "key")
}
func (m *MockService) Encrypt(value string) (*string, error) {
	return util.Encrypt(value, "key")
}
