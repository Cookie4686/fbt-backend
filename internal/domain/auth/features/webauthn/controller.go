package webauthn

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
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewWebAuthnServiceHandler(&Server{service, repo}, opts...)
}

func (s *Server) GetUserPasskey(ctx context.Context, in *authv1.WebAuthnServiceGetUserPasskeyRequest) (*authv1.WebAuthnServiceGetUserPasskeyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	passkey, err := s.repo.GetPasskey(ctx, in.PasskeyId)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceGetUserPasskeyResponse{Passkey: passkey.ToProto()}, nil
}

func (s *Server) CreateUserPasskey(ctx context.Context, in *authv1.WebAuthnServiceCreateUserPasskeyRequest) (*authv1.WebAuthnServiceCreateUserPasskeyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	passkey := &model.Passkey{
		PasskeyID:      in.Passkey.PasskeyId,
		PublicKey:      in.Passkey.PublicKey,
		UserID:         auth.Session.UserId,
		WebauthnUserID: in.Passkey.WebauthnId,
		Counter:        in.Passkey.Counter,
		DeviceType:     in.Passkey.DeviceType,
		BackedUp:       in.Passkey.BackedUp,
		Transports:     in.Passkey.Transports,
	}

	err = s.repo.CreatePasskey(ctx, passkey)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceCreateUserPasskeyResponse{Passkey: passkey.ToProto()}, nil
}

func (s *Server) UpdatePasskeyCounter(ctx context.Context, in *authv1.WebAuthnServiceUpdatePasskeyCounterRequest) (*authv1.WebAuthnServiceUpdatePasskeyCounterResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	passkey, err := s.repo.UpdatePasskeyCounter(ctx, in.PasskeyId, in.Counter)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, passkey.UserID, false)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceUpdatePasskeyCounterResponse{Session: session.ToProto()}, nil
}
