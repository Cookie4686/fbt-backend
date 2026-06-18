package session

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"net/http"
)

type con struct {
	service service.Service
}

func NewController(service service.Service) Controller {
	return Controller(con{service: service})
}

func (s con) Validate(ctx context.Context, auth *model.Auth) (*ValidateResponse, error) {
	auth, err := s.service.Validate(ctx, auth.Session.Id)
	if err != nil {
		return nil, err
	}
	return &ValidateResponse{StatusCode: http.StatusOK, Payload: auth}, nil
}

func (s con) Logout(ctx context.Context, auth *model.Auth) (*LogoutResponse, error) {
	err := s.service.InvalidateSession(ctx, &model.Session{Id: auth.Session.Id})
	if err != nil {
		return nil, err
	}

	return &LogoutResponse{StatusCode: http.StatusOK}, nil
}
