package oauth

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

func NewFeature(d *util.Dependency, s service.Service, m middleware.Middleware) *Feature {
	h := &handler{
		Dependency: d,
		controller: NewController(s, d.DB),
	}

	return &Feature{
		Register:    http.HandlerFunc(h.Register),
		Login:       http.HandlerFunc(h.Login),
		AUTH_Status: m.Auth(http.HandlerFunc(h.Status)),
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if payload, err := util.ExtractPayload[RegisterPayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.Register(ctx, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if payload, err := util.ExtractPayload[LoginPayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.Login(ctx, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Status(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Status(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
