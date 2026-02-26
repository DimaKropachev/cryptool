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

func Decrypt(inPath, outPath string, password []byte) error {
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
			f.PB = progressbar.New(progressbar.PrefixDecrypt+": "+f.Info.Name(), f.Info.Size())

			err := decryptFile(f, password)
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Fprintf(os.Stdout, "Directory %s successfully decrypted", nodeInfo.Name())
	} else {
		outPath = file.CreateOutputFilePath(file.Decrpt, inPath, outPath)

		pb := progressbar.New(progressbar.PrefixDecrypt+": "+nodeInfo.Name(), nodeInfo.Size())

		file := &models.File{
			Name:    nodeInfo.Name(),
			Info:    nodeInfo,
			Path:    inPath,
			PB:      pb,
			OutPath: outPath,
		}

		err := decryptFile(file, password)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "File %s successfully decrypted", nodeInfo.Name())
	}

	return nil
}

func decryptFile(f *models.File, password []byte) error {
	inFile, err := os.OpenFile(f.Path, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(f.OutPath, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error accessing the output file: %w", err)
	}
	defer outFile.Close()

	f.PB.Start()
	header,n, err := crypto.DecryptHeader(inFile)
	if err != nil {
		f.PB.Finish()
		return fmt.Errorf("error reading header: %w", err)
	}
	f.PB.Add(n)

	alg, err := algorithms.CreateAlgorithmByID(int(header.AlgID), password, header.Salt)
	if err != nil {
		f.PB.Finish()
		return fmt.Errorf("error creating algorithm: %w", err)
	}

	content, errs, err := file.ReadEncryptedFile(inFile, alg.GetNonceSize(), int(header.BlockSize), alg.GetTagSize())
	if err != nil {
		f.PB.Finish()
		return fmt.Errorf("error reading input file: %w", err)
	}

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

			f.PB.Add(len(ciphertext.Nonce)+len(ciphertext.Buf))

			if _, err := outFile.Write(plaintext); err != nil {
				f.PB.Finish()
				return err
			}
		case err := <-errs:
			if err != nil {
				f.PB.Finish()
				return fmt.Errorf("error reading file: %w", err)
			}
		}
	}

	return nil
}
