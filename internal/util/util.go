package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
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

func Decrypt(encryptedValue string, encryptionKey string) (*string, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, err
	}
	value, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return nil, err
	}
	ciphertext, err := DecryptGCM(value, key)
	if err != nil {
		return nil, err
	}
	decryptedValue := string(ciphertext)
	return &decryptedValue, nil
}

func Encrypt(value string, encryptionKey string) (*string, error) {
	key, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return nil, err
	}
	ciphertext, err := EncryptGCM([]byte(value), key)
	if err != nil {
		return nil, err
	}
	encryptedValue := base64.StdEncoding.EncodeToString(ciphertext)
	return &encryptedValue, nil
}

func EncryptGCM(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func DecryptGCM(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
