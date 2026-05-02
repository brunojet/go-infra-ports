package rest

import (
	"net/http"

	restcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
)

var (
	validRequestMethods = map[string]struct{}{
		DefaultMethodName: {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodPatch:  {},
	}
)

type registryOptions struct {
	Requests          map[string]RestRequestSpec
	RequestsEnvelopes map[string]RestRequestSpec
	Responses         map[int]RestResponseSpec
	ResponseEnvelopes map[int]RestEnvelopeSpec
	Informations      map[int]RestResponseSpec
	Redirections      map[int]RestResponseSpec
	Problems          map[int]RestResponseSpec
}

type RegistryOption func(*registryOptions)

func newRegistryOptions() *registryOptions {
	return &registryOptions{
		Requests:          map[string]RestRequestSpec{DefaultMethodName: &DefaultRestRequest{}},
		RequestsEnvelopes: make(map[string]RestRequestSpec),
		Responses:         map[int]RestResponseSpec{DefaultStatusCode: &DefaultRestResponse{}},
		ResponseEnvelopes: make(map[int]RestEnvelopeSpec),
		Informations:      map[int]RestResponseSpec{DefaultStatusCode: &DefaultRestResponse{}},
		Redirections:      map[int]RestResponseSpec{DefaultStatusCode: &DefaultRestResponse{}},
		Problems:          map[int]RestResponseSpec{DefaultStatusCode: &DefaultRestResponse{}},
	}
}

func newRegistryConfig(options ...RegistryOption) *registryOptions {
	cfg := newRegistryOptions()
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func (o *registryOptions) registerRequest(spec RestRequestSpec, target map[string]RestRequestSpec, methods ...string) {
	if spec == nil {
		panic(errRestRegisterRequestSpecNil)
	}
	if len(methods) == 0 {
		methods = []string{DefaultMethodName}
	}
	for _, method := range methods {
		if _, ok := validRequestMethods[method]; !ok {
			panic(errRestRegisterRequestInvalidMethod(method))
		}
		target[method] = spec
	}
}

func (o *registryOptions) registerResponse(spec RestResponseSpec, target map[int]RestResponseSpec, firstStatusCode, lastStatusCode int, statusCodes ...int) {
	if spec == nil {
		panic(errRestRegisterResponseSpecNil)
	}
	if len(statusCodes) == 0 {
		statusCodes = []int{DefaultStatusCode}
	}
	for _, code := range statusCodes {
		if code == DefaultStatusCode {
			target[code] = spec
			continue
		}
		if code < firstStatusCode || code > lastStatusCode {
			panic(errRestRegisterResponseOutOfRange(code, firstStatusCode, lastStatusCode))
		}
		target[code] = spec
	}
}

func (o *registryOptions) registerResponseEnvelope(spec restcts.RestEnvelopeSpec, statusCodes ...int) {
	if spec == nil {
		panic(errRestRegisterResponseSpecNil)
	}
	if len(statusCodes) == 0 {
		o.ResponseEnvelopes[DefaultStatusCode] = spec
		return
	}
	for _, code := range statusCodes {
		if code != DefaultStatusCode && (code < http.StatusOK || code > http.StatusOK+99) {
			panic(errRestRegisterResponseOutOfRange(code, http.StatusOK, http.StatusOK+99))
		}
		o.ResponseEnvelopes[code] = spec
	}
}

func WithRequest(spec RestRequestSpec, methods ...string) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerRequest(spec, ro.Requests, methods...)
	}
}

func WithRequestEnvelope(spec RestRequestSpec, methods ...string) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerRequest(spec, ro.RequestsEnvelopes, methods...)
	}
}

func WithResponse(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerResponse(spec, ro.Responses, http.StatusOK, http.StatusOK+99, statusCodes...)
	}
}

func WithResponseEnvelope(spec restcts.RestEnvelopeSpec, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerResponseEnvelope(spec, statusCodes...)
	}
}

func WithInformation(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerResponse(spec, ro.Informations, http.StatusContinue, http.StatusContinue+99, statusCodes...)
	}
}

func WithRedirection(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerResponse(spec, ro.Redirections, http.StatusMultipleChoices, http.StatusMultipleChoices+99, statusCodes...)
	}
}

func WithProblem(spec RestResponseSpec, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		ro.registerResponse(spec, ro.Problems, http.StatusBadRequest, http.StatusBadRequest+199, statusCodes...)
	}
}
