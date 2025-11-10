package algorithms

import (
	"github.com/DimaKropachev/cryptool/pkg/crypto"
	"github.com/DimaKropachev/cryptool/pkg/crypto/algorithms/aes"
	"github.com/DimaKropachev/cryptool/pkg/crypto/algorithms/chacha20"
)

const (
	IDAES128GCM        = 1
	IDAES192GCM        = 2
	IDAES256GCM        = 3
	IDCHACHA20POLY1305 = 4

	AlgAES128GCM        = "aes128-gcm"
	AlgAES192GCM        = "aes192-gcm"
	AlgAES256GCM        = "aes256-gcm"
	AlgCHACHA20POLY1305 = "chacha20-poly1305"
)

type CipherAlgorithm interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext, nonce []byte) ([]byte, error)
	GetNonceSize() int
	GetTagSize() int
}

func CreateAlgorithmByName(algorithm string, password, salt []byte) (CipherAlgorithm, int, error) {
	var (
		alg CipherAlgorithm
		id  int
		err error
	)

	if len(password) != 0 {
		switch algorithm {
		case AlgAES128GCM:
			id = IDAES128GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 16), salt)
		case AlgAES192GCM:
			id = IDAES192GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 24), salt)
		case AlgAES256GCM:
			id = IDAES256GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 32), salt)
		case AlgCHACHA20POLY1305:
			id = IDCHACHA20POLY1305
			alg, err = chacha20.NewChaCha20Poly1305(crypto.GenerateKeyFromPassword(password, salt, 32), salt)
		}
		if err != nil {
			return nil, id, err
		}
	} else {
		// TODO: что делать если пользователь не ввел пароль
		switch algorithm {
		case AlgAES128GCM:
			id = IDAES128GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKey(16), salt)
		case AlgAES192GCM:
			id = IDAES192GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKey(24), salt)
		case AlgAES256GCM:
			id = IDAES256GCM
			alg, err = aes.NewAESGCM(crypto.GenerateKey(32), salt)
		case AlgCHACHA20POLY1305:
			id = IDCHACHA20POLY1305
			alg, err = chacha20.NewChaCha20Poly1305(crypto.GenerateKey(32), salt)
		}
		if err != nil {
			return nil, id, err
		}
	}

	return alg, id, nil
}

func CreateAlgorithmByID(algorithm int, password, salt []byte) (CipherAlgorithm, error) {
	var (
		alg CipherAlgorithm
		err error
	)

	if len(password) != 0 {
		switch algorithm {
		case IDAES128GCM:
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 16), salt)
		case IDAES192GCM:
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 24), salt)
		case IDAES256GCM:
			alg, err = aes.NewAESGCM(crypto.GenerateKeyFromPassword(password, salt, 32), salt)
		case IDCHACHA20POLY1305:
			alg, err = chacha20.NewChaCha20Poly1305(crypto.GenerateKeyFromPassword(password, salt, 32), salt)
		}
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: что делать если пользователь не ввел пароль
	}

	return alg, nil
}
