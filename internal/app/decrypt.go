package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DimaKropachev/cryptool/pkg/crypto"
	"github.com/DimaKropachev/cryptool/pkg/crypto/algorithms"
	"github.com/DimaKropachev/cryptool/pkg/file"
	"github.com/DimaKropachev/cryptool/pkg/models"
	"github.com/DimaKropachev/cryptool/pkg/progressbar"
)

func Decrypt(inPath, outPath string, password []byte) error {
	inPath = filepath.Clean(inPath)
	outPath = filepath.Clean(outPath)

	nodeInfo, err := os.Stat(inPath)
	if err != nil {
		return fmt.Errorf("error receiving information about an input data: %w", err)
	}

	if nodeInfo.IsDir() {

	} else {
		pb := progressbar.New(progressbar.PrefixDecrypt+": "+nodeInfo.Name(), nodeInfo.Size())

		file := &models.File{
			Name: nodeInfo.Name(),
			Info: nodeInfo,
			Path: inPath,
			PB:   pb,
		}

		err := decryptFile(file, outPath, password)
		if err != nil {
			return err
		}
	}

	in, err := os.OpenFile(inPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer in.Close()

	header, err := crypto.DecryptHeader(in)
	if err != nil {
		// error decrypting header
		return err
	}

	alg, err := algorithms.CreateAlgorithmByID(int(header.AlgID), password, header.Salt)
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

func decryptFile(f *models.File, outPath string, password []byte) error {
	inFile, err := os.OpenFile(f.Path, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outPath, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error accessing the output file: %w", err)
	}
	defer outFile.Close()

	header, err := crypto.DecryptHeader(inFile)
	if err != nil {
		return fmt.Errorf("error reading header: %w", err)
	}

	alg, err := algorithms.CreateAlgorithmByID(int(header.AlgID), password, header.Salt)
	if err != nil {
		return fmt.Errorf("error creating algorithm: %w", err)
	}

	content, errs, err := file.ReadEncryptedFile(inFile, alg.GetNonceSize(), int(header.BlockSize), alg.GetTagSize())
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	f.PB.Start()
READ:
	for {
		select {
		case ciphertext, ok := <-content:
			if !ok {
				f.PB.Finish()
				break READ
			}

			plaintext, err := alg.Decrypt(ciphertext.Buf, ciphertext.Nonce)
			if err != nil {
				f.PB.Finish()
				return err
			}

			if _, err := outFile.Write(plaintext); err != nil {
				f.PB.Finish()
				return err
			}

			f.PB.Add(int(header.BlockSize))
		case err := <-errs:
			if err != nil {
				f.PB.Finish()
				return fmt.Errorf("error reading file: %w", err)
			}
		}
	}

	return nil
}
