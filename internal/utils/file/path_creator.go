package file

import (
	"path/filepath"
	"strings"
)

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
