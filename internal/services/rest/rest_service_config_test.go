package rest

import (
	"testing"
)

func TestWithUpstreamMapper_SetsMapper(t *testing.T) {
	mapper := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	opt := WithUpstreamMapper[testCreatePayload, testUpdatePayload](mapper)

	o := newRestServiceOptions[testCreatePayload, testResponse, testUpdatePayload]([]RestServiceOption{opt})

	if o.upstream != mapper {
		t.Fatal("expected upstream to be the provided mapper")
	}
}

func TestWithUpstreamMapper_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream mapper")
		}
	}()
	WithUpstreamMapper[testCreatePayload, testUpdatePayload](nil)
}

func TestWithDownstreamMapper_SetsMapper(t *testing.T) {
	mapper := &DefaultRestDownstreamMapper[testResponse]{}
	opt := WithDownstreamMapper[testResponse](mapper)

	o := newRestServiceOptions[testCreatePayload, testResponse, testUpdatePayload]([]RestServiceOption{opt})

	if o.downstream != mapper {
		t.Fatal("expected downstream to be the provided mapper")
	}
}

func TestWithDownstreamMapper_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil downstream mapper")
		}
	}()
	WithDownstreamMapper[testResponse](nil)
}

func TestNewRestServiceOptions_Defaults(t *testing.T) {
	o := newRestServiceOptions[testCreatePayload, testResponse, testUpdatePayload](nil)

	if o.upstream == nil {
		t.Fatal("expected default upstream mapper to be set")
	}
	if o.downstream == nil {
		t.Fatal("expected default downstream mapper to be set")
	}
}

func TestNewRestServiceOptions_WithBothOptions(t *testing.T) {
	customUpstream := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	customDownstream := &DefaultRestDownstreamMapper[testResponse]{}

	o := newRestServiceOptions[testCreatePayload, testResponse, testUpdatePayload]([]RestServiceOption{
		WithUpstreamMapper[testCreatePayload, testUpdatePayload](customUpstream),
		WithDownstreamMapper[testResponse](customDownstream),
	})

	if o.upstream != customUpstream {
		t.Fatal("expected custom upstream mapper")
	}
	if o.downstream != customDownstream {
		t.Fatal("expected custom downstream mapper")
	}
}

func TestNewRestServiceOptions_UnknownOptionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown option type")
		}
	}()
	newRestServiceOptions[testCreatePayload, *testResponse, testUpdatePayload]([]RestServiceOption{
		unknownOption{},
	})
}

func TestNewRestService_WithDownstreamOption(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	custom := &DefaultDownstreamCapture{}

	svc, err := NewRestService[testCreatePayload, testResponse, testUpdatePayload](
		repo,
		WithDownstreamMapper[testResponse](custom),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

// unknownOption is an unrecognized RestServiceOption used to test the runtime panic.
type unknownOption struct{}

// DefaultDownstreamCapture embeds DefaultRestDownstreamMapper to satisfy the interface.
type DefaultDownstreamCapture struct {
	DefaultRestDownstreamMapper[testResponse]
}
