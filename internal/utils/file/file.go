package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DimaKropachev/cryptool/internal/utils"
)

func CalculateOptimalBlockSize(path string) (int, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	fileSize := int(fileInfo.Size())

	freeRAM, err := utils.GetFreeRAM()
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
	err := validateFilePath(inputFilePath)
	if err != nil {
		return "", err
	}
	err = validateFileName(outputFileName)
	if err != nil {
		return "", err
	}
	err = validateDirPath(outputDir)
	if err != nil {
		return "", err
	}

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
	err := validateFilePath(inputFilePath)
	if err != nil {
		return "", err
	}
	err = validateFileName(outputFileName)
	if err != nil {
		return "", err
	}
	err = validateDirPath(outputDir)
	if err != nil {
		return "", err
	}

	var newPath string

	if !strings.HasSuffix(inputFilePath, ".crpt") {
		return "", pathError(ActionValidate, inputFilePath, ErrInvalidFileExtension)
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

func validateFilePath(path string) error {
	if path == "" {
		return pathError(ActionValidate, path, ErrEmptyPath)
	}

	dir, file := filepath.Split(path)

	if dir != "" {
		err := validateDirPath(dir)
		if err != nil {
			pe, _ := err.(*PathError)
			return pathError(ActionValidate, path, pe.Err)

		}
	}

	if file != "" {
		err := validateFileName(file)
		if err != nil {
			return pathError(ActionValidate, path, err)
		}
	}

	return nil
}

func validateFileName(fileName string) error {
	if len(fileName) == 0 {
		return pathError(ActionValidate, fileName, ErrEmptyFileName)
	}

	forbiddenChar := `\/:*?"<>|`

	if strings.ContainsAny(fileName, forbiddenChar) {
		return pathError(ActionValidate, fileName, ErrForbiddenCharsFileName)
	}

	if strings.ReplaceAll(fileName, ".", "") == "" {
		return pathError(ActionValidate, fileName, ErrDotsFileName)
	}

	return nil
}

func validateDirPath(dirPath string) error {
	if len(dirPath) == 0 {
		return pathError(ActionValidate, dirPath, ErrEmptyDirPath)
	}

	forbiddenChar := `\/:*?"<>|`

	if strings.ContainsAny(dirPath, forbiddenChar) {
		return pathError(ActionValidate, dirPath, ErrForbiddenCharsDirPath)
	}

	var c1, c2, sep int
	for _, char := range dirPath {
		switch char {
		case '/':
			sep++
			c1++
		case filepath.Separator:
			sep++
			c2++
		default:
			sep = 0
		}

		if sep > 2 {
			return pathError(ActionValidate, dirPath, ErrSepInvalidSyntax)
		}
	}
	if c1 > 0 && c2 > 0 {
		return pathError(ActionValidate, dirPath, ErrSepInvalidSyntax)
	}

	folders := strings.Split(dirPath, string(filepath.Separator))
	for _, folder := range folders {
		if strings.ReplaceAll(folder, ".", "") == "" {
			return pathError(ActionValidate, dirPath, ErrFolderDotsDirPath)
		}
	}

	return nil
}
