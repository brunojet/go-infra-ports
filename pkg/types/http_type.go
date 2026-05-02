package types

import (
	"net/http"
	"net/url"
)

// HTTPRequesOptions represents HTTP request metadata used by adapters.
type HTTPRequesOptions struct {
	// Header carries outbound HTTP headers.
	Header http.Header
	// Query carries outbound URL query parameters.
	Query url.Values
}

// HTTPResponseOptions represents HTTP response metadata used by adapters.
type HTTPResponseOptions struct {
	// StatusCode is the HTTP status to be returned.
	StatusCode int
	// Header carries inbound or outbound HTTP headers.
	Header http.Header
}
