package rest

import (
	"encoding/json"
	"net/http"
	"testing"
)

// --- cloneConfig ---

func TestCloneConfig_CreatesIndependentMaps(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	// use non-default keys to verify independence
	keyX := RestMethod(99)
	keyY := RestMethod(100)
	registry.cfg.Requests[keyX] = &DefaultRestRequest{}

	cloned := registry.cloneConfig()
	if cloned.Requests[keyX] == nil {
		t.Fatal("expected cloned config to copy requests map")
	}

	cloned.Requests[keyY] = &DefaultRestRequest{}
	if registry.cfg.Requests[keyY] != nil {
		t.Fatal("expected cloned config map to be independent from source")
	}
}

// --- resolveResponseSpec (method) ---

func TestResolveResponseSpec_ByStatusClass(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.Informations[http.StatusContinue] = &testResponse{}
	registry.cfg.Redirections[http.StatusMovedPermanently] = &testResponse{}
	registry.cfg.Problems[http.StatusBadRequest] = &testResponse{}
	registry.cfg.Responses[http.StatusOK] = &testResponse{}

	cases := []struct {
		status int
		label  string
	}{
		{http.StatusContinue, "information"},
		{http.StatusMovedPermanently, "redirection"},
		{http.StatusBadRequest, "problem"},
		{http.StatusOK, "response"},
	}
	for _, tc := range cases {
		spec, err := registry.resolveResponseSpec(tc.status)
		if err != nil {
			t.Fatalf("unexpected error for %s (status %d): %v", tc.label, tc.status, err)
		}
		if _, ok := spec.(*testResponse); !ok {
			t.Fatalf("expected testResponse for %s (status %d), got %T", tc.label, tc.status, spec)
		}
	}
}

func TestResolveResponseSpec_ReturnsErrorWhenNil(t *testing.T) {
	registry := &restRegistry{cfg: &registryOptions{
		Responses: map[int]RestResponseSpec{},
	}}
	_, err := registry.resolveResponseSpec(http.StatusContinue + 600)
	if err == nil {
		t.Fatal("expected error when spec is nil, got nil")
	}
}

// --- resolveRequestSpec ---

func TestResolveRequestSpec_Hit(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.Requests[MethodCreate] = &testLocalRequest{}

	resolved, err := registry.newRequestSpec(MethodCreate)
	if err != nil {
		t.Fatalf("unexpected error for direct hit: %v", err)
	}
	if _, ok := resolved.(*testLocalRequest); !ok {
		t.Fatalf("expected testLocalRequest for direct hit, got %T", resolved)
	}
}

func TestResolveRequestSpec_FallbackToDefault(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}

	resolved, err := registry.newRequestSpec(MethodCreate)
	if err != nil {
		t.Fatalf("unexpected error for fallback: %v", err)
	}
	if _, ok := resolved.(*DefaultRestRequest); !ok {
		t.Fatalf("expected DefaultRestRequest fallback, got %T", resolved)
	}
}

// --- resolveRequestEnvelopeSpec ---

func TestResolveRequestEnvelopeSpec_Hit(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.RequestsEnvelopes[MethodCreate] = &testLocalRequest{}

	resolved := registry.resolveRequestEnvelopeSpec(MethodCreate)
	if _, ok := resolved.(*testLocalRequest); !ok {
		t.Fatalf("expected testLocalRequest for direct hit, got %T", resolved)
	}
}

func TestResolveRequestEnvelopeSpec_FallbackToDefault(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.RequestsEnvelopes[MethodCreate] = &testLocalRequest{}

	resolved := registry.resolveRequestEnvelopeSpec(MethodCreate)
	if _, ok := resolved.(*testLocalRequest); !ok {
		t.Fatalf("expected testLocalRequest via MethodCreate fallback, got %T", resolved)
	}
}

func TestResolveRequestEnvelopeSpec_ReturnsNilWhenUnconfigured(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}

	resolved := registry.resolveRequestEnvelopeSpec(MethodCreate)
	if resolved != nil {
		t.Fatalf("expected nil when unconfigured, got %T", resolved)
	}
}

// --- resolveResponseEnvelopeSpec ---

func TestResolveResponseEnvelopeSpec_Hit(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.ResponseEnvelopes[http.StatusOK] = &responseEnvelope{}

	resolved := registry.resolveResponseEnvelopeSpec(http.StatusOK)
	if resolved == nil {
		t.Fatal("expected non-nil envelope spec for direct hit")
	}
}

func TestResolveResponseEnvelopeSpec_FallbackToDefault(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}
	registry.cfg.ResponseEnvelopes[DefaultStatusCode] = &responseEnvelope{}

	resolved := registry.resolveResponseEnvelopeSpec(http.StatusOK)
	if resolved == nil {
		t.Fatal("expected non-nil envelope spec via DefaultStatusCode fallback")
	}
}

func TestResolveResponseEnvelopeSpec_ReturnsNilWhenUnconfigured(t *testing.T) {
	registry := &restRegistry{cfg: newRegistryOptions()}

	resolved := registry.resolveResponseEnvelopeSpec(http.StatusOK)
	if resolved != nil {
		t.Fatalf("expected nil when unconfigured, got %T", resolved)
	}
}

// --- test helpers ---

type testLocalRequest struct {
	Body json.RawMessage
}

func (r *testLocalRequest) New() RestRequestSpec { return &testLocalRequest{} }
func (r *testLocalRequest) SetBody(b RestRequestSpec) {
	if raw, ok := b.(*DefaultRestRequest); ok {
		r.Body = append([]byte(nil), raw.Body...)
	}
}
