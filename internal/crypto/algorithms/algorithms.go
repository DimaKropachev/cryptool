package algorithms

const (
	MagicNum = "CRPT"

	AlgAES128GCM        = 1
	AlgAES192GCM        = 2
	AlgAES256GCM        = 3
	AlgCHACHA20POLY1305 = 4
)

type CipherAlgorithm interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext, nonce []byte) ([]byte, error)
	GetNonceSize() int
	GetTagSize() int
}
