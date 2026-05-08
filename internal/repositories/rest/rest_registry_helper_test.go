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
	reqDefault := NewDataSpecOf[DefaultRestRequest]()
	respDefault := NewDataSpecOf[DefaultRestResponse]()
	src.requests[MethodCreate] = reqDefault
	src.responses[201] = respDefault
	src.responseEnvelopes[201] = NewEnvelopeSpec("data", "meta")
	src.informations[101] = respDefault
	src.redirections[301] = respDefault
	src.problems[400] = respDefault

	mergeRegistryOptions(dst, src)

	if dst.requests[MethodCreate] == nil || dst.responses[201] == nil || dst.responseEnvelopes[201] == nil {
		t.Fatal("expected merged config to include request/response maps")
	}
	if dst.informations[101] == nil || dst.redirections[301] == nil || dst.problems[400] == nil {
		t.Fatal("expected merged config to include status class maps")
	}
}

// --- resolveResponseInstance ---

func TestResolveResponseInstance_Success(t *testing.T) {
	instance, err := resolveResponseInstance(NewDataSpecOf[DefaultRestResponse]())
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
	resolved, err := resolveResponseSlice(NewDataSpecOf[DefaultRestResponse](), 3)
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
	source := map[int]RestDataSpec{
		defaultStatusCode:  NewDataSpecOf[DefaultRestResponse](),
		http.StatusCreated: NewDataSpecOf[sampleT](),
	}
	resolved := resolveResponseSpec(http.StatusCreated, source)
	if _, ok := resolved.(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected dataSpecOf[sampleT] for direct hit, got %T", resolved)
	}
}

func TestResolveResponseSpec_FallbackToDefault(t *testing.T) {
	source := map[int]RestDataSpec{
		defaultStatusCode: NewDataSpecOf[DefaultRestResponse](),
	}
	resolved := resolveResponseSpec(http.StatusOK, source)
	if _, ok := resolved.(*dataSpecOf[DefaultRestResponse]); !ok {
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

func (n *nilNewSpec) New() RestDataSpec             { return nil }
func (n *nilNewSpec) NewSlice(_ int) []RestDataSpec { return nil }
func (n *nilNewSpec) SetBody(_ any) error           { return nil }
func (n *nilNewSpec) Body() any                     { return nil }
func (n *nilNewSpec) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) }
func (n *nilNewSpec) UnmarshalJSON(_ []byte) error  { return nil }

type badSliceSpec struct{}

func (b *badSliceSpec) New() RestDataSpec             { return &badSliceSpec{} }
func (b *badSliceSpec) NewSlice(_ int) []RestDataSpec { return []RestDataSpec{} }
func (b *badSliceSpec) SetBody(_ any) error           { return nil }
func (b *badSliceSpec) Body() any                     { return nil }
func (b *badSliceSpec) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) }
func (b *badSliceSpec) UnmarshalJSON(_ []byte) error  { return nil }
