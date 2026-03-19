package util

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
)

func GenerateBase32UUID() string {
	bytes := make([]byte, 15)
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func GenerateBase64UUID() string {
	bytes := make([]byte, 18)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}
