package rest

import (
	"errors"
	"fmt"
)

var (
	errRepositoryMissingHttpClient       = errors.New("rest repository: HttpClient must not be nil")
	errRepositoryMissingRegistry         = errors.New("rest repository: RestRegistry must not be nil")
	errRepositoryPathMethodNotConfigured = errors.New("rest repository: no path template configured for path method")
	errRepositoryRequestBodyNil          = errors.New("rest repository: request body must not be nil")
	errRepositoryNilHTTPResponse         = errors.New("rest repository: HttpClient returned nil response")
)

// (removed unused path/template error helpers)

func errRepositoryInvalidBasePath(err error) error {
	return fmt.Errorf("invalid basePath: %w", err)
}

func errRepositoryPathMethodNotConfiguredf(pathMethod RestMethod) error {
	return fmt.Errorf("%w: %d", errRepositoryPathMethodNotConfigured, pathMethod)
}

func errRepositoryRequestBodyNilf(pathMethod RestMethod) error {
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
