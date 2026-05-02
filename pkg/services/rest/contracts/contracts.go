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
)

// RestUpstreamMapper maps downstream request data to upstream request contracts.
type RestUpstreamMapper[C, U any] interface {
	ToUpstreamPost(payload C, ids Identifiers, upsPayload *any) error
	ToUpstreamPut(payload U, ids Identifiers, upsPayload *any) error
	ToUpstreamPatch(payload U, ids Identifiers, upsPayload *any) error
	ToUpstreamIdentifiers(reqIds Identifiers, pathIds *Identifiers, upsIds *Identifiers) error
	ToUpstreamQuery(reqQuery url.Values, upsQuery *url.Values) error
	ToUpstreamHeaders(reqHeader http.Header, upsHeader *http.Header) error
}

// RestDownstreamMapper maps upstream responses into downstream response contracts.
type RestDownstreamMapper[R any] interface {
	ToDownstreamStatusCode(statusCode int, downstreamStatusCode *int) error
	ToDownstreamHeaders(upsHeader http.Header, downstreamHeader *http.Header) error
	ToDownstreamResponse(upsPayload any, payload *R) error
	ToDownstreamInformation(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
	ToDownstreamRedirection(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
	ToDownstreamProblem(statusCode int, upsPayload any, serviceMeta *ServiceMeta) error
}
