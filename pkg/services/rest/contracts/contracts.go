// Package contracts defines public REST service mapping contracts.
package contracts

import (
	"net/http"
	"net/url"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type (
	// Identifiers aliases request identifiers used by REST mappers.
	Identifiers = types.Identifiers
	// ServiceMeta aliases shared service metadata propagated to downstream mappings.
	ServiceMeta = svccts.ServiceMeta
	// ResponseMeta aliases transport-level metadata returned with a response (e.g. pagination).
	ResponseMeta = types.ResponseMeta
)

// RestUpstreamMapper maps downstream request data to upstream request contracts.
type RestUpstreamMapper[C, U any] interface {
	ToUpstreamPost(dwsPayload C, ids Identifiers, upsPayload *any) error
	ToUpstreamPut(dwsPayload C, ids Identifiers, upsPayload *any) error
	ToUpstreamPatch(dwsPayload U, ids Identifiers, upsPayload *any) error
	ToUpstreamQuery(dwsQuery, upsQuery url.Values) error
	ToUpstreamHeaders(dwsHeader, upsHeader http.Header) error
}

// RestDownstreamMapper maps upstream responses into downstream response contracts.
type RestDownstreamMapper[R any] interface {
	ToDownstreamStatusCode(upsStatusCode int, dwsStatusCode *int) error
	ToDownstreamHeaders(upsHeader, dwsHeader http.Header) error
	ToDownstreamResponse(upsPayload any, dwsPayload *R) error
	ToDownstreamResponseMeta(meta ResponseMeta, serviceMeta *ServiceMeta) error
	ToDownstreamInformation(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
	ToDownstreamRedirection(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
	ToDownstreamProblem(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
}
