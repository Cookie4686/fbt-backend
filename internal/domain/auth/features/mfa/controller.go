package mfa

import (
	"context"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/interceptor"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pquerna/otp/totp"
)

type Server struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewMFAServiceHandler(&Server{service, repo}, opts...)
}

func (s *Server) Status(ctx context.Context, in *authv1.MFAServiceStatusRequest) (*authv1.MFAServiceStatusResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	userMfaList, err := s.repo.GetMFAList(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}
	return &authv1.MFAServiceStatusResponse{TotpEnabled: userMfaList.Totp}, nil
}

func (s *Server) TOTPValidate(ctx context.Context, in *authv1.MFAServiceTOTPValidateRequest) (*authv1.MFAServiceTOTPValidateResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	userTotp, err := s.repo.GetTOTP(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	secret, err := s.service.Decrypt(userTotp.Key)
	if err != nil {
		return nil, err
	}

	if !totp.Validate(in.Code, *secret) {
		return nil, errors.BadRequest
	}

	session, err := s.service.CreateSession(ctx, auth.User.Id, true)
	if err != nil {
		return nil, err
	}

	return &authv1.MFAServiceTOTPValidateResponse{
		Session: session.ToProto(),
	}, nil
}

func (s *Server) TOTPUpsertKey(ctx context.Context, in *authv1.MFAServiceTOTPUpsertKeyRequest) (*authv1.MFAServiceTOTPUpsertKeyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := s.service.Encrypt(in.Key)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpsertTOTP(ctx, *encryptedKey, auth.User.Id)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, auth.User.Id, true)
	if err != nil {
		return nil, err
	}
	err = s.service.InvalidateSession(ctx, &auth.Session)
	if err != nil {
		return nil, err
	}

	return &authv1.MFAServiceTOTPUpsertKeyResponse{Session: session.ToProto()}, nil
}
