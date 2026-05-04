package rest

import (
	"maps"
	"net/http"
	"net/url"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type (
	Identifiers                  = types.Identifiers
	RestUpstreamMapper[C, U any] = svccts.RestUpstreamMapper[C, U]
	RestDownstreamMapper[R any]  = svccts.RestDownstreamMapper[R]
)

type DefaultRestUpstreamMapper[C, U any] struct{}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPost(payload C, ids Identifiers, upsPayload *svccts.RestRequestSpec) error {
	if spec, ok := any(payload).(svccts.RestRequestSpec); ok {
		*upsPayload = spec
	}
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPut(payload C, ids Identifiers, upsPayload *svccts.RestRequestSpec) error {
	if spec, ok := any(payload).(svccts.RestRequestSpec); ok {
		*upsPayload = spec
	}
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPatch(payload U, ids Identifiers, upsPayload *svccts.RestRequestSpec) error {
	if spec, ok := any(payload).(svccts.RestRequestSpec); ok {
		*upsPayload = spec
	}
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamQuery(reqQuery, upsQuery url.Values) error {
	maps.Copy(upsQuery, reqQuery)
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamHeaders(reqHeader, upsHeader http.Header) error {
	maps.Copy(upsHeader, reqHeader)
	return nil
}

type DefaultRestDownstreamMapper[R any] struct{}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamStatusCode(statusCode int, downstreamStatusCode *int) error {
	*downstreamStatusCode = statusCode
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamHeaders(upsHeader, downstreamHeader http.Header) error {
	maps.Copy(downstreamHeader, upsHeader)
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamResponse(upsPayload any, payload *R) error {
	// Try direct type assertion first (for value types)
	if casted, ok := upsPayload.(R); ok {
		*payload = casted
		return nil
	}

	// Try pointer type assertion (for pointer types)
	castedPtr, ok := upsPayload.(*R)
	if !ok {
		return errRestDownstreamMapperResponseTypeAssertion
	}
	*payload = *castedPtr
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamInformation(statusCode int, upsPayload RestResponseSpec, serviceMeta *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamResponseMeta(_ svccts.ResponseMeta, _ *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamRedirection(statusCode int, upsPayload RestResponseSpec, serviceMeta *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamProblem(statusCode int, upsPayload RestResponseSpec, serviceMeta *svccts.ServiceMeta) error {
	return nil
}
