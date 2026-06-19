package mock

import "fbt/backend/internal/domain/auth/features/credentials"

var (
	ValidUser = credentials.RegisterPayload{
		Username: "test",
		Password: "12345678",
		Email:    "test@email.com",
	}
)
