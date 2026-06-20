package user

import (
	"context"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

type Feature struct {
	GetByUsername http.HandlerFunc
}

type Handler struct {
	*util.Dependency
	controller Controller
}

func NewFeature(d *util.Dependency, service service.Service) *Feature {
	handler := &Handler{
		Dependency: d,
		controller: NewController(service),
	}

	return &Feature{
		GetByUsername: http.HandlerFunc(handler.GetByUsername),
	}
}

func (h *Handler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	params := mux.Vars(r)
	username := params["username"]

	if response, err := h.controller.GetByUsername(ctx, username); err != nil {
		util.SendError(h.Logger, w, r, err)
	} else if err := util.SendJson(w, response); err != nil {
		util.SendError(h.Logger, w, r, err)
	}
}
