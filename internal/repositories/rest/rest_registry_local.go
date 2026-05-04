package rest

import (
	"net/http"
)

func (r *restRegistry) cloneConfig() *registryOptions {
	cloned := newRegistryOptions()
	mergeRegistryOptions(cloned, r.cfg)
	return cloned
}

func (r *restRegistry) resolveResponseSpec(status int) (RestResponseSpec, error) {
	switch {
	case status >= http.StatusContinue && status <= http.StatusContinue+99:
		return resolveResponseSpec(status, r.cfg.Informations), nil
	case status >= http.StatusOK && status <= http.StatusOK+99:
		return resolveResponseSpec(status, r.cfg.Responses), nil
	case status >= http.StatusMultipleChoices && status <= http.StatusMultipleChoices+99:
		return resolveResponseSpec(status, r.cfg.Redirections), nil
	case status >= http.StatusBadRequest && status <= http.StatusBadRequest+199:
		return resolveResponseSpec(status, r.cfg.Problems), nil
	default:
		return nil, errRestResolveResponseSpecNil
	}
}

func (r *restRegistry) newResponseSpec(status int) (RestResponseSpec, error) {
	prototype, err := r.resolveResponseSpec(status)
	if err != nil {
		return nil, err
	}
	instance := prototype.New()
	if instance == nil {
		return nil, errRestResolveResponseSpecNewNil
	}
	return instance, nil
}

func (r *restRegistry) newRequestSpec(restMethod RestMethod) (RestRequestSpec, error) {
	spec, ok := r.cfg.Requests[restMethod]
	if !ok {
		return nil, errRestNewRequestSpecNotFound(restMethod)
	}
	instance := spec.New()
	if instance == nil {
		return nil, errRestNewRequestSpecNil
	}
	return instance, nil
}

func (r *restRegistry) resolveRequestEnvelopeSpec(restMethod RestMethod) RestRequestSpec {
	if spec, ok := r.cfg.RequestsEnvelopes[restMethod]; ok && spec != nil {
		return spec
	}
	return nil
}

func (r *restRegistry) resolveResponseEnvelopeSpec(status int) RestEnvelopeSpec {
	if spec, ok := r.cfg.ResponseEnvelopes[status]; ok && spec != nil {
		return spec
	}
	if spec, ok := r.cfg.ResponseEnvelopes[DefaultStatusCode]; ok && spec != nil {
		return spec
	}
	return nil
}
