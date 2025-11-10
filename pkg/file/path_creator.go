package file

import (
	"path/filepath"
	"strings"
)

func CreateOutPath(inPath, outPath string) string {
	var newPath string

	outPathDir, outPathFile := filepath.Split(outPath)
	if outPathFile == "" {
		_, inputFile := filepath.Split(inPath)
		newPath = outPathDir + inputFile + ".crpt"
	} else {
		newPath = outPath
	}

	return newPath
}

func CreatePathDecryptedFile(inputFilePath, outputFileName, outputDir string) (string, error) {
	err := ValidateFilePath(inputFilePath)
	if err != nil {
		return "", err
	}
	err = ValidateFileName(outputFileName)
	if err != nil {
		return "", err
	}
	err = ValidateDirPath(outputDir)
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
