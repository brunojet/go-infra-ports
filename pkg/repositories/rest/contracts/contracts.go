// Package contracts defines public REST repository contracts.
package contracts

import (
	"context"
	"encoding/json"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// RestMethod selects which HTTP verbs are registered for a route entry.
// It is used as a bitmask combining method constants (e.g. MethodCreate,
// MethodList, MethodGet, MethodUpdate, MethodSave, MethodDelete).
type RestMethod uint8

// Method* constants identify the REST operations used by registries
// and repositories (e.g. create, list, get, update, save, delete).
const (
	MethodCreate RestMethod = 1 << iota // POST  — collection
	MethodList                          // GET   — collection
	MethodGet                           // GET   — instance
	MethodUpdate                        // PUT   — instance
	MethodSave                          // PATCH — instance
	MethodDelete                        // DELETE — instance
)

// RestDataSpec is a unified V2 data spec used for both requests and responses.
// It provides construction, slice allocation, body assignment and JSON
// marshal/unmarshal hooks to support generic data-backed specs.
//
// NOTE (design): Consider changing `SetBody` and `Body` to return an `error`
// in a future major version to allow the spec implementation to validate
// the provided value and surface clear failures when values do not match the
// expected concrete type (for example wrong type or `nil` when not allowed).
// Example signatures:
//
//	`SetBody(body any) error`
//	`Body() (any, error)`
//
// This improves safety and explicit error handling but is a breaking change
// and should be accompanied by a migration/adapter strategy.
type RestDataSpec interface {
	// New returns a new empty data spec instance of the same concrete type.
	New() RestDataSpec
	// NewSlice returns a typed data spec slice with the requested size.
	NewSlice(n int) []RestDataSpec
	// SetBody assigns the request/response body from another value.
	SetBody(body any) error
	// Body returns the underlying request/response body value.
	Body() any
	// MarshalJSON marshals the underlying data for outbound requests.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals inbound response bytes into the underlying data.
	UnmarshalJSON(data []byte) error
}

// RestEnvelopeSpec defines envelope payload contracts for the V2 registry API.
// It mirrors RestEnvelopeSpec but is used by V2 envelope implementations.
type RestEnvelopeSpec interface {
	// New returns a new empty envelope spec instance of the same concrete type.
	New() RestEnvelopeSpec
	// EnvelopeData returns the inner data payload extracted from the envelope.
	EnvelopeData() json.RawMessage
	// EnvelopeMeta returns the metadata extracted from the envelope.
	EnvelopeMeta() types.ResponseMeta
	// SetEnvelopeMeta assigns the metadata extracted from the envelope.
	SetEnvelopeData(data json.RawMessage)
	// SetEnvelopeMeta assigns the metadata extracted from the envelope.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals inbound response bytes into the underlying envelope.
	UnmarshalJSON(data []byte) error
}

// RestRegistry is the V2 registry contract used by the newer REST
// repository/registry implementations. It exposes mapping and codec
// operations for request/response marshaling, envelope handling, and
// allocation of request spec instances used by the service layer.
//
// The interface mirrors the older `RestRegistry` API but operates on the
// unified `RestDataSpec`/`RestEnvelopeSpecV2` abstractions introduced in V2.
type RestRegistry interface {
	// Merge combines the current registry with another and returns the merged result.
	Merge(other RestRegistry) RestRegistry
	// ResolveRequest marshals a request spec into requestBody.
	ResolveRequest(body RestDataSpec, requestBody *[]byte) error
	// ResolveEnvelopeRequest envelopes the current request body in place.
	ResolveEnvelopeRequest(method RestMethod, dataBody *[]byte) error
	// ResolveResponse unmarshals a single response payload into body.
	ResolveResponse(status int, responseBody []byte, body *RestDataSpec) error
	// ResolveResponses unmarshals a collection response payload into bodies.
	ResolveResponses(status int, responseBody []byte, bodies *[]RestDataSpec) error
	// ResolveEnvelopeResponse extracts payload and metadata from an envelope.
	ResolveEnvelopeResponse(status int, dataBody *[]byte, meta *types.ResponseMeta) error
	// NewRequestSpec allocates a request spec instance for service layer usage.
	// Returns an error if the spec's New method returns nil.
	NewRequestSpec(method RestMethod) (RestDataSpec, error)
	// ReleaseRequestSpec releases a request spec instance used by the service layer.
	ReleaseRequestSpec(spec RestDataSpec)
}

// RestRequest wraps the request body with REST-specific context.
type RestRequest struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Data carries the create payload.
	Data RestDataSpec
}

// Release releases any resources held by the request, such as underlying data specs. After calling Release, the RestRequest should not be used unless re-populated by a repository method. This is a no-op for the current implementation but is included for future-proofing and to signal intent around request lifecycle management. Consumers should call Release when they are done with the RestRequest to allow for proper cleanup and resource management in future implementations.
// In the future this may be expanded to include other cleanup operations as needed (e.g. closing request body readers if the contract evolves to include them). This is a no-op for the current implementation but is included for future-proofing and to signal intent around request lifecycle management. Consumers should call Release when they are done with the RestRequest to allow for proper cleanup and resource management in future implementations.
func (r *RestRequest) Release() {
	r.Data = nil
}

// RestResponse represents one REST repository output.
type RestResponse struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information RestDataSpec
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection RestDataSpec
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem RestDataSpec
	// Data carries the decoded successful payload.
	Data RestDataSpec
}

// Release releases any resources held by the response, such as underlying data specs.
// In the future this may be expanded to include other cleanup operations as needed (e.g. closing response body readers if the contract evolves to include them). After calling Release, the RestResponse should not be used unless re-populated by a repository method. This is a no-op for the current implementation but is included for future-proofing and to signal intent around response lifecycle management. Consumers should call Release when they are done with the RestResponse to allow for proper cleanup and resource management in future implementations.
func (r *RestResponse) Release() {
	r.Data = nil
	r.Information = nil
	r.Redirection = nil
	r.Problem = nil
}

// RestResponses represents many REST repository outputs.
type RestResponses struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information RestDataSpec
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection RestDataSpec
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem RestDataSpec
	// Data carries the decoded successful payload collection.
	Data []RestDataSpec
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
	// NewRequest allocates a new RestRequest for the given method, pre-populating any method-specific fields (e.g. Data with the expected spec type). This is used by the service layer to obtain properly initialized request instances for mapping and repository calls.
	NewRequest(method RestMethod) (*RestRequest, error)
	// NewResponse allocates a new RestResponse with properly initialized data specs. This is used by the service layer to obtain properly initialized response instances for repository calls and mapping. The implementation should ensure that the Data, Information, Redirection, and Problem fields are initialized with the appropriate concrete types expected by the repository and mappers.
	NewResponse() *RestResponse
}
