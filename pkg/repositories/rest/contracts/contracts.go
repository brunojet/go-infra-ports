// Package contracts defines public REST repository contracts.
package contracts

import (
	"context"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// RestCreate wraps the request body with REST-specific context.
type RestCreate[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the create payload.
	Body C
}

// RestUpdate wraps update payload and REST context.
type RestUpdate[U any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the update payload.
	Body U
}

// RestSave wraps save payload and REST context.
type RestSave[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the save payload.
	Body C
}

// RestResponse represents one REST repository output.
type RestResponse[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Data carries the payload.
	Data R
}

// RestResponses represents many REST repository outputs.
type RestResponses[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Data carries the payload collection.
	Data []R
}

// RestRepository is the outbound port for HTTP/REST adapters.
type RestRepository[C, R, U any] interface {
	// Create creates a new entity using a REST adapter.
	Create(ctx context.Context, request RestCreate[C], response *RestResponse[R]) error
	// List returns a collection of entities.
	List(ctx context.Context, reqCtx types.RequestContext, response *RestResponses[R]) error
	// Get returns one entity identified in reqCtx.
	Get(ctx context.Context, reqCtx types.RequestContext, response *RestResponse[R]) error
	// Update replaces entity state.
	Update(ctx context.Context, request RestUpdate[U], response *RestResponse[R]) error
	// Save applies partial changes or upsert semantics.
	Save(ctx context.Context, request RestSave[C], response *RestResponse[R]) error
	// Delete removes one entity identified in reqCtx.
	Delete(ctx context.Context, reqCtx types.RequestContext, response *RestResponse[R]) error
}
