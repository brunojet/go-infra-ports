// Package contracts defines public REST service mapping contracts.
package contracts

import (
	"net/http"
	"net/url"

	rpocts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
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
	// RestRequestSpec aliases the repository request payload contract for use in REST mapper signatures.
	RestRequestSpec = rpocts.RestRequestSpec
	// RestResponseSpec aliases the repository response payload contract for use in REST mapper signatures.
	RestResponseSpec = rpocts.RestResponseSpec
)

// RestUpstreamMapper maps downstream request data to upstream request contracts.
type RestUpstreamMapper[C, U any] interface {
	ToUpstreamPost(payload C, ids Identifiers, upsPayload *RestRequestSpec) error
	ToUpstreamPut(payload C, ids Identifiers, upsPayload *RestRequestSpec) error
	ToUpstreamPatch(payload U, ids Identifiers, upsPayload *RestRequestSpec) error
	ToUpstreamQuery(reqQuery, upsQuery url.Values) error
	ToUpstreamHeaders(reqHeader, upsHeader http.Header) error
}

// RestDownstreamMapper maps upstream responses into downstream response contracts.
type RestDownstreamMapper[R any] interface {
	ToDownstreamStatusCode(statusCode int, downstreamStatusCode *int) error
	ToDownstreamHeaders(upsHeader, downstreamHeader http.Header) error
	ToDownstreamResponse(upsPayload any, payload *R) error
	ToDownstreamResponseMeta(meta ResponseMeta, serviceMeta *ServiceMeta) error
	ToDownstreamInformation(statusCode int, upsPayload RestResponseSpec, serviceMeta *ServiceMeta) error
	ToDownstreamRedirection(statusCode int, upsPayload RestResponseSpec, serviceMeta *ServiceMeta) error
	ToDownstreamProblem(statusCode int, upsPayload RestResponseSpec, serviceMeta *ServiceMeta) error
}
