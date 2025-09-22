package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher" 

	"github.com/DimaKropachev/cryptool/internal/utils"
)

type AESGCM struct {
	gcm       cipher.AEAD
	NonceSize int
	TagSize   int
}

func NewAESGCM(password []byte, keySize int, salt []byte) (*AESGCM, error) {
	key := utils.GenerateKeyFromPassword(password, salt, keySize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	aes := &AESGCM{
		gcm:       gcm,
		NonceSize: gcm.NonceSize(),
		TagSize:   gcm.Overhead(),
	}

	return aes, nil
}

func (aes *AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	// generate nonce
	nonce := utils.GenerateNonce(aes.NonceSize)

	if _, err := result.Write(nonce); err != nil {
		return nil, err
	}

	// encipher plaintext
	ciphertext := aes.gcm.Seal(nil, nonce, plaintext, nil)

	if _, err := result.Write(ciphertext); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (aes *AESGCM) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	plaintext, err := aes.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (aes *AESGCM) GetNonceSize() int {
	return aes.NonceSize
}

func (aes *AESGCM) GetTagSize() int {
	return aes.TagSize
}