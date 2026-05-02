package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

// --- mergeRegistryOptions ---

func TestMergeRegistryOptions_CopiesAllMaps(t *testing.T) {
	dst := newRegistryOptions()
	src := newRegistryOptions()
	src.Requests["m"] = &DefaultRestRequest{}
	src.Responses[201] = &DefaultRestResponse{}
	src.ResponseEnvelopes[201] = &responseEnvelope{}
	src.Informations[101] = &DefaultRestResponse{}
	src.Redirections[301] = &DefaultRestResponse{}
	src.Problems[400] = &DefaultRestResponse{}

	mergeRegistryOptions(dst, src)

	if dst.Requests["m"] == nil || dst.Responses[201] == nil || dst.ResponseEnvelopes[201] == nil {
		t.Fatal("expected merged config to include request/response maps")
	}
	if dst.Informations[101] == nil || dst.Redirections[301] == nil || dst.Problems[400] == nil {
		t.Fatal("expected merged config to include status class maps")
	}
}

// --- resolveResponseInstance ---

func TestResolveResponseInstance_Success(t *testing.T) {
	instance, err := resolveResponseInstance(&DefaultRestResponse{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instance == nil {
		t.Fatal("expected non-nil instance")
	}
}

func TestResolveResponseInstance_ReturnsErrorWhenNewReturnsNil(t *testing.T) {
	_, err := resolveResponseInstance(&nilNewSpec{})
	if err == nil {
		t.Fatal("expected error when New() returns nil")
	}
}

// --- resolveResponseSlice ---

func TestResolveResponseSlice_ReturnsCorrectLength(t *testing.T) {
	resolved, err := resolveResponseSlice(&DefaultRestResponse{}, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resolved) != 3 {
		t.Fatalf("unexpected slice length: got %d, want 3", len(resolved))
	}
}

func TestResolveResponseSlice_ReturnsErrorOnLengthMismatch(t *testing.T) {
	_, err := resolveResponseSlice(&badSliceSpec{}, 3)
	if err == nil {
		t.Fatal("expected error on length mismatch, got nil")
	}
}

// --- resolveResponseSpec (free function) ---

func TestResolveResponseSpec_HitByStatus(t *testing.T) {
	source := map[int]RestResponseSpec{
		DefaultStatusCode:  &DefaultRestResponse{},
		http.StatusCreated: &testResponse{},
	}
	resolved := resolveResponseSpec(http.StatusCreated, source)
	if _, ok := resolved.(*testResponse); !ok {
		t.Fatalf("expected testResponse for direct hit, got %T", resolved)
	}
}

func TestResolveResponseSpec_FallbackToDefault(t *testing.T) {
	source := map[int]RestResponseSpec{
		DefaultStatusCode: &DefaultRestResponse{},
	}
	resolved := resolveResponseSpec(http.StatusOK, source)
	if _, ok := resolved.(*DefaultRestResponse); !ok {
		t.Fatalf("expected DefaultRestResponse fallback, got %T", resolved)
	}
}

// --- marshalInto ---

func TestMarshalInto_Success(t *testing.T) {
	var out []byte
	if err := marshalInto(&out, map[string]string{"k": "v"}, func(e error) error { return e }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !json.Valid(out) {
		t.Fatal("expected valid JSON output")
	}
}

func TestMarshalInto_Error(t *testing.T) {
	sentinel := errors.New("wrap")
	var out []byte
	err := marshalInto(&out, make(chan int), func(e error) error { return sentinel })
	if err != sentinel {
		t.Fatalf("expected wrapped sentinel error, got %v", err)
	}
}

// --- unmarshalInto ---

func TestUnmarshalInto_Success(t *testing.T) {
	var dst map[string]string
	if err := unmarshalInto([]byte(`{"k":"v"}`), &dst, func(e error) error { return e }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["k"] != "v" {
		t.Fatalf("expected 'v', got %q", dst["k"])
	}
}

func TestUnmarshalInto_Error(t *testing.T) {
	sentinel := errors.New("wrap")
	var dst map[string]string
	err := unmarshalInto([]byte(`not-json`), &dst, func(e error) error { return sentinel })
	if err != sentinel {
		t.Fatalf("expected wrapped sentinel error, got %v", err)
	}
}

// --- test helpers ---

type nilNewSpec struct{}

func (n *nilNewSpec) New() RestResponseSpec             { return nil }
func (n *nilNewSpec) NewSlice(_ int) []RestResponseSpec { return nil }

type badSliceSpec struct{}

func (b *badSliceSpec) New() RestResponseSpec             { return &badSliceSpec{} }
func (b *badSliceSpec) NewSlice(_ int) []RestResponseSpec { return []RestResponseSpec{} }
