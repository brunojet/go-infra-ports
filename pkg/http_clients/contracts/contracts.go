// Package contracts defines the public HttpClient outbound port contract.
package contracts

import (
	"context"
	"net/http"
)

// HttpClient is the outbound port for executing HTTP requests.
// Any implementation whose Do method matches this signature satisfies the
// interface, enabling adapter modules (e.g. go-infra-adapters) to be wired
// at the application bootstrap layer without a compile-time dependency between
// this module and the adapter module.
type HttpClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
