package session

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
	service    *service.AuthService
	controller Controller
}

func NewFeature(d *dependency.Dependency, service *service.AuthService) *Feature {
	handler := &Handler{
		Dependency: d,
		service:    service,
		controller: NewController(service),
	}

	return &Feature{
		Logout:   http.HandlerFunc(handler.Logout),
		Validate: http.HandlerFunc(handler.Validate),
	}
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Validate(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Logout(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
