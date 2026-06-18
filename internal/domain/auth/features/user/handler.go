package user

import (
	"context"
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
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
