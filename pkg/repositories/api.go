// Package repositories exposes the public repositories API and REST registry facade.
package repositories

import (
	rest "github.com/brunojet/go-infra-ports/internal/repositories/rest"
	restcontracts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
)

const (
	// DefaultMethodName is the default key used by registry request mappings.
	DefaultMethodName = rest.DefaultMethodName
	// DefaultStatusCode is the default key used by registry response mappings.
	DefaultStatusCode = rest.DefaultStatusCode
)

// RestRegistry aliases the public REST registry contract.
type RestRegistry = restcontracts.RestRegistry

// RestRequestSpec aliases the request payload contract used by RestRegistry.
type RestRequestSpec = restcontracts.RestRequestSpec

// RestResponseSpec aliases the response payload contract used by RestRegistry.
type RestResponseSpec = restcontracts.RestResponseSpec

// RestEnvelopeSpec aliases the envelope payload contract used by RestRegistry.
type RestEnvelopeSpec = restcontracts.RestEnvelopeSpec

// RestRequest aliases REST request DTO used by RestRepository ports.
type RestRequest = restcontracts.RestRequest

// RestResponse aliases REST single response DTO used by RestRepository ports.
type RestResponse = restcontracts.RestResponse

// RestResponses aliases REST list response DTO used by RestRepository ports.
type RestResponses = restcontracts.RestResponses

// RestRepository aliases the outbound REST repository port contract.
type RestRepository = restcontracts.RestRepository

// RegistryOption configures RestRegistry construction.
type RegistryOption = rest.RegistryOption

// DefaultRestRequest provides a raw-body request spec default implementation.
type DefaultRestRequest = rest.DefaultRestRequest

// DefaultRestResponse provides a raw-body response spec default implementation.
type DefaultRestResponse = rest.DefaultRestResponse

// NewRestRegistry builds a REST registry with the provided options.
func NewRestRegistry(options ...RegistryOption) RestRegistry {
	return rest.NewRestRegistry(options...)
}

// WithRequest registers request specs for methods (or default method when empty).
func WithRequest(spec RestRequestSpec, methods ...string) RegistryOption {
	return rest.WithRequest(spec, methods...)
}

// WithRequestEnvelope registers request envelope specs for methods.
func WithRequestEnvelope(spec RestRequestSpec, methods ...string) RegistryOption {
	return rest.WithRequestEnvelope(spec, methods...)
}

// WithResponse registers 2xx response specs (or default status when empty).
func WithResponse(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return rest.WithResponse(spec, statusCodes...)
}

// WithResponseEnvelope registers envelope specs for 2xx responses.
func WithResponseEnvelope(spec RestEnvelopeSpec, statusCodes ...int) RegistryOption {
	return rest.WithResponseEnvelope(spec, statusCodes...)
}

// WithInformation registers 1xx response specs.
func WithInformation(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return rest.WithInformation(spec, statusCodes...)
}

// WithRedirection registers 3xx response specs.
func WithRedirection(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return rest.WithRedirection(spec, statusCodes...)
}

// WithProblem registers 4xx and 5xx response specs.
func WithProblem(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return rest.WithProblem(spec, statusCodes...)
}
