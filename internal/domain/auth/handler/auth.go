package handler

import (
	"context"
	"encoding/json"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/util"
	"net/http"
)

func (s *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		SessionId string `json:"token"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	auth, err := s.repo.Validate(ctx, body.SessionId)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	util.SendJson(w, http.StatusOK, &auth)
}

func (s *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		SessionId string `json:"token"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	err := s.repo.InvalidateSession(ctx, &model.Session{Id: body.SessionId})
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
