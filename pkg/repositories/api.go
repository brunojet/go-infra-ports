// Package repositories exposes the public repositories API and REST registry facade.
package repositories

import (
	"github.com/brunojet/go-infra-ports/internal/repositories/rest"
)

const (
	// DefaultStatusCode is the default key used by registry response mappings.
	DefaultStatusCode = rest.DefaultStatusCode
)

// RestRegistry aliases the public REST registry contract.
type RestRegistry = rest.RestRegistry

// RestRequestSpec aliases the request payload contract used by RestRegistry.
type RestRequestSpec = rest.RestRequestSpec

// RestResponseSpec aliases the response payload contract used by RestRegistry.
type RestResponseSpec = rest.RestResponseSpec

// RestEnvelopeSpec aliases the envelope payload contract used by RestRegistry.
type RestEnvelopeSpec = rest.RestEnvelopeSpec

// RestRequest aliases REST request DTO used by RestRepository ports.
type RestRequest = rest.RestRequest

// RestResponse aliases REST single response DTO used by RestRepository ports.
type RestResponse = rest.RestResponse

// RestResponses aliases REST list response DTO used by RestRepository ports.
type RestResponses = rest.RestResponses

// RestRepository aliases the outbound REST repository port contract.
type RestRepository = rest.RestRepository

// RegistryOption configures RestRegistry construction.
type RegistryOption = rest.RegistryOption

// DefaultRestRequest provides a raw-body request spec default implementation.
type DefaultRestRequest = rest.DefaultRestRequest

// DefaultRestResponse provides a raw-body response spec default implementation.
type DefaultRestResponse = rest.DefaultRestResponse

// RestMethod aliases the path method type used for RestRepository path resolution.
type RestMethod = rest.RestMethod

// NewRestRegistry builds a REST registry with the provided options.
func NewRestRegistry(options ...RegistryOption) RestRegistry {
	return rest.NewRestRegistry(options...)
}

// WithRequest registers request specs for methods.
func WithRequest(spec RestRequestSpec, methods rest.RestMethod) RegistryOption {
	return rest.WithRequest(spec, methods)
}

// WithRequestEnvelope registers request envelope specs for methods.
func WithRequestEnvelope(spec RestRequestSpec, methods rest.RestMethod) RegistryOption {
	return rest.WithRequestEnvelope(spec, methods)
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

// HttpClient is the outbound port contract for executing HTTP requests.
type HttpClient = rest.HttpClient

// RepositoryOption configures RestRepository construction.
type RepositoryOption = rest.RepositoryOption

// NewRestRepository builds a RestRepository with the provided options.
// HttpClient is required; panics at startup if missing or if any option is invalid.
func NewRestRepository(opts ...RepositoryOption) RestRepository {
	return rest.NewRestRepository(opts...)
}

// WithHttpClient sets the HttpClient used for HTTP transport. Required.
func WithHttpClient(c HttpClient) RepositoryOption {
	return rest.WithHttpClient(c)
}

// WithRegistryOpt sets the RestRegistry used for request/response marshaling. Required.
func WithRegistryOpt(r RestRegistry) RepositoryOption {
	return rest.WithRegistry(r)
}

// WithBaseURL sets the base URL prepended to all operation paths.
func WithBaseURL(url string) RepositoryOption {
	return rest.WithBasePath(url)
}

// WithPath registers the URL path template for a given operation.
// The required identifiers are extracted from the template and validated at runtime.
func WithPath(methods RestMethod, pathTemplate string) RepositoryOption {
	return rest.WithPath(methods, pathTemplate)
}

// WithRepositoryHeader adds a default request header sent on every operation.
func WithRepositoryHeader(key, value string) RepositoryOption {
	return rest.WithHeader(key, value)
}
