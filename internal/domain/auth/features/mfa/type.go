package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"net/http"
)

type Feature struct {
	AUTH_MFAStatus     http.Handler
	AUTH_TOTPValidate  http.Handler
	AUTH_TOTPUpsertKey http.Handler
}

type Controller interface {
	MFAStatus(context.Context, *model.Auth) (*MFAStatusResponse, error)
	TOTPValidate(context.Context, *model.Auth, *TOTPValidatePayload) (*TOTPValidateResponse, error)
	TOTPUpsertKey(context.Context, *model.Auth, *TOTPSetupPayload) (*TOTPSetupResponse, error)
}

type Repo interface {
	GetMFAList(ctx context.Context, userID string) (*model.MfaList, error)
	GetTOTP(ctx context.Context, userID string) (*model.MfaTotp, error)
	UpsertTOTP(ctx context.Context, key string, userID string) error
}
