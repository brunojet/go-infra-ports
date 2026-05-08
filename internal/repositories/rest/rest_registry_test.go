package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type foreignRegistry struct{}

func (f *foreignRegistry) Merge(other RestRegistry) RestRegistry { return other }
func (f *foreignRegistry) ResolveRequest(body RestDataSpec, requestBody *[]byte) error {
	return nil
}
func (f *foreignRegistry) ResolveEnvelopeRequest(_ RestMethod, dataBody *[]byte) error { return nil }
func (f *foreignRegistry) ResolveResponse(status int, responseBody []byte, body *RestDataSpec) error {
	return nil
}
func (f *foreignRegistry) ResolveResponses(status int, responseBody []byte, bodies *[]RestDataSpec) error {
	return nil
}
func (f *foreignRegistry) ResolveEnvelopeResponse(status int, dataBody *[]byte, meta *types.ResponseMeta) error {
	return nil
}
func (f *foreignRegistry) NewRequestSpec(_ RestMethod) (RestDataSpec, error) {
	return NewDataSpecOf[DefaultRestRequest](), nil
}
func (f *foreignRegistry) ReleaseRequestSpec(spec RestDataSpec) {}

type testResponse struct {
	Name string `json:"name"`
}

func TestResolveEnvelopeRequest(t *testing.T) {
	registry := NewRestRegistry(WithRequestEnvelope("data", MethodCreate))

	payload := []byte(`{"name":"jet"}`)
	if err := registry.ResolveEnvelopeRequest(MethodCreate, &payload); err != nil {
		t.Fatalf("resolve envelope request failed: %v", err)
	}

	expected := []byte(`{"data":{"name":"jet"}}`)
	if !bytes.Equal(payload, expected) {
		t.Fatalf("unexpected envelope payload: got %s, want %s", payload, expected)
	}
}

func TestResolveEnvelopeRequest_NoSetup_ReturnsSilently(t *testing.T) {
	registry := NewRestRegistry()
	payload := []byte(`{"name":"jet"}`)

	if err := registry.ResolveEnvelopeRequest(MethodCreate, &payload); err != nil {
		t.Fatalf("expected no error when envelope is not configured, got %v", err)
	}

	expected := []byte(`{"name":"jet"}`)
	if !bytes.Equal(payload, expected) {
		t.Fatalf("expected payload to remain unchanged: got %s, want %s", payload, expected)
	}
}

func TestResolveEnvelopeResponseAndResolveResponses(t *testing.T) {
	registry := NewRestRegistry(
		WithResponseOf[testResponse](http.StatusOK),
		WithResponseEnvelope("data", "meta", http.StatusOK),
	)

	responseBody := []byte(`{"data":[{"name":"a"},{"name":"b"}],"meta":{"page":1}}`)
	meta := types.ResponseMeta{}

	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &responseBody, &meta); err != nil {
		t.Fatalf("resolve envelope response failed: %v", err)
	}

	expectedData := []byte(`[{"name":"a"},{"name":"b"}]`)
	if !bytes.Equal(responseBody, expectedData) {
		t.Fatalf("unexpected extracted data body: got %s, want %s", responseBody, expectedData)
	}
	if meta["page"] != float64(1) {
		t.Fatalf("unexpected meta value: got %v", meta["page"])
	}

	var bodies []RestDataSpec
	if err := registry.ResolveResponses(http.StatusOK, responseBody, &bodies); err != nil {
		t.Fatalf("resolve responses failed: %v", err)
	}
	if len(bodies) != 2 {
		t.Fatalf("unexpected bodies length: got %d, want 2", len(bodies))
	}

	firstAny := bodies[0].Body()
	first, ok := firstAny.(testResponse)
	if !ok {
		t.Fatalf("unexpected first body type: %T", bodies[0].Body())
	}
	if first.Name != "a" {
		t.Fatalf("unexpected first body name: got %s, want a", first.Name)
	}
}

func TestResolveResponseUsesDefaultRawBody(t *testing.T) {
	registry := NewRestRegistry()

	var body RestDataSpec
	raw := []byte(`{"hello":"world"}`)
	if err := registry.ResolveResponse(http.StatusOK, raw, &body); err != nil {
		t.Fatalf("resolve response failed: %v", err)
	}

	// Default response spec uses DefaultRestResponse (alias to json.RawMessage)
	if rb, ok := body.Body().(DefaultRestResponse); ok {
		if !bytes.Equal([]byte(rb), raw) {
			t.Fatalf("unexpected default raw body: got %s, want %s", []byte(rb), raw)
		}
	} else {
		t.Fatalf("unexpected body type: %T", body.Body())
	}
}

func TestResolveEnvelopeResponse_NoSetup_Passthrough(t *testing.T) {
	registry := NewRestRegistry()

	original := []byte(`{"data":{"name":"jet"},"meta":{"page":1}}`)
	body := append([]byte(nil), original...)
	meta := types.ResponseMeta{}

	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &body, &meta); err != nil {
		t.Fatalf("expected no error when envelope is not configured, got %v", err)
	}
	if !bytes.Equal(body, original) {
		t.Fatalf("expected body to remain unchanged: got %s, want %s", body, original)
	}
	if len(meta) != 0 {
		t.Fatalf("expected meta to remain empty, got %v", meta)
	}
}

func TestMergeCombinesRegistries(t *testing.T) {
	left := NewRestRegistry(WithRequestOf[DefaultRestRequest](MethodCreate))
	right := NewRestRegistry(WithRequestOf[testResponse](MethodUpdate))

	merged := left.Merge(right)

	postSpec, err := merged.NewRequestSpec(MethodCreate)
	if err != nil {
		t.Fatalf("unexpected error for post spec: %v", err)
	}
	if _, ok := postSpec.(*dataSpecOf[DefaultRestRequest]); !ok {
		t.Fatalf("unexpected post spec type: %T", postSpec)
	}

	putSpec, err := merged.NewRequestSpec(MethodUpdate)
	if err != nil {
		t.Fatalf("unexpected error for put spec: %v", err)
	}
	if _, ok := putSpec.(*dataSpecOf[testResponse]); !ok {
		t.Fatalf("unexpected put spec type: %T", putSpec)
	}
}

func expectPanic(t *testing.T, fn func(), msg string) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic: %s", msg)
		}
	}()
	fn()
}

func TestMergeWithNilAndForeignRegistryPanics(t *testing.T) {
	left := NewRestRegistry(WithRequestOf[DefaultRestRequest](MethodCreate))

	expectPanic(t, func() { left.Merge(nil) }, "nil other")
	expectPanic(t, func() { left.Merge(&foreignRegistry{}) }, "foreign other")
}

// --- error-path tests ---

type marshalFailRequest struct {
	C chan int `json:"c"`
}

type nilNewEnvelopeSpec struct{}

func (n *nilNewEnvelopeSpec) New() RestEnvelopeSpec             { return nil }
func (n *nilNewEnvelopeSpec) EnvelopeData() json.RawMessage     { return json.RawMessage(`null`) }
func (n *nilNewEnvelopeSpec) EnvelopeMeta() types.ResponseMeta  { return types.ResponseMeta{} }
func (n *nilNewEnvelopeSpec) SetEnvelopeData(_ json.RawMessage) {}
func (n *nilNewEnvelopeSpec) MarshalJSON() ([]byte, error)      { return json.Marshal(nil) }
func (n *nilNewEnvelopeSpec) UnmarshalJSON(_ []byte) error      { return nil }

func TestResolveRequest_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry()
	if err := registry.ResolveRequest(NewDataSpecOf[DefaultRestRequest](), nil); err == nil {
		t.Fatal("expected error for nil request body pointer")
	}

	payload := []byte{}
	if err := registry.ResolveRequest(nil, &payload); err == nil {
		t.Fatal("expected error for nil request spec")
	}

	if err := registry.ResolveRequest(NewDataSpecOf[marshalFailRequest](), &payload); err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestResolveEnvelopeRequest_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithRequestEnvelope("data", MethodCreate))
	payload := []byte(`{"x":1}`)
	if err := registry.ResolveEnvelopeRequest(MethodCreate, nil); err == nil {
		t.Fatal("expected error for nil dataBody")
	}

	// replace envelope with one whose New returns nil
	if rImpl, ok := registry.(*restRegistry); ok {
		rImpl.cfg.requestsEnvelopes[MethodCreate] = &nilNewEnvelopeSpec{}
	} else {
		t.Fatal("expected concrete *restRegistry to modify cfg")
	}
	if err := registry.ResolveEnvelopeRequest(MethodCreate, &payload); err == nil {
		t.Fatal("expected error when envelope New returns nil")
	}
}

func TestResolveEnvelopeResponse_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithResponseEnvelope("data", "meta", http.StatusOK))
	body := []byte(`{"data":{"x":1},"meta":{}}`)
	meta := types.ResponseMeta{}

	if err := registry.ResolveEnvelopeResponse(http.StatusOK, nil, &meta); err == nil {
		t.Fatal("expected error for nil dataBody")
	}
	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &body, nil); err == nil {
		t.Fatal("expected error for nil meta")
	}
	badJSON := []byte(`not-json`)
	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &badJSON, &meta); err == nil {
		t.Fatal("expected unmarshal error")
	}

	registry = NewRestRegistry(WithResponseEnvelope("data", "meta", http.StatusOK))
	if rImpl, ok := registry.(*restRegistry); ok {
		rImpl.cfg.responseEnvelopes[http.StatusOK] = &nilNewEnvelopeSpec{}
	} else {
		t.Fatal("expected concrete *restRegistry to modify cfg")
	}
	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &body, &meta); err == nil {
		t.Fatal("expected error when envelope New returns nil")
	}
}
