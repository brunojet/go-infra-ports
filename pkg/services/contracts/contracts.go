// Package contracts defines public service-layer contracts.
package contracts

import (
	"context"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// ServiceCreate carries input for create operations.
type ServiceCreate[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the create payload.
	Body C
}

// ServiceUpdate carries input for update operations.
type ServiceUpdate[U any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the update payload.
	Body U
}

// ServiceSave carries input for full-replacement save operations (PUT semantics).
type ServiceSave[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the save payload.
	Body C
}

// ServiceMeta contains transport-neutral metadata returned by a service.
type ServiceMeta struct {
	// Message is a human-readable summary.
	Message string
	// Location is the canonical URI used by redirection/create flows.
	Location string
	// Code is a domain-specific code, not an HTTP status.
	Code string
	// Metadata stores extension key/value pairs.
	Metadata map[string]any
}

// ServiceResponse represents a single-entity service output.
type ServiceResponse[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Meta carries status-independent metadata.
	Meta ServiceMeta
	// Data carries the payload.
	Data R
}

// ServiceResponses represents a multi-entity service output.
type ServiceResponses[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Meta carries status-independent metadata.
	Meta ServiceMeta
	// Data carries the payload collection.
	Data []R
}

// Service defines CRUD operations expected from application services.
type Service[C, R, U any] interface {
	// Create creates a new entity.
	Create(ctx context.Context, request ServiceCreate[C], response *ServiceResponse[R]) error
	// List returns a collection of entities.
	List(ctx context.Context, reqCtx types.RequestContext, response *ServiceResponses[R]) error
	// Get returns one entity identified in reqCtx.
	Get(ctx context.Context, reqCtx types.RequestContext, response *ServiceResponse[R]) error
	// Update replaces entity state.
	Update(ctx context.Context, request ServiceUpdate[U], response *ServiceResponse[R]) error
	// Save replaces the full entity state (PUT semantics).
	Save(ctx context.Context, request ServiceSave[C], response *ServiceResponse[R]) error
	// Delete removes one entity identified in reqCtx.
	Delete(ctx context.Context, reqCtx types.RequestContext, response *ServiceResponse[R]) error
}
