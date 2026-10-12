package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func EncryptAES(plainText []byte, key []byte) (string, error) {
	// 1. create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 2. generate nonce (number used once)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// 3. Seal (encrypts the data)
	// cipherText := nonce + cipherText + tag
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	// 4. base64 encode
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(cipherBase64 string, key []byte) ([]byte, error) {
	// 1. base64 decode
	cipherText, err := base64.URLEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, err
	}

	// 2. create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 3. extract nonce
	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// 4. Open (decrypts the data)
	return aesGCM.Open(nil, nonce, cipherText, nil)
}
