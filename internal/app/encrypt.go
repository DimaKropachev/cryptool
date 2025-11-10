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

func Encrypt(algorithm, inPath, outPath string, password []byte) error {
	inPath = filepath.Clean(inPath)
	outPath = filepath.Clean(outPath)

	nodeInfo, err := os.Stat(inPath)
	if err != nil {
		return fmt.Errorf("error receiving information about an input data: %w", err)
	}

	if nodeInfo.IsDir() {
		_, err := file.ReadDirectory(inPath)
		if err != nil {
			return err
		}
		// TODO: что дклать если пользователь указал директорию
	} else {
		pb := progressbar.New(progressbar.PrefixEncrypt+": "+nodeInfo.Name(), nodeInfo.Size())

		file := &models.File{
			Name: nodeInfo.Name(),
			Info: nodeInfo,
			Path: inPath,
			PB:   pb,
		}

		err := encryptFile(file, outPath, algorithm, password)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "File %s successfully encrypted", nodeInfo.Name())
	return nil
}

func encryptFile(f *models.File, outPath, algorithm string, password []byte) error {
	salt := crypto.GenerateSalt(crypto.DefaultSaltSize)

	alg, algID, err := algorithms.CreateAlgorithmByName(algorithm, password, salt)
	if err != nil {
		return err
	}

	blockSize, err := CalculateOptimalBlockSize(int(f.Info.Size()))
	if err != nil {
		return err
	}

	header := crypto.NewHeader(algID, blockSize, len(salt), alg.GetNonceSize(), salt)
	encHeader, err := crypto.EncryptHeader(header)
	if err != nil {
		return err
	}

	if outPath == "." {
		outPath = f.Name + ".crpt"
	}

	outFile, err := os.OpenFile(outPath, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error accessing the output file: %w", err)
	}
	defer outFile.Close()

	if _, err := outFile.Write(encHeader); err != nil {
		return fmt.Errorf("error writing the header: %w", err)
	}

	content, errs, err := file.ReadDecryptedFile(f.Path, blockSize)
	if err != nil {
		return err
	}

	f.PB.Start()
READ:
	for {
		select {
		case plaintext, ok := <-content:
			if !ok {
				f.PB.Finish()
				break READ
			}

			ciphertext, err := alg.Encrypt(plaintext)
			if err != nil {
				f.PB.Finish()
				return err
			}

			if _, err := outFile.Write(ciphertext); err != nil {
				f.PB.Finish()
				return err
			}

			f.PB.Add(blockSize)
		case err := <-errs:
			if err != nil {
				f.PB.Finish()
				return fmt.Errorf("error reading file: %w", err)
			}
		}
	}
	return nil
}

func encryptDirectory(fs []*models.File, algorithm string, password []byte) error {

	return nil
}
