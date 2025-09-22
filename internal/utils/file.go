package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CalculateOptimalBlockSize(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	fileSize := int(fileInfo.Size())

	freeRAM, err := GetFreeRAM()
	if err != nil {
		return 0, err
	}

	if fileSize < int(freeRAM)/2 {
		return fileSize, nil
	} else {
		quarterRAM := freeRAM / 4
		if quarterRAM > 2*1024*1024*1024 {
			return 2 * 1024 * 1024 * 1024, nil
		}
		return int(quarterRAM), nil
	}
}

func CreatePathEncryptedFile(inputFilePath, outputFileName, outputDir string) (string, error) {
	var newPath string

	inputDir, inputFile := filepath.Split(inputFilePath)

	if outputFileName == "" {
		outputFileName = inputFile + ".crpt"
	} else {
		outputFileName += ".crpt"
	}

	if outputDir != "" {
		newPath = filepath.Join(outputDir, outputFileName)
	} else {
		newPath = filepath.Join(inputDir, outputFileName)
	}

	return newPath, nil
}

func CreatePathDecryptedFile(inputFilePath, outputFileName, outputDir string) (string, error) {
	var newPath string

	if !strings.HasSuffix(inputFilePath, ".crpt") {
		return "", fmt.Errorf("")
	}

	inputDir, inputFile := filepath.Split(inputFilePath)

	if outputFileName == "" {
		outputFileName = strings.TrimSuffix(inputFile, ".crpt")
	}

	if outputDir != "" {
		newPath = filepath.Join(outputDir, outputFileName)
	} else {
		newPath = filepath.Join(inputDir, outputFileName)
	}

	return newPath, nil
}

func ReadDecryptedFile(path string, blockSize int) (<-chan []byte, <-chan error, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, fmt.Errorf("file not found")
		}
		return nil, nil, err
	}

	outCh := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer f.Close()
		defer close(outCh)
		defer close(errCh)

		buf := make([]byte, blockSize)

		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			outCh <- buf[:n]
		}
	}()

	return outCh, errCh, nil
}

type Content struct {
	Nonce []byte
	Buf   []byte
}

func ReadEncryptedFile(f *os.File, nonceSize, blockSize, tagSize int) (<-chan Content, <-chan error, error) {
	outCh := make(chan Content)
	errCh := make(chan error)

	go func() {
		defer close(outCh)
		defer close(errCh)

		nonce := make([]byte, nonceSize)
		buf := make([]byte, blockSize+tagSize)
		for {
			n, err := f.Read(nonce)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			m, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				errCh <- err
			}

			outCh <- Content{
				Nonce: nonce[:n],
				Buf:   buf[:m],
			}
		}
	}()

	return outCh, errCh, nil
}
