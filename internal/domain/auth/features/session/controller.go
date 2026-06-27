package session

import (
	"context"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/interceptor"
	"net/http"

	"connectrpc.com/connect"
)

type Server struct {
	service service.Service
}

func NewServiceHandler(service service.Service, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewSessionServiceHandler(&Server{service}, opts...)
}

func (s *Server) Validate(ctx context.Context, in *authv1.SessionServiceValidateRequest) (*authv1.SessionServiceValidateResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	return &authv1.SessionServiceValidateResponse{
		Session: auth.Session.ToProto(),
		User:    auth.User.ToProto(),
	}, nil
}

func (s *Server) Logout(ctx context.Context, in *authv1.SessionServiceLogoutRequest) (*authv1.SessionServiceLogoutResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.service.InvalidateSession(ctx, &model.Session{Id: auth.Session.Id})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
