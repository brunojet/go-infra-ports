package rest

import (
	"errors"
	"fmt"
)

var (
	errRestServiceRepositoryNil       = errors.New("rest service repository is nil")
	errRestServiceResponseNil         = errors.New("rest service response is nil")
	errRestServiceInvalidNon2xxStatus = errors.New("invalid non-2xx status code")
	errRestServiceNilResponseData     = errors.New("rest response data is nil for 2xx status code")
)

func errRestServiceUpstreamMappingFailed(operation string, err error) error {
	return fmt.Errorf("rest service upstream mapping failed at %s: %w", operation, err)
}

func errRestServiceDownstreamMappingFailed(operation string, err error) error {
	return fmt.Errorf("rest service downstream mapping failed at %s: %w", operation, err)
}
