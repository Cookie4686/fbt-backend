package mfa

import (
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
)

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
