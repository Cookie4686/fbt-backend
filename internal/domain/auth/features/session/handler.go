package session

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
	handler := &handler{
		Dependency: d,
		controller: NewController(service),
	}

	return &Feature{
		AUTH_Logout:   m.Auth(http.HandlerFunc(handler.Logout)),
		AUTH_Validate: m.Auth(http.HandlerFunc(handler.Validate)),
	}
}

func (h *handler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Validate(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.Logout(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
