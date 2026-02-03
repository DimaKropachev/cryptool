package app

import (
	"fmt"
	"log"
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

	nodeInfo, err := os.Stat(inPath)
	if err != nil {
		return fmt.Errorf("error receiving information about an input data: %w", err)
	}

	if nodeInfo.IsDir() {
		outPath, err = file.CreateOutputDirPath(inPath, outPath)
		if err != nil {
			return err
		}
		if outPath[len(outPath)-1] != filepath.Separator {
			outPath += string(filepath.Separator)
		}

		files, err := file.ReadDirectory(inPath)
		if err != nil {
			return err
		}

		for _, f := range files {
			f.OutPath = file.CreateOutputFilePath(file.Encrpt, f.Path, outPath)

			f.PB = progressbar.New(progressbar.PrefixEncrypt+": "+f.Info.Name(), f.Info.Size())

			err := encryptFile(f, algorithm, password)
			if err != nil {
				log.Println(err)
			}
		}

	} else {
		// Создаем путь к выходнуму файлу
		outPath = file.CreateOutputFilePath(file.Encrpt, inPath, outPath)

		// Создаем прогресс бар
		pb := progressbar.New(progressbar.PrefixEncrypt+": "+nodeInfo.Name(), nodeInfo.Size())

		file := &models.File{
			Name:    nodeInfo.Name(),
			Info:    nodeInfo,
			Path:    inPath,
			OutPath: outPath,
			PB:      pb,
		}

		err := encryptFile(file, algorithm, password)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "File %s successfully encrypted", nodeInfo.Name())
	return nil
}

func encryptFile(f *models.File, algorithm string, password []byte) error {
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

	outFile, err := os.OpenFile(f.OutPath, os.O_CREATE, 0644)
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
