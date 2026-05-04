// Package rest exposes the public REST service API, including the service constructor,
// functional options, and default mapper implementations.
package rest

import (
	"github.com/brunojet/go-infra-ports/internal/services/rest"
	rpocts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	svcrestcts "github.com/brunojet/go-infra-ports/pkg/services/rest/contracts"
)

type (
	//nolint:revive // Kept for API clarity and compatibility with constructor signatures.
	// RestServiceOption is a marker interface for functional options accepted by NewRestService.
	RestServiceOption = rest.RestServiceOption

	// DefaultRestUpstreamMapper is the default no-op upstream mapper implementation.
	DefaultRestUpstreamMapper[C, U any] = rest.DefaultRestUpstreamMapper[C, U]

	// DefaultRestDownstreamMapper is the default downstream mapper implementation.
	DefaultRestDownstreamMapper[R any] = rest.DefaultRestDownstreamMapper[R]
)

// WithUpstreamMapper returns a RestServiceOption that sets a custom upstream mapper.
// Panics if mapper is nil.
func WithUpstreamMapper[C, U any](mapper svcrestcts.RestUpstreamMapper[C, U]) RestServiceOption {
	return rest.WithUpstreamMapper[C, U](mapper)
}

// WithDownstreamMapper returns a RestServiceOption that sets a custom downstream mapper.
// Panics if mapper is nil.
func WithDownstreamMapper[R any](mapper svcrestcts.RestDownstreamMapper[R]) RestServiceOption {
	return rest.WithDownstreamMapper[R](mapper)
}

// NewRestService creates a REST-backed service with injectible mappers.
// Pass WithUpstreamMapper or WithDownstreamMapper to override the defaults.
// Returns an error if repo is nil.
func NewRestService[C, R, U any](
	repo rpocts.RestRepository,
	opts ...RestServiceOption,
) (svccts.Service[C, R, U], error) {
	return rest.NewRestService[C, R, U](repo, opts...)
}
