package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	AlgAES256GCM = 1
	AlgAES192GCM = 2
	AlgAES128GCM = 3
)

type Header struct {
	Algorithm uint16
	Salt      [16]byte
	Nonce     [12]byte
	Iteration uint32
}

// func main() {
// 	message := "secret data"
// 	password := []byte("qwerty123")
// 	alg := "aes-gcm-128"
// 	iterations := 100000

// 	ciphertext, err := aesEncrypt(message, string(password), alg, iterations)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(ciphertext)
// 	fmt.Println()

// 	res, err := aesDecrypt(ciphertext, password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println(res)

// 	// key := make([]byte, 32)
// 	// if _, err := io.ReadFull(rand.Reader, key); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println("key =", fmt.Sprintf("%x", key), "len(key)", len(key))

// 	// plaintext := []byte("secret data")

// 	// block, err := aes.NewCipher(key)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println(block.BlockSize())

// 	// gcm, err := cipher.NewGCM(block)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println("Nonce size =", gcm.NonceSize())

// 	// nonce := make([]byte, gcm.NonceSize())
// 	// if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
// 	// log.Printf("Encrypted: %x\n", ciphertext)
// }

func aesEncrypt(plaintext, password, alg string, iteratoins int) ([]byte, error) {
	var (
		key    []byte
		header Header
	)

	salt := generateSalt()
	header.Salt = [16]byte(salt)

	switch alg {
	case "aes-gcm-128":
		key = createKeyFromPassword([]byte(password), salt, 16, iteratoins)
		header.Algorithm = 3
	case "aes-gcm-192":
		generateSalt()
		key = createKeyFromPassword([]byte(password), salt, 24, iteratoins)
		header.Algorithm = 2
	case "aes-gcm-256":
		salt := generateSalt()
		key = createKeyFromPassword([]byte(password), salt, 32, iteratoins)
		header.Algorithm = 1
	}

	fmt.Println("key =", fmt.Sprintf("%x", key), "len(key) =", len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := generateNonce(gcm.NonceSize())
	header.Nonce = [12]byte(nonce)
	header.Iteration = uint32(iteratoins)

	binHeader := bytes.Buffer{}
	if err := binary.Write(&binHeader, binary.NativeEndian, header); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), binHeader.Bytes())

	return append(binHeader.Bytes(), ciphertext...), nil
}

func createKeyFromPassword(password, salt []byte, n, iterations int) []byte {
	key := pbkdf2.Key(
		password,
		salt,
		iterations,
		n,
		sha3.New256,
	)
	return key
}

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func generateNonce(n int) []byte {
	nonce := make([]byte, n)
	rand.Read(nonce)
	return nonce
}

func aesDecrypt(ciphertext, password []byte) (string, error) {
	headerSize := binary.Size(Header{})

	if len([]byte(ciphertext)) <= headerSize {
		return "", fmt.Errorf("cipher text is too small")
	}
	headerBytes := []byte(ciphertext)[:headerSize]

	var header Header
	buffer := bytes.NewReader(headerBytes)
	err := binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		return "", err
	}

	fmt.Printf("Algorithm ID: %d\n", header.Algorithm)
	fmt.Printf("Salt: %x\n", header.Salt)
	fmt.Printf("Nonce: %x\n", header.Nonce)

	ciphertext = ciphertext[headerSize:]

	var n int
	switch header.Algorithm {
	case 1:
		n = 32
	case 2:
		n = 24
	case 3:
		n = 16
	}
	key := createKeyFromPassword(password, header.Salt[:], n, int(header.Iteration))
	fmt.Printf("key = %x len(key) = %d\n", key, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, header.Nonce[:], []byte(ciphertext), headerBytes)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
