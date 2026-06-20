package mfa

import (
	"context"
	"fbt/backend/internal/domain/auth/middleware"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"
)

type handler struct {
	*util.Dependency
	controller Controller
}

func NewFeature(d *util.Dependency, service service.Service, m middleware.Middleware) *Feature {
	h := &handler{
		Dependency: d,
		controller: NewController(service, d.DB),
	}

	return &Feature{
		AUTH_MFAStatus:     m.Auth(http.HandlerFunc(h.MFAStatus)),
		AUTH_TOTPValidate:  m.Auth(http.HandlerFunc(h.TOTPValidate)),
		AUTH_TOTPUpsertKey: m.Auth(http.HandlerFunc(h.TOTPUpsertKey)),
	}
}

func (h *handler) MFAStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.MFAStatus(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) TOTPValidate(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) TOTPUpsertKey(w http.ResponseWriter, r *http.Request) {
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
