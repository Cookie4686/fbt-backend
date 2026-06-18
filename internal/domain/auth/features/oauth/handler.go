package oauth

import (
	"context"
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"
)

type Handler struct {
	*dependency.Dependency
	controller Controller
}

func NewFeature(d *dependency.Dependency, service service.Service) *Feature {
	handler := &Handler{
		Dependency: d,
		controller: NewController(service, d.DB),
	}

	return &Feature{
		Register: http.HandlerFunc(handler.Register),
		Login:    http.HandlerFunc(handler.Login),
		Status:   http.HandlerFunc(handler.Status),
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Status(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
