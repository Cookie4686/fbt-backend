package util

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
)

const base32BitEntropy = 15
const base64BitEntropy = 18

func GenerateBase32UUID() string {
	bytes := make([]byte, base32BitEntropy)
	rand.Read(bytes)

	return base32.StdEncoding.EncodeToString(bytes)
}

func GenerateBase64UUID() string {
	bytes := make([]byte, base64BitEntropy)
	rand.Read(bytes)

	return base64.StdEncoding.EncodeToString(bytes)
}
