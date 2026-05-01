package net_http

import (
	"net/http"
	"testing"
)

func TestNewMiddlewareRegistry_ReturnsNonNil(t *testing.T) {
	if got := NewMiddlewareRegistry(); got == nil {
		t.Fatalf("expected non-nil registry")
	}
}

func TestMiddlewareRegistry_RegisterPanicsOnNilRegistry(t *testing.T) {
	var registry *MiddlewareRegistry
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil registry")
		}
	}()
	registry.Register(func(next http.Handler) http.Handler { return next })
}

func TestMiddlewareRegistry_RegisterPanicsOnNilMiddleware(t *testing.T) {
	registry := NewMiddlewareRegistry()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil middleware")
		}
	}()
	registry.Register(nil)
}

func TestMiddlewareRegistry_MiddlewaresReturnsCopy(t *testing.T) {
	registry := NewMiddlewareRegistry()
	registry.Register(func(next http.Handler) http.Handler { return next })

	mws := registry.Middlewares()
	if len(mws) != 1 {
		t.Fatalf("expected 1 middleware, got %d", len(mws))
	}
	mws[0] = nil

	got := registry.Middlewares()
	if len(got) != 1 || got[0] == nil {
		t.Fatalf("expected defensive copy from Middlewares")
	}
}

func TestMiddlewareRegistry_ResetClearsMiddlewares(t *testing.T) {
	registry := NewMiddlewareRegistry()
	registry.Register(func(next http.Handler) http.Handler { return next })
	registry.Reset()

	if got := len(registry.Middlewares()); got != 0 {
		t.Fatalf("expected empty registry after reset, got %d", got)
	}
}

func TestMiddlewareRegistry_MiddlewaresNilRegistry(t *testing.T) {
	var registry *MiddlewareRegistry
	if got := registry.Middlewares(); got != nil {
		t.Fatalf("expected nil from nil registry Middlewares, got %v", got)
	}
}

func TestMiddlewareRegistry_ResetNilRegistry(t *testing.T) {
	var registry *MiddlewareRegistry
	registry.Reset() // must not panic
}

func TestMiddlewareRegistry_RegisterSkipsEmptyList(t *testing.T) {
	registry := NewMiddlewareRegistry()
	registry.Register() // no-op, must not panic
	if got := len(registry.Middlewares()); got != 0 {
		t.Fatalf("expected 0 middlewares after empty register, got %d", got)
	}
}
