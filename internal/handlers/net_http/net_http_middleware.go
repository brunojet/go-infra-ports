package net_http

import "sync"

// MiddlewareRegistry centralizes reusable middlewares for net/http handlers.
// It is safe for concurrent registration and reads during application bootstrap.
type MiddlewareRegistry struct {
	mu          sync.RWMutex
	middlewares []Middleware
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{}
}

// Register appends middlewares in the provided order.
func (r *MiddlewareRegistry) Register(mws ...Middleware) {
	if r == nil {
		panic("net_http: MiddlewareRegistry.Register: registry must not be nil")
	}
	validated := cloneMiddlewares(mws, "MiddlewareRegistry.Register")
	if len(validated) == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, validated...)
}

// Middlewares returns a copy of the registered middleware chain.
func (r *MiddlewareRegistry) Middlewares() []Middleware {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Middleware(nil), r.middlewares...)
}

// Reset clears the registry. This exists mainly for tests and controlled bootstrap reconfiguration.
func (r *MiddlewareRegistry) Reset() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = nil
}

func cloneMiddlewares(mws []Middleware, caller string) []Middleware {
	if len(mws) == 0 {
		return nil
	}
	cloned := make([]Middleware, len(mws))
	for i, mw := range mws {
		if mw == nil {
			panic("net_http: " + caller + ": middleware must not be nil")
		}
		cloned[i] = mw
	}
	return cloned
}
