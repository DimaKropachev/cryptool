package file

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	Decrpt = "Decrypt"
	Encrpt = "Encypt"
)

// CreateOutputFilePath создает путь к выходноу файлу.
// Если outPath = "", то функция вернет путь к входному файлу только с подписью ".crpt".
// Если outPath указывает на директорию, то вернется значение этой директории в которой будет файл с именем входного только с подписью ".crpt".
// Если outPath уже содержит путь к файлу, то CreateOutputFile вернет его без изменений.
func CreateOutputFilePath(action string, inPath, outPath string) string {
	// expansion - расширение
	var expansion string
	switch action {
	case Decrpt:
		expansion = ".dec"
	case Encrpt:
		expansion = ".enc"
	}

	if outPath == "" {
		return inPath + expansion
	}

	outPathDir, outPathFile := filepath.Split(outPath)
	switch outPathFile {
	case "":
		_, inputFile := filepath.Split(inPath)
		return outPathDir + string(filepath.Separator) + inputFile + expansion
	case ".":
		_, inputFile := filepath.Split(inPath)
		return "." + string(filepath.Separator) + inputFile + expansion
	}
	return outPath
}

// CreateOuputDirPath создает путь к выходной директории.
// Если outPath = "", то вернется значение пути к входной директории с подписью (crpt).
// Если outPath содержит путь к конкретному файлу, то будет возвращен только значение директории из этого пути.
func CreateOutputDirPath(inPath, outPath string) (string, error) {
	var newPath string

	if outPath == "" {
		newPath = inPath + "(crpt)"
	} else {
		newPath = outPath
	}

	err := os.MkdirAll(newPath, 0644)
	if err != nil {
		return "", fmt.Errorf("error creating output directory: %w", err)
	}
	return newPath, nil
}
