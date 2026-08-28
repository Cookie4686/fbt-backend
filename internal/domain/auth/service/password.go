package service

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	passwordSaltSize = 16
	idKeyTime        = 2
	idKeyMemory      = 19 * 1024
	idKeyThread      = 1
	idKeyLen         = 32
)

func (s *service) GenerateSalt() []byte {
	salt := make([]byte, passwordSaltSize)
	rand.Read(salt)

	return salt
}

func (s *service) HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, idKeyTime, idKeyMemory, idKeyThread, idKeyLen)
}
