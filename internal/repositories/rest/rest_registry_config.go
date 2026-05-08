package rest

import (
	"net/http"
)

const defaultStatusCode = 0

type registryOptions struct {
	requests          map[RestMethod]RestDataSpec
	requestsEnvelopes map[RestMethod]RestEnvelopeSpec
	responses         map[int]RestDataSpec
	responseEnvelopes map[int]RestEnvelopeSpec
	informations      map[int]RestDataSpec
	redirections      map[int]RestDataSpec
	problems          map[int]RestDataSpec
}

type RegistryOption func(*registryOptions)

func newRegistryOptions() *registryOptions {
	reqDefault := NewDataSpecOf[DefaultRestRequest]()
	respDefault := NewDataSpecOf[DefaultRestResponse]()
	return &registryOptions{
		requests: map[RestMethod]RestDataSpec{
			MethodCreate: reqDefault,
			MethodUpdate: reqDefault,
			MethodSave:   reqDefault,
		},
		requestsEnvelopes: make(map[RestMethod]RestEnvelopeSpec),
		responses:         map[int]RestDataSpec{defaultStatusCode: respDefault},
		responseEnvelopes: make(map[int]RestEnvelopeSpec),
		informations:      map[int]RestDataSpec{defaultStatusCode: respDefault},
		redirections:      map[int]RestDataSpec{defaultStatusCode: respDefault},
		problems:          map[int]RestDataSpec{defaultStatusCode: respDefault},
	}
}

func newRegistryConfig(options ...RegistryOption) *registryOptions {
	cfg := newRegistryOptions()
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func (o *registryOptions) registerRequest(spec RestDataSpec, target map[RestMethod]RestDataSpec, methods RestMethod) {
	if spec == nil {
		panic(errRestRegisterRequestSpecNil)
	}
	if methods&allWriteMethods == 0 {
		panic(errRestRegisterRequestInvalidMethod(methods))
	}
	for _, method := range writeMethodsList {
		if methods&method != 0 {
			target[method] = spec
		}
	}
}

func (o *registryOptions) registerResponse(spec RestDataSpec, target map[int]RestDataSpec, firstStatusCode, lastStatusCode int, statusCodes ...int) {
	if spec == nil {
		panic(errRestRegisterResponseSpecNil)
	}
	if len(statusCodes) == 0 {
		statusCodes = []int{defaultStatusCode}
	}
	for _, code := range statusCodes {
		if code == defaultStatusCode {
			target[code] = spec
			continue
		}
		if code < firstStatusCode || code > lastStatusCode {
			panic(errRestRegisterResponseOutOfRange(code, firstStatusCode, lastStatusCode))
		}
		target[code] = spec
	}
}

func (o *registryOptions) registerResponseEnvelope(spec RestEnvelopeSpec, statusCodes ...int) {
	if spec == nil {
		panic(errRestRegisterResponseSpecNil)
	}
	if len(statusCodes) == 0 {
		o.responseEnvelopes[defaultStatusCode] = spec
		return
	}
	for _, code := range statusCodes {
		if code != defaultStatusCode && (code < http.StatusOK || code > http.StatusOK+99) {
			panic(errRestRegisterResponseOutOfRange(code, http.StatusOK, http.StatusOK+99))
		}
		o.responseEnvelopes[code] = spec
	}
}

func WithRequestOf[T any](methods RestMethod) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewDataSpecOf[T]()
		ro.registerRequest(spec, ro.requests, methods)
	}
}

func WithRequestEnvelope(dataField string, methods RestMethod) RegistryOption {
	return func(ro *registryOptions) {
		if methods&allWriteMethods == 0 {
			panic(errRestRegisterRequestInvalidMethod(methods))
		}
		spec := NewEnvelopeSpec(dataField, "")
		for _, method := range writeMethodsList {
			if methods&method != 0 {
				ro.requestsEnvelopes[method] = spec
			}
		}
	}
}

func WithResponseOf[T any](statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewDataSpecOf[T]()
		ro.registerResponse(spec, ro.responses, http.StatusOK, http.StatusOK+99, statusCodes...)
	}
}

func WithResponseEnvelope(dataField, metaField string, statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewEnvelopeSpec(dataField, metaField)
		ro.registerResponseEnvelope(spec, statusCodes...)
	}
}

func WithInformationOf[T any](statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewDataSpecOf[T]()
		ro.registerResponse(spec, ro.informations, http.StatusContinue, http.StatusContinue+99, statusCodes...)
	}
}

func WithRedirectionOf[T any](statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewDataSpecOf[T]()
		ro.registerResponse(spec, ro.redirections, http.StatusMultipleChoices, http.StatusMultipleChoices+99, statusCodes...)
	}
}

func WithProblemOf[T any](statusCodes ...int) RegistryOption {
	return func(ro *registryOptions) {
		spec := NewDataSpecOf[T]()
		ro.registerResponse(spec, ro.problems, http.StatusBadRequest, http.StatusBadRequest+199, statusCodes...)
	}
}
