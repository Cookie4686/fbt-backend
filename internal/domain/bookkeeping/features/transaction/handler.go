package transaction

import (
	"context"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/bookkeeping/service"
	"fbt/backend/internal/util"
	"net/http"
)

type handler struct {
	*util.Dependency
	controller Controller
}

func NewFeature(d *util.Dependency, service service.Service) *Feature {
	h := &handler{
		Dependency: d,
		controller: NewController(service, NewRepo(d.DB)),
	}

	return &Feature{
		GetAll: http.HandlerFunc(h.GetAll),
		Create: http.HandlerFunc(h.Create),
		Update: http.HandlerFunc(h.Update),
		Delete: http.HandlerFunc(h.Delete),
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if response, err := h.controller.GetAll(ctx, auth); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if payload, err := util.ExtractPayload[CreatePayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.Create(ctx, auth, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if payload, err := util.ExtractPayload[UpdatePayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.Update(ctx, auth, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	auth := ctx.Value("auth").(*model.Auth)

	if payload, err := util.ExtractPayload[DeletePayload](r); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if response, err := h.controller.Delete(ctx, auth, payload); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
