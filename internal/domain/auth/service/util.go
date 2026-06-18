package service

import "fbt/backend/internal/util"

func (s *service) Decrypt(encryptedValue string) (*string, error) {
	return util.Decrypt(encryptedValue, s.CFG.ENCRYPTION_KEY)
}

func (s *service) Encrypt(value string) (*string, error) {
	return util.Encrypt(value, s.CFG.ENCRYPTION_KEY)
}
