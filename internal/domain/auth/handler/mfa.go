package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp/totp"
)

func (s *AuthHandler) GetUserMFAList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	params := mux.Vars(r)
	userID := params["id"]

	userMfaList, err := s.repo.GetMFAList(ctx, userID)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	util.SendJson(w, http.StatusOK, &userMfaList)
}

func (s *AuthHandler) TOTPValidate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		Code   string `json:"code"`
		UserID string `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	userTotp, err := s.repo.GetTOTP(ctx, body.UserID)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	secret, err := s.Decrypt(userTotp.Key)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	isPassed := totp.Validate(body.Code, *secret)

	response := struct {
		IsPassed bool `json:"is_passed"`
	}{IsPassed: isPassed}

	util.SendJson(w, http.StatusOK, &response)
}

func (s *AuthHandler) TOTPUpsertKey(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	body := struct {
		Key    string `json:"key"`
		UserID string `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}

	encryptedKey, err := s.Encrypt(body.Key)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}
	err = s.repo.UpsertTOTP(ctx, *encryptedKey, body.UserID)
	if err != nil {
		util.SendError(s.logger, w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *AuthHandler) Decrypt(encryptedValue string) (*string, error) {
	key, err := base64.StdEncoding.DecodeString(s.cfg.ENCRYPTION_KEY)
	if err != nil {
		return nil, err
	}
	value, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return nil, err
	}
	ciphertext, err := util.DecryptGCM(value, key)
	if err != nil {
		return nil, err
	}
	decryptedValue := string(ciphertext)
	return &decryptedValue, nil
}

func (s *AuthHandler) Encrypt(value string) (*string, error) {
	key, err := base64.StdEncoding.DecodeString(s.cfg.ENCRYPTION_KEY)
	if err != nil {
		return nil, err
	}
	ciphertext, err := util.EncryptGCM([]byte(value), key)
	if err != nil {
		return nil, err
	}
	encryptedValue := base64.StdEncoding.EncodeToString(ciphertext)
	return &encryptedValue, nil
}
