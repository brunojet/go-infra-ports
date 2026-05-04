package rest

import (
	"errors"
	"fmt"
)

var (
	errRepositoryMissingHttpClient       = errors.New("rest repository: HttpClient must not be nil")
	errRepositoryMissingRegistry         = errors.New("rest repository: RestRegistry must not be nil")
	errRepositoryBaseURLEmpty            = errors.New("rest repository: base URL must not be empty")
	errRepositoryPathInvalidChars        = errors.New("rest repository: path has invalid characters")
	errRepositoryPathInvalidStructure    = errors.New("rest repository: path has invalid structure")
	errRepositoryPathMethodNotConfigured = errors.New("rest repository: no path template configured for path method")
	errRepositoryRequestBodyNil          = errors.New("rest repository: request body must not be nil")
	errRepositoryNilHTTPResponse         = errors.New("rest repository: HttpClient returned nil response")
	errRepositoryCollectionPathHasIDs    = errors.New("rest repository: collection path must not contain {key} placeholders; use WithPath for parameterized paths")
	errRepositoryInstancePathMissingID   = errors.New("rest repository: instance path must contain at least one {key} placeholder; use WithPath for non-parameterized paths")
)

func errRepositoryPathInvalidCharsf(path string) error {
	return fmt.Errorf("%w: %q", errRepositoryPathInvalidChars, path)
}

func errRepositoryPathInvalidStructuref(path string) error {
	return fmt.Errorf("%w: %q", errRepositoryPathInvalidStructure, path)
}

func errRepositoryInvalidPathTemplate(err error) error {
	return fmt.Errorf("invalid path template: %w", err)
}

func errRepositoryInvalidBasePath(err error) error {
	return fmt.Errorf("invalid basePath: %w", err)
}

func errRepositoryPathMethodNotConfiguredf(pathMethod PathMethod) error {
	return fmt.Errorf("%w: %d", errRepositoryPathMethodNotConfigured, pathMethod)
}

func errRepositoryRequestBodyNilf(pathMethod PathMethod) error {
	return fmt.Errorf("%w for operation %d", errRepositoryRequestBodyNil, pathMethod)
}

func errRepositoryBuildRequest(err error) error {
	return fmt.Errorf("rest repository: build request: %w", err)
}

func errRepositoryExecuteRequest(err error) error {
	return fmt.Errorf("rest repository: execute request: %w", err)
}

func errRepositoryReadResponseBody(err error) error {
	return fmt.Errorf("rest repository: read response body: %w", err)
}

func errRepositoryResolveRequest(err error) error {
	return fmt.Errorf("rest repository: resolve request: %w", err)
}
func errRepositoryResolveEnvelopeRequest(err error) error {
	return fmt.Errorf("rest repository: resolve envelope request: %w", err)
}

func errRepositoryResolveEnvelopeResponse(err error) error {
	return fmt.Errorf("rest repository: resolve envelope response: %w", err)
}

func errRepositoryResolveResponse(err error) error {
	return fmt.Errorf("rest repository: resolve response: %w", err)
}
