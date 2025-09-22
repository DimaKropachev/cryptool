package chacha20

import ( 
	"bytes"
	"crypto/cipher"

	"github.com/DimaKropachev/cryptool/internal/utils"
	"golang.org/x/crypto/chacha20poly1305"
)

type ChaCha20Poly1305 struct {
	aead      cipher.AEAD
	NonceSize int
	TagSize   int
}

func NewChaCha20Poly1305(password []byte, keySize int, salt []byte) (*ChaCha20Poly1305, error) {
	key := utils.GenerateKeyFromPassword(password, salt, keySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	return &ChaCha20Poly1305{
		aead:      aead,
		NonceSize: aead.NonceSize(),
		TagSize:   aead.Overhead(),
	}, nil
}

func (chacha20 *ChaCha20Poly1305) Encrypt(plaintext []byte) ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	nonce := utils.GenerateNonce(chacha20.NonceSize)

	if _, err := result.Write(nonce); err != nil {
		return nil, err
	}

	ciphertext := chacha20.aead.Seal(nil, nonce, plaintext, nil)

	if _, err := result.Write(ciphertext); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (chacha20 *ChaCha20Poly1305) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
	plaintext, err := chacha20.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (chacha20 *ChaCha20Poly1305) GetNonceSize() int {
	return chacha20.NonceSize
}

func (chacha20 *ChaCha20Poly1305) GetTagSize() int {
	return chacha20.TagSize
}