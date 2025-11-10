package file

import (
	"errors"
	"fmt"
)

type Action string

const (
	ActionValidate = "validate"
	ActionScanning = "scanning"
)

var (
	ErrEmptyPath = errors.New("path cannot be empty")

	ErrSepInvalidSyntax         = errors.New("separator invalid syntax")
	ErrEmptyDirPath          = errors.New("directory path cannot be empty")
	ErrForbiddenCharsDirPath = errors.New("directory path cannot contain forbidden characters")
	ErrFolderDotsDirPath     = errors.New("directory path cannot contain a folder with a name consisting only dots")

	ErrEmptyFileName          = errors.New("file name cannot be empty")
	ErrForbiddenCharsFileName = errors.New("file name cannot contain forbidden characters")
	ErrDotsFileName           = errors.New("file name cannot contain a folder with a name consisting only dots")

	ErrInvalidFileExtension = errors.New("invalid file invalid")
)

type PathError struct {
	Action Action
	Path   string
	Err    error
}

func (pe *PathError) Error() string {
	return fmt.Sprintf("%s [%s]: %v", pe.Action, pe.Path, pe.Err)
}

func pathError(action Action, path string, err error) error {
	return &PathError{
		Action: action,
		Path:   path,
		Err:    err,
	}
}
