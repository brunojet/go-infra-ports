// Package contracts defines public REST repository contracts.
package contracts

import (
	"context"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// RestRequestV2 wraps the request body with REST-specific context.
type RestRequestV2 struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Data carries the create payload.
	Data any
}

// Release releases any resources held by the request, such as underlying data specs. After calling Release, the RestRequest should not be used unless re-populated by a repository method. This is a no-op for the current implementation but is included for future-proofing and to signal intent around request lifecycle management. Consumers should call Release when they are done with the RestRequest to allow for proper cleanup and resource management in future implementations.
// In the future this may be expanded to include other cleanup operations as needed (e.g. closing request body readers if the contract evolves to include them). This is a no-op for the current implementation but is included for future-proofing and to signal intent around request lifecycle management. Consumers should call Release when they are done with the RestRequest to allow for proper cleanup and resource management in future implementations.
func (r *RestRequestV2) Release() {
	r.Data = nil
}

// RestResponseV2 represents one REST repository output.
type RestResponseV2 struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information any
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection any
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem any
	// Data carries the decoded successful payload.
	Data any
}

// Release releases any resources held by the response, such as underlying data specs.
// In the future this may be expanded to include other cleanup operations as needed (e.g. closing response body readers if the contract evolves to include them). After calling Release, the RestResponse should not be used unless re-populated by a repository method. This is a no-op for the current implementation but is included for future-proofing and to signal intent around response lifecycle management. Consumers should call Release when they are done with the RestResponse to allow for proper cleanup and resource management in future implementations.
func (r *RestResponseV2) Release() {
	r.Data = nil
	r.Information = nil
	r.Redirection = nil
	r.Problem = nil
}

// RestResponsesV2 represents many REST repository outputs.
type RestResponsesV2 struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Information carries decoded informational payloads for 1xx responses.
	Information any
	// Redirection carries decoded redirection payloads for 3xx responses.
	Redirection any
	// Problem carries decoded problem payloads for 4xx/5xx responses.
	Problem any
	// Data carries the decoded successful payload collection.
	Data []any
}

// RestRepositoryV2 is the outbound port for HTTP/REST adapters.
type RestRepositoryV2 interface {
	// Create creates a new entity using a REST adapter.
	Create(ctx context.Context, request *RestRequestV2, response *RestResponseV2) error
	// List returns a collection of entities.
	List(ctx context.Context, reqCtx types.RequestContext, response *RestResponsesV2) error
	// Get returns one entity identified in reqCtx.
	Get(ctx context.Context, reqCtx types.RequestContext, response *RestResponseV2) error
	// Update replaces entity state.
	Update(ctx context.Context, request *RestRequestV2, response *RestResponseV2) error
	// Save applies partial changes or upsert semantics.
	Save(ctx context.Context, request *RestRequestV2, response *RestResponseV2) error
	// Delete removes one entity identified in reqCtx.
	Delete(ctx context.Context, reqCtx types.RequestContext, response *RestResponseV2) error
	// NewRequest allocates a new RestRequest for the given method, pre-populating any method-specific fields (e.g. Data with the expected spec type). This is used by the service layer to obtain properly initialized request instances for mapping and repository calls.
	NewRequest(method RestMethod) (*RestRequestV2, error)
	// NewResponse allocates a new RestResponse with properly initialized data specs. This is used by the service layer to obtain properly initialized response instances for repository calls and mapping. The implementation should ensure that the Data, Information, Redirection, and Problem fields are initialized with the appropriate concrete types expected by the repository and mappers.
	NewResponse() *RestResponseV2
}
