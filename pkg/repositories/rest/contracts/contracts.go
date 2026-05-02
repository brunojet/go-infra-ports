// Package contracts defines public REST repository contracts.
package contracts

import (
	"context"
	"encoding/json"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// RestRequestSpec defines the request payload contract resolved by RestRegistry.
//
// Implementations are responsible for carrying decoded request payload data and
// accepting raw body bytes when needed.
type RestRequestSpec interface {
	// New returns a new empty request spec instance of the same concrete type.
	New() RestRequestSpec
	// SetBody assigns a raw JSON payload to the request spec.
	SetBody(body json.RawMessage)
}

// RestResponseSpec defines the response payload contract resolved by RestRegistry.
//
// Implementations are used for successful, informational, redirection, and
// problem responses.
type RestResponseSpec interface {
	// New returns a new empty response spec instance of the same concrete type.
	New() RestResponseSpec
	// NewSlice returns a typed response spec slice with the requested size.
	NewSlice(n int) []RestResponseSpec
}

// RestEnvelopeSpec defines envelope payload contracts carrying data+meta.
//
// After unmarshaling the raw response into the spec, the registry extracts the
// payload via EnvelopeData and the metadata via EnvelopeMeta.
type RestEnvelopeSpec interface {
	// New returns a new empty envelope spec instance of the same concrete type.
	New() RestEnvelopeSpec
	// EnvelopeData returns the inner data payload extracted from the envelope.
	EnvelopeData() json.RawMessage
	// EnvelopeMeta returns the metadata extracted from the envelope.
	EnvelopeMeta() types.ResponseMeta
}

// RestRegistry defines mapping and codec operations used by RestRepository.
//
// The registry is responsible for transparent marshal and unmarshal behavior,
// including request/response envelopes and metadata extraction.
type RestRegistry interface {
	// Merge combines the current registry with another and returns the merged result.
	Merge(other RestRegistry) RestRegistry
	// ResolveRequest marshals a request spec into requestBody using the provided key.
	ResolveRequest(name string, body RestRequestSpec, requestBody *[]byte) error
	// ResolveEnvelopeRequest envelopes the current request body in place.
	ResolveEnvelopeRequest(name string, dataBody *[]byte) error
	// ResolveResponse unmarshals a single response payload into body.
	ResolveResponse(status int, responseBody []byte, body *RestResponseSpec) error
	// ResolveResponses unmarshals a collection response payload into bodies.
	ResolveResponses(status int, responseBody []byte, bodies *[]RestResponseSpec) error
	// ResolveEnvelopeResponse extracts payload and metadata from an envelope.
	ResolveEnvelopeResponse(status int, dataBody *[]byte, meta *types.ResponseMeta) error
	// NewRequestSpec allocates a request spec instance for service layer usage.
	// Returns an error if the spec's New method returns nil.
	NewRequestSpec(name string) (RestRequestSpec, error)
	// ReleaseRequestSpec releases a request spec instance used by the service layer.
	ReleaseRequestSpec(spec RestRequestSpec)
}

// RestRequest wraps the request body with REST-specific context.
type RestRequest struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the create payload.
	Body RestRequestSpec
}

// RestResponse represents one REST repository output.
type RestResponse struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information RestResponseSpec
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection RestResponseSpec
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem RestResponseSpec
	// Data carries the decoded successful payload.
	Data RestResponseSpec
}

// RestResponses represents many REST repository outputs.
type RestResponses struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information RestResponseSpec
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection RestResponseSpec
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem RestResponseSpec
	// Data carries the decoded successful payload collection.
	Data []RestResponseSpec
}

// RestRepository is the outbound port for HTTP/REST adapters.
type RestRepository interface {
	// Create creates a new entity using a REST adapter.
	Create(ctx context.Context, request RestRequest, response *RestResponse) error
	// List returns a collection of entities.
	List(ctx context.Context, reqCtx types.RequestContext, response *RestResponses) error
	// Get returns one entity identified in reqCtx.
	Get(ctx context.Context, reqCtx types.RequestContext, response *RestResponse) error
	// Update replaces entity state.
	Update(ctx context.Context, request RestRequest, response *RestResponse) error
	// Save applies partial changes or upsert semantics.
	Save(ctx context.Context, request RestRequest, response *RestResponse) error
	// Delete removes one entity identified in reqCtx.
	Delete(ctx context.Context, reqCtx types.RequestContext, response *RestResponse) error
}
