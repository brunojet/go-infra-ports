// Package contracts defines public repository-layer contracts.
package contracts

import (
	"context"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// RepositoryCreate carries input for create operations.
type RepositoryCreate[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the create payload.
	Body C
}

// RepositoryUpdate carries input for update operations.
type RepositoryUpdate[U any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the update payload.
	Body U
}

// RepositorySave carries input for save operations.
type RepositorySave[C any] struct {
	// Context carries transport metadata.
	Context types.RequestContext
	// Body carries the save payload.
	Body C
}

// RepositoryResponse represents a single-entity repository output.
type RepositoryResponse[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Data carries the payload.
	Data R
}

// RepositoryResponses represents a multi-entity repository output.
type RepositoryResponses[R any] struct {
	// Context carries response transport metadata.
	Context types.ResponseContext
	// Data carries the payload collection.
	Data []R
}

// Repository is a transport/storage-agnostic outbound port.
// Adapters (REST, GORM, Redis, etc.) implement this interface.
type Repository[C, R, U any] interface {
	// Create creates a new entity in the outbound system.
	Create(ctx context.Context, request RepositoryCreate[C], response *RepositoryResponse[R]) error
	// List returns a collection of entities.
	List(ctx context.Context, reqCtx types.RequestContext, response *RepositoryResponses[R]) error
	// Get returns one entity identified in reqCtx.
	Get(ctx context.Context, reqCtx types.RequestContext, response *RepositoryResponse[R]) error
	// Update replaces entity state.
	Update(ctx context.Context, request RepositoryUpdate[U], response *RepositoryResponse[R]) error
	// Save applies partial changes or upsert semantics.
	Save(ctx context.Context, request RepositorySave[C], response *RepositoryResponse[R]) error
	// Delete removes one entity identified in reqCtx.
	Delete(ctx context.Context, reqCtx types.RequestContext, response *RepositoryResponse[R]) error
}
