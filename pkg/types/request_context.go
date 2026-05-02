// Package types defines shared transport-level types.
package types

import (
	"net/http"
	"net/url"
)

// RequestContext carries transport metadata that business rules or adapters may use.
// Observability concerns (correlation ID, tenant, user) are handled by middlewares via context.Context.
type RequestContext struct {
	Query       url.Values
	Headers     http.Header
	Identifiers map[string]string // path params extracted and validated by the handler
}

// ResponseContext carries transport metadata populated by adapters after execution.
type ResponseContext struct {
	// StatusCode is the HTTP status to emit when provided by a lower layer.
	StatusCode int
	// Headers carries additional response headers.
	Headers http.Header
}
