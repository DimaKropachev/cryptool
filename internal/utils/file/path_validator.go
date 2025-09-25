package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

func validateFilePath(path string) error {
	if path == "" {
		return pathError(ActionValidate, path, ErrEmptyPath)
	}

	dir, file := filepath.Split(path)
	fmt.Println(dir, file)

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

	forbiddenChar := `:*?"<>|`
	dirPath = strings.ReplaceAll(dirPath, "/", string(filepath.Separator))
	folders := strings.Split(dirPath, string(filepath.Separator))
	for i, folder := range folders {
		if i == 0 {
			if folder == "." {
				continue
			}
		}
		if len(folder) == 0 {
			continue
		}

		if strings.ReplaceAll(folder, ".", "") == "" {
			return pathError(ActionValidate, dirPath, ErrFolderDotsDirPath)
		}

		if strings.ContainsAny(folder, forbiddenChar) {
			return pathError(ActionValidate, dirPath, ErrForbiddenCharsDirPath)
		}
	}

	return nil
}
