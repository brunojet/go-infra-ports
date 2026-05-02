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

func (r *restRegistry) resolveRequestSpec(name string) RestRequestSpec {
	if spec, ok := r.cfg.Requests[name]; ok && spec != nil {
		return spec
	}
	return r.cfg.Requests[DefaultMethodName]
}

func (r *restRegistry) resolveRequestEnvelopeSpec(name string) RestRequestSpec {
	if spec, ok := r.cfg.RequestsEnvelopes[name]; ok && spec != nil {
		return spec
	}
	if spec, ok := r.cfg.RequestsEnvelopes[DefaultMethodName]; ok && spec != nil {
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
