package app

import (
	"fmt"
	"os"

	"github.com/DimaKropachev/cryptool/internal/crypto/algorithms"
	"github.com/DimaKropachev/cryptool/internal/crypto/algorithms/aes"
	"github.com/DimaKropachev/cryptool/internal/crypto/algorithms/chacha20"
	"github.com/DimaKropachev/cryptool/internal/utils/file"
	"github.com/DimaKropachev/cryptool/internal/utils"
)

func EncryptAndSave(algorithm, inputPath, outputDirPath, outputFileName string, password []byte, saltSize int) error {
	var (
		algID int
		alg   algorithms.CipherAlgorithm
		err   error
		salt  = utils.GenerateSalt(saltSize)
	)

	switch algorithm {
	case "aes128-gcm":
		algID = algorithms.AlgAES128GCM
		alg, err = aes.NewAESGCM(password, 16, salt)
	case "aes192-gcm":
		algID = algorithms.AlgAES192GCM
		alg, err = aes.NewAESGCM(password, 24, salt)
	case "aes256-gcm":
		algID = algorithms.AlgAES256GCM
		alg, err = aes.NewAESGCM(password, 32, salt)
	case "chacha20-poly1305":
		algID = algorithms.AlgCHACHA20POLY1305
		alg, err = chacha20.NewChaCha20Poly1305(password, 32, salt)
	}
	if err != nil {
		return err
	}

	outPath, err := file.CreatePathEncryptedFile(inputPath, outputFileName, outputDirPath)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(outPath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	blockSize, err := file.CalculateOptimalBlockSize(inputPath)
	if err != nil {
		return err
	}

	header := utils.NewHeader(algID, blockSize, saltSize, alg.GetNonceSize(), salt)
	encHeader, err := utils.EncryptHeader(header)
	if err != nil {
		return fmt.Errorf("error encrypting header: %w", err)
	}

	// writting header
	if _, err := out.Write(encHeader); err != nil {
		return fmt.Errorf("error writting header from file")
	}

	content, errs, err := file.ReadDecryptedFile(inputPath, blockSize)
	if err != nil {
		return err
	}

READ:
	for {
		select {
		case plaintext, ok := <-content:
			if !ok {
				break READ
			}

			ciphertext, err := alg.Encrypt(plaintext)
			if err != nil {
				return err
			}

			if _, err = out.Write(ciphertext); err != nil {
				return err
			}

		case err, ok := <-errs:
			if !ok {
				break READ
			}
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	return nil
}

func DecryptAndSave(inputPath, outputDirPath, outputFileName string, password []byte) error {
	outPath, err := file.CreatePathDecryptedFile(inputPath, outputFileName, outputDirPath)
	if err != nil {
		return err
	}

	in, err := os.OpenFile(inputPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer in.Close()

	header, err := utils.DecryptHeader(in)
	if err != nil {
		return err
	}

	var alg algorithms.CipherAlgorithm
	switch header.AlgID {
	case algorithms.AlgAES128GCM:
		alg, err = aes.NewAESGCM(password, 16, header.Salt)
	case algorithms.AlgAES192GCM:
		alg, err = aes.NewAESGCM(password, 24, header.Salt)
	case algorithms.AlgAES256GCM:
		alg, err = aes.NewAESGCM(password, 32, header.Salt)
	case algorithms.AlgCHACHA20POLY1305:
		alg, err = chacha20.NewChaCha20Poly1305(password, 32, header.Salt)
	}
	if err != nil {
		return err
	}

	out, err := os.OpenFile(outPath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	content, errs, err := file.ReadEncryptedFile(in, alg.GetNonceSize(), int(header.BlockSize), alg.GetTagSize())
	if err != nil {
		return err
	}

READ:
	for {
		select {
		case ciphertext, ok := <-content:
			if !ok {
				break READ
			}

			plaintext, err := alg.Decrypt(ciphertext.Buf, ciphertext.Nonce)
			if err != nil {
				return err
			}

			if _, err := out.Write(plaintext); err != nil {
				return err
			}

		case err, ok := <-errs:
			if !ok {
				break READ
			}
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	return nil
}
