// Package handlers exposes the public handler API and re-exports for consumers.
package handlers

import (
	"net/http"
	"regexp"

	"github.com/brunojet/go-infra-ports/internal/handlers/net_http"
)

// RouteMethod is a bitmask that selects which HTTP routes a handler registers.
type RouteMethod = net_http.RouteMethod

const (
	// MethodCreate registers the collection POST route.
	MethodCreate RouteMethod = net_http.MethodCreate
	// MethodList registers the collection GET route.
	MethodList RouteMethod = net_http.MethodList
	// MethodGet registers the instance GET route.
	MethodGet RouteMethod = net_http.MethodGet
	// MethodUpdate registers the instance PUT route.
	MethodUpdate RouteMethod = net_http.MethodUpdate
	// MethodSave registers the instance PATCH route.
	MethodSave RouteMethod = net_http.MethodSave
	// MethodDelete registers the instance DELETE route.
	MethodDelete RouteMethod = net_http.MethodDelete
	// AllCollectionMethods combines MethodCreate and MethodList.
	AllCollectionMethods RouteMethod = net_http.AllCollectionMethods
	// AllInstanceMethods combines MethodGet, MethodUpdate, MethodSave, and MethodDelete.
	AllInstanceMethods RouteMethod = net_http.AllInstanceMethods
)

// Middleware wraps HTTP handlers.
type Middleware = net_http.Middleware

// MiddlewareRegistry stores a middleware chain for handler registration.
type MiddlewareRegistry = net_http.MiddlewareRegistry

// HandlerOption configures the net/http handler adapter.
type HandlerOption = net_http.HandlerOption

// NewNetHttpHandler builds and registers a net/http CRUD handler on mux.
func NewNetHttpHandler[C, R, U any](
	mux *http.ServeMux,
	svc net_http.Service[C, R, U],
	options ...HandlerOption,
) net_http.NetHttpHandler {
	return net_http.NewNetHttpHandler(mux, svc, options...)
}

// NewMiddlewareRegistry creates an empty middleware registry.
func NewMiddlewareRegistry() *MiddlewareRegistry {
	return net_http.NewMiddlewareRegistry()
}

// WithCollection configures collection routes and path parameter validation.
func WithCollection(methods RouteMethod, pathFmt string, parentsValidation ...*regexp.Regexp) HandlerOption {
	return net_http.WithCollection(methods, pathFmt, parentsValidation...)
}

// WithInstance configures instance routes and path parameter validation.
func WithInstance(methods RouteMethod, pathFmt string, parentsValidation ...*regexp.Regexp) HandlerOption {
	return net_http.WithInstance(methods, pathFmt, parentsValidation...)
}

// WithMiddlewareRegistry applies middlewares registered in registry to all routes.
func WithMiddlewareRegistry(registry *MiddlewareRegistry) HandlerOption {
	return net_http.WithMiddlewareRegistry(registry)
}

// WithInformation overrides the default builder for 1xx responses.
func WithInformation(fn func(status int, meta net_http.ServiceMeta) any) HandlerOption {
	return net_http.WithInformation(fn)
}

// WithEnvelope overrides the default builder for 2xx response payloads.
func WithEnvelope(fn func(status int, meta net_http.ServiceMeta, resp any) any) HandlerOption {
	return net_http.WithEnvelope(fn)
}

// WithRedirection overrides the default builder for 3xx responses.
func WithRedirection(fn func(status int, meta net_http.ServiceMeta) (body any, location string)) HandlerOption {
	return net_http.WithRedirection(fn)
}

// WithProblemDetails overrides the default builder for 4xx/5xx responses.
func WithProblemDetails(fn func(status int, instance string, meta net_http.ServiceMeta) any) HandlerOption {
	return net_http.WithProblemDetails(fn)
}
