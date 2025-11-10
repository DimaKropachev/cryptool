package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	DefaultSaltSize = 16
)

func GenerateSalt(size int) []byte {
	salt := make([]byte, size)
	rand.Read(salt)
	return salt
}

func GenerateNonce(size int) []byte {
	nonce := make([]byte, size)
	rand.Read(nonce)
	return nonce
}

func GenerateKeyFromPassword(password, salt []byte, keySize int) []byte {
	key := pbkdf2.Key(
		password,
		salt,
		100000,
		keySize,
		sha3.New256,
	)

	return key
}

func GenerateKey(keySize int) []byte {
	key := make([]byte, keySize)
	rand.Read(key)
	return key
}
