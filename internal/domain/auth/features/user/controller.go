package user

import (
	"context"
	"fbt/backend/internal/domain/auth/service"
	"net/http"
)

type Controller interface {
	GetByUsername(ctx context.Context, username string) (*GetByUsernameResponse, error)
}

type con struct {
	service service.Service
}

func NewController(service service.Service) Controller {
	return Controller(con{service: service})
}

func (s con) GetByUsername(ctx context.Context, username string) (*GetByUsernameResponse, error) {
	user, err := s.service.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &GetByUsernameResponse{StatusCode: http.StatusOK, Payload: user}, nil
}
