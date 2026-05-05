package http_util

import (
	"errors"
	"fmt"
)

var (
	errInvalidPathTemplate            = errors.New("invalid path template")
	errRepositoryBaseURLEmpty         = errors.New("base URL must not be empty")
	errRepositoryPathInvalidChars     = errors.New("path has invalid characters")
	errRepositoryPathInvalidStructure = errors.New("path has invalid structure")
)

func errRepositoryPathInvalidCharsf(path string) error {
	return fmt.Errorf("%w: %q", errRepositoryPathInvalidChars, path)
}

func errRepositoryPathInvalidStructuref(path string) error {
	return fmt.Errorf("%w: %q", errRepositoryPathInvalidStructure, path)
}

func errRepositoryInvalidPathTemplate(err error) error {
	return fmt.Errorf("%w: %v", errInvalidPathTemplate, err)
}
