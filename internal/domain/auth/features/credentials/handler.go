package credentials

import (
	"context"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"
)

type handler struct {
	*util.Dependency
	controller Controller
}

func NewFeature(d *util.Dependency, service service.Service) *Feature {
	handler := &handler{
		Dependency: d,
		controller: NewController(service, NewRepo(d.DB)),
	}

	return &Feature{
		Register: http.HandlerFunc(handler.Register),
		Login:    http.HandlerFunc(handler.Login),
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
