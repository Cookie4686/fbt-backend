package handler

import (
	"context"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	params := mux.Vars(r)
	username := params["username"]

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	util.SendJson(w, http.StatusOK, &user)
}
