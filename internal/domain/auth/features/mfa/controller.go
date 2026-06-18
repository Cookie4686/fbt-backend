package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pquerna/otp/totp"
)

type con struct {
	service *service.AuthService
	repo    *repo
}

func NewController(service *service.AuthService, db *pgxpool.Pool) Controller {
	return Controller(con{service: service, repo: newRepo(db)})
}

func (s con) MFAStatus(ctx context.Context, auth *model.Auth) (*MFAStatusResponse, error) {
	userMfaList, err := s.repo.GetMFAList(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}
	return &MFAStatusResponse{StatusCode: http.StatusOK, Payload: userMfaList}, nil
}

func (s con) TOTPValidate(ctx context.Context, auth *model.Auth, body *TOTPValidatePayload) (*TOTPValidateResponse, error) {
	userTotp, err := s.repo.GetTOTP(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}

	secret, err := util.Decrypt(userTotp.Key, s.service.CFG.ENCRYPTION_KEY)
	if err != nil {
		return nil, err
	}

	isPassed := totp.Validate(body.Code, *secret)

	return &TOTPValidateResponse{
		StatusCode: http.StatusOK,
		Payload:    &TOTPValidateResponsePayload{IsPassed: isPassed}}, nil
}

func (s con) TOTPUpsertKey(ctx context.Context, auth *model.Auth, body *TOTPSetupPayload) (*TOTPSetupResponse, error) {
	encryptedKey, err := util.Encrypt(body.Key, s.service.CFG.ENCRYPTION_KEY)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpsertTOTP(ctx, *encryptedKey, auth.User.Id)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, auth.User.Id)
	if err != nil {
		return nil, err
	}
	err = s.service.InvalidateSession(ctx, &auth.Session)
	if err != nil {
		return nil, err
	}

	return &TOTPSetupResponse{StatusCode: http.StatusOK, Payload: session}, nil
}
