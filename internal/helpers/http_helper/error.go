package http_helper

import (
	"errors"
	"fmt"
)

var (
	errNilURL                  = errors.New("url must not be nil")
	errNilDstHeader            = errors.New("dst header must not be nil")
	errRepositoryBasePathEmpty = errors.New("base path must not be empty")
	errPathTooLong             = errors.New("path length exceeds maximum")
	errPathInvalidChars        = errors.New("path has invalid characters")
	errPathInvalidStructure    = errors.New("path has invalid structure")
	errInvalidPathTemplateErr  = errors.New("invalid path template")
	errPathParametersExceedMax = errors.New("number of path parameters exceeds maximum")
)

func errPathTooLongf(p string, maxAllowed int) error {
	return fmt.Errorf("%w of %d: %q", errPathTooLong, maxAllowed, p)
}

func errPathInvalidCharsf(p string) error {
	return fmt.Errorf("%w: %q", errPathInvalidChars, p)
}

func errPathInvalidStructuref(p string) error {
	return fmt.Errorf("%w: %q", errPathInvalidStructure, p)
}

func errInvalidPathTemplate(err error) error {
	return fmt.Errorf("%w: %w", errInvalidPathTemplateErr, err)
}

func errPathParametersExceedMaxf(count, maxAllowed int) error {
	return fmt.Errorf("%w of %d: %d", errPathParametersExceedMax, maxAllowed, count)
}
