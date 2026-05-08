package rest

import (
	"net/http"
)

func (r *restRegistry) cloneConfig() *registryOptions {
	cloned := newRegistryOptions()
	mergeRegistryOptions(cloned, r.cfg)
	return cloned
}

func (r *restRegistry) resolveResponseSpec(status int) (RestDataSpec, error) {
	switch {
	case status >= http.StatusContinue && status <= http.StatusContinue+99:
		return resolveResponseSpec(status, r.cfg.informations), nil
	case status >= http.StatusOK && status <= http.StatusOK+99:
		return resolveResponseSpec(status, r.cfg.responses), nil
	case status >= http.StatusMultipleChoices && status <= http.StatusMultipleChoices+99:
		return resolveResponseSpec(status, r.cfg.redirections), nil
	case status >= http.StatusBadRequest && status <= http.StatusBadRequest+199:
		return resolveResponseSpec(status, r.cfg.problems), nil
	default:
		return nil, errRestResolveResponseSpecNil
	}
}

func (r *restRegistry) newResponseSpec(status int) (RestDataSpec, error) {
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

func (r *restRegistry) newRequestSpec(restMethod RestMethod) (RestDataSpec, error) {
	spec, ok := r.cfg.requests[restMethod]
	if !ok {
		return nil, errRestNewRequestSpecNotFound(restMethod)
	}
	instance := spec.New()
	if instance == nil {
		return nil, errRestNewRequestSpecNil
	}
	return instance, nil
}

func (r *restRegistry) resolveRequestEnvelopeSpec(restMethod RestMethod) RestEnvelopeSpec {
	if spec, ok := r.cfg.requestsEnvelopes[restMethod]; ok && spec != nil {
		return spec
	}
	return nil
}

func (r *restRegistry) resolveResponseEnvelopeSpec(status int) RestEnvelopeSpec {
	if spec, ok := r.cfg.responseEnvelopes[status]; ok && spec != nil {
		return spec
	}
	if spec, ok := r.cfg.responseEnvelopes[defaultStatusCode]; ok && spec != nil {
		return spec
	}
	return nil
}
