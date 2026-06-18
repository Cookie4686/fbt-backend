package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

type Feature struct {
	MFAStatus     http.HandlerFunc
	TOTPValidate  http.HandlerFunc
	TOTPUpsertKey http.HandlerFunc
}

type Controller interface {
	MFAStatus(context.Context, *model.Auth) (*MFAStatusResponse, error)
	TOTPValidate(context.Context, *model.Auth, *TOTPValidatePayload) (*TOTPValidateResponse, error)
	TOTPUpsertKey(context.Context, *model.Auth, *TOTPSetupPayload) (*TOTPSetupResponse, error)
}

type TOTPValidatePayload struct {
	Code string `json:"code"`
}
type TOTPValidateResponsePayload struct {
	IsPassed bool `json:"is_passed"`
}
type TOTPValidateResponse = util.Response[TOTPValidateResponsePayload]

type TOTPSetupPayload struct {
	Key string `json:"key"`
}
type TOTPSetupResponse = util.Response[model.Session]

type MFAStatusPayload struct{}
type MFAStatusResponse = util.Response[model.MfaList]
