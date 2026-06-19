package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"
)

type Handler struct {
	*util.Dependency
	controller Controller
}

func NewFeature(d *util.Dependency, service service.Service) *Feature {
	handler := &Handler{
		Dependency: d,
		controller: NewController(service, d.DB),
	}

	return &Feature{
		MFAStatus:     http.HandlerFunc(handler.MFAStatus),
		TOTPValidate:  http.HandlerFunc(handler.TOTPValidate),
		TOTPUpsertKey: http.HandlerFunc(handler.TOTPUpsertKey),
	}
}

func (h *Handler) MFAStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.MFAStatus(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *Handler) TOTPValidate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if payload, err := util.ExtractPayload[TOTPValidatePayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.TOTPValidate(ctx, auth, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *Handler) TOTPUpsertKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if payload, err := util.ExtractPayload[TOTPValidatePayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.TOTPValidate(ctx, auth, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
