package main

import (
	"errors"
	"fmt"
)

const (
	ActionValidate = "validate"
)

type PathError struct {
	Action string
	Path   string
	Err    error
}

func (pe *PathError) Error() string {
	return fmt.Sprintf("%s %s: %v", pe.Action, pe.Path, pe.Err)
}

func pathError(action, path string, err error) error {
	return &PathError{
		Action: action,
		Path:   path,
		Err:    err,
	}
}

var (
	ErrForbiddenChars = errors.New("file name must not contain forbidden characters")
)
