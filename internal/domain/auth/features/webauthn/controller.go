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

	credential, err := s.repo.GetPasskey(ctx, in.RpId, in.CredentialId)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceGetUserPasskeyResponse{Credential: credential.ToProto()}, nil
}

func (s *Server) CreateUserPasskey(ctx context.Context, in *authv1.WebAuthnServiceCreateUserPasskeyRequest) (*authv1.WebAuthnServiceCreateUserPasskeyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	credential := &model.WebAuthnCredential{
		UserID:       auth.Session.UserId,
		RpID:         in.Credential.RpId,
		CredentialID: in.Credential.CredentialId,
		PublicKey:    in.Credential.PublicKey,
		Aaguid:       in.Credential.Aaguid,
		Counter:      in.Credential.Counter,
		DeviceType:   in.Credential.DeviceType,
		Transports:   in.Credential.Transports,
		BackupState:  in.Credential.BackedUp,
	}

	err = s.repo.CreatePasskey(ctx, credential)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceCreateUserPasskeyResponse{Credential: credential.ToProto()}, nil
}

func (s *Server) UpdatePasskeyCounter(ctx context.Context, in *authv1.WebAuthnServiceUpdatePasskeyCounterRequest) (*authv1.WebAuthnServiceUpdatePasskeyCounterResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	passkey, err := s.repo.UpdatePasskeyCounter(ctx, in.RpId, in.CredentialId, in.Counter)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, passkey.UserID, false)
	if err != nil {
		return nil, err
	}

	return &authv1.WebAuthnServiceUpdatePasskeyCounterResponse{Session: session.ToProto()}, nil
}
