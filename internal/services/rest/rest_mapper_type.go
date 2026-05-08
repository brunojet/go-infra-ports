package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/brunojet/go-infra-ports/internal/helpers/http_helper"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type (
	Identifiers                  = types.Identifiers
	RestUpstreamMapper[C, U any] = svccts.RestUpstreamMapper[C, U]
	RestDownstreamMapper[R any]  = svccts.RestDownstreamMapper[R]
)

type DefaultRestUpstreamMapper[C, U any] struct{}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPost(payload C, ids Identifiers, upsPayload *any) error {
	if upsPayload == nil {
		return fmt.Errorf("%s: upsPayload cannot be nil", "ToUpstreamPost")
	}
	*upsPayload = payload
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPut(payload C, ids Identifiers, upsPayload *any) error {
	if upsPayload == nil {
		return fmt.Errorf("%s: upsPayload cannot be nil", "ToUpstreamPut")
	}
	*upsPayload = payload
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamPatch(payload U, ids Identifiers, upsPayload *any) error {
	if upsPayload == nil {
		return fmt.Errorf("%s: upsPayload cannot be nil", "ToUpstreamPatch")
	}
	*upsPayload = payload
	return nil
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamQuery(reqQuery, upsQuery url.Values) error {
	return http_helper.ApplyQueryParams(upsQuery, reqQuery)
}

func (m *DefaultRestUpstreamMapper[C, U]) ToUpstreamHeaders(reqHeader, upsHeader http.Header) error {
	return http_helper.ApplyHeaderParams(upsHeader, reqHeader)
}

type DefaultRestDownstreamMapper[R any] struct{}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamStatusCode(statusCode int, dwsStatusCode *int) error {
	if dwsStatusCode == nil {
		return fmt.Errorf("ToDownstreamStatusCode: dwsStatusCode cannot be nil")
	}
	*dwsStatusCode = statusCode
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamHeaders(upsHeader, dwsHeader http.Header) error {
	return http_helper.ApplyHeaderParams(dwsHeader, upsHeader)
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamResponse(upsPayload any, payload *R) error {
	if payload == nil {
		return fmt.Errorf("ToDownstreamResponse: payload cannot be nil")
	}
	// Simple, strict behavior: upstream payload must be exactly the same
	// type as the downstream generic `R`. No reflection or pointer/value
	// coercion is performed — if types differ, return an error.
	if casted, ok := upsPayload.(R); ok {
		*payload = casted
		return nil
	}

	return fmt.Errorf("ToDownstreamResponse: failed to cast upsPayload to target type")
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamResponseMeta(_ svccts.ResponseMeta, _ *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamInformation(statusCode int, upsPayload any, serviceMeta *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamRedirection(statusCode int, upsPayload any, serviceMeta *svccts.ServiceMeta) error {
	return nil
}

func (m *DefaultRestDownstreamMapper[R]) ToDownstreamProblem(statusCode int, upsPayload any, serviceMeta *svccts.ServiceMeta) error {
	return nil
}
