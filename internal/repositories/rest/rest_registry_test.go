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
func (f *foreignRegistry) ResolveRequest(body RestRequestSpec, requestBody *[]byte) error {
	return nil
}

// ResolveEnvelopeRequest implements the RestRegistry contract using RestMethod.
func (f *foreignRegistry) ResolveEnvelopeRequest(_ RestMethod, dataBody *[]byte) error { return nil }
func (f *foreignRegistry) ResolveResponse(status int, responseBody []byte, body *RestResponseSpec) error {
	return nil
}
func (f *foreignRegistry) ResolveResponses(status int, responseBody []byte, bodies *[]RestResponseSpec) error {
	return nil
}
func (f *foreignRegistry) ResolveEnvelopeResponse(status int, dataBody *[]byte, meta *types.ResponseMeta) error {
	return nil
}

func (f *foreignRegistry) NewRequestSpec(_ RestMethod) (RestRequestSpec, error) {
	return &DefaultRestRequest{}, nil
}
func (f *foreignRegistry) ReleaseRequestSpec(spec RestRequestSpec) {}

type testEnvelopeRequest struct {
	Data json.RawMessage `json:"data"`
}

func (t *testEnvelopeRequest) New() RestRequestSpec {
	return &testEnvelopeRequest{}
}

func (t *testEnvelopeRequest) SetBody(body RestRequestSpec) {
	if raw, ok := body.(*DefaultRestRequest); ok {
		t.Data = append([]byte(nil), raw.Body...)
	}
}

type testResponse struct {
	Name string `json:"name"`
}

func (t *testResponse) New() RestResponseSpec {
	return &testResponse{}
}

func (t *testResponse) NewSlice(n int) []RestResponseSpec {
	out := make([]RestResponseSpec, n)
	for i := range out {
		out[i] = &testResponse{}
	}
	return out
}

func TestResolveEnvelopeRequest(t *testing.T) {
	registry := NewRestRegistry(WithRequestEnvelope(&testEnvelopeRequest{}, MethodCreate))

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
		WithResponse(&testResponse{}, http.StatusOK),
		WithResponseEnvelope(&responseEnvelope{}, http.StatusOK),
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

	var bodies []RestResponseSpec
	if err := registry.ResolveResponses(http.StatusOK, responseBody, &bodies); err != nil {
		t.Fatalf("resolve responses failed: %v", err)
	}
	if len(bodies) != 2 {
		t.Fatalf("unexpected bodies length: got %d, want 2", len(bodies))
	}

	first, ok := bodies[0].(*testResponse)
	if !ok {
		t.Fatalf("unexpected first body type: %T", bodies[0])
	}
	if first.Name != "a" {
		t.Fatalf("unexpected first body name: got %s, want a", first.Name)
	}
}

func TestResolveResponseUsesDefaultRawBody(t *testing.T) {
	registry := NewRestRegistry()

	var body RestResponseSpec
	raw := []byte(`{"hello":"world"}`)
	if err := registry.ResolveResponse(http.StatusOK, raw, &body); err != nil {
		t.Fatalf("resolve response failed: %v", err)
	}

	defaultBody, ok := body.(*DefaultRestResponse)
	if !ok {
		t.Fatalf("unexpected body type: %T", body)
	}
	if !bytes.Equal(defaultBody.Body, raw) {
		t.Fatalf("unexpected default raw body: got %s, want %s", defaultBody.Body, raw)
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
	left := NewRestRegistry(WithRequest(&DefaultRestRequest{}, MethodCreate))
	right := NewRestRegistry(WithRequest(&testEnvelopeRequest{}, MethodUpdate))

	merged := left.Merge(right)

	postSpec, err := merged.NewRequestSpec(MethodCreate)
	if err != nil {
		t.Fatalf("unexpected error for post spec: %v", err)
	}
	if _, ok := postSpec.(*DefaultRestRequest); !ok {
		t.Fatalf("unexpected post spec type: %T", postSpec)
	}

	putSpec, err := merged.NewRequestSpec(MethodUpdate)
	if err != nil {
		t.Fatalf("unexpected error for put spec: %v", err)
	}
	if _, ok := putSpec.(*testEnvelopeRequest); !ok {
		t.Fatalf("unexpected put spec type: %T", putSpec)
	}
}

func TestMergeWithNilAndForeignRegistryReturnsClone(t *testing.T) {
	left := NewRestRegistry(WithRequest(&DefaultRestRequest{}, MethodCreate))

	mergedNil := left.Merge(nil)
	if _, err := mergedNil.NewRequestSpec(MethodCreate); err != nil {
		t.Fatalf("expected merged registry with nil other to preserve entries: %v", err)
	}

	mergedForeign := left.Merge(&foreignRegistry{})
	if _, err := mergedForeign.NewRequestSpec(MethodCreate); err != nil {
		t.Fatalf("expected merged registry with foreign other to preserve entries: %v", err)
	}
}

func TestResolveRequest_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry()
	if err := registry.ResolveRequest(&DefaultRestRequest{}, nil); err == nil {
		t.Fatal("expected error for nil request body pointer")
	}

	payload := []byte{}
	if err := registry.ResolveRequest(nil, &payload); err == nil {
		t.Fatal("expected error for nil request spec")
	}

	if err := registry.ResolveRequest(&marshalFailRequest{}, &payload); err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestResolveEnvelopeRequest_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithRequestEnvelope(&nilEnvelopeRequest{}, MethodCreate))
	payload := []byte(`{"x":1}`)
	if err := registry.ResolveEnvelopeRequest(MethodCreate, nil); err == nil {
		t.Fatal("expected error for nil dataBody")
	}
	if err := registry.ResolveEnvelopeRequest(MethodCreate, &payload); err == nil {
		t.Fatal("expected error when envelope New returns nil")
	}

	registry = NewRestRegistry(WithRequestEnvelope(&marshalFailEnvelopeRequest{}, MethodCreate))
	if err := registry.ResolveEnvelopeRequest(MethodCreate, &payload); err == nil {
		t.Fatal("expected marshal error for envelope")
	}
}

func TestResolveResponse_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithResponse(&testResponse{}, http.StatusOK))
	if err := registry.ResolveResponse(http.StatusOK, []byte(`{"name":"ok"}`), nil); err == nil {
		t.Fatal("expected error for nil body pointer")
	}

	var out RestResponseSpec
	if err := registry.ResolveResponse(http.StatusOK, []byte(`not-json`), &out); err == nil {
		t.Fatal("expected unmarshal error")
	}

	registry = NewRestRegistry(WithResponse(&nilNewResponseSpec{}, http.StatusOK))
	if err := registry.ResolveResponse(http.StatusOK, []byte(`{"name":"ok"}`), &out); err == nil {
		t.Fatal("expected error when response New returns nil")
	}

	registry = NewRestRegistry(WithResponse(&testResponse{}, http.StatusOK))
	if err := registry.ResolveResponse(http.StatusContinue+600, []byte(`{"name":"ok"}`), &out); err == nil {
		t.Fatal("expected error for status class outside supported ranges")
	}
}

func TestResolveResponses_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithResponse(&testResponse{}, http.StatusOK))
	if err := registry.ResolveResponses(http.StatusOK, []byte(`[]`), nil); err == nil {
		t.Fatal("expected error for nil bodies pointer")
	}

	var out []RestResponseSpec
	if err := registry.ResolveResponses(http.StatusOK, []byte(`not-json`), &out); err == nil {
		t.Fatal("expected unmarshal raw list error")
	}

	registry = NewRestRegistry(WithResponse(&nilItemSliceResponse{}, http.StatusOK))
	if err := registry.ResolveResponses(http.StatusOK, []byte(`[{}]`), &out); err == nil {
		t.Fatal("expected error when NewSlice contains nil item")
	}

	registry = NewRestRegistry(WithResponse(&testResponse{}, http.StatusOK))
	if err := registry.ResolveResponses(http.StatusOK, []byte(`[1]`), &out); err == nil {
		t.Fatal("expected item unmarshal error")
	}

	registry = NewRestRegistry(WithResponse(&badLenSliceResponse{}, http.StatusOK))
	if err := registry.ResolveResponses(http.StatusOK, []byte(`[{}]`), &out); err == nil {
		t.Fatal("expected error when NewSlice length mismatches payload count")
	}

	registry = NewRestRegistry(WithResponse(&testResponse{}, http.StatusOK))
	if err := registry.ResolveResponses(http.StatusContinue+600, []byte(`[{}]`), &out); err == nil {
		t.Fatal("expected error for status class outside supported ranges")
	}
}

func TestResolveEnvelopeResponse_ErrorPaths(t *testing.T) {
	registry := NewRestRegistry(WithResponseEnvelope(&responseEnvelope{}, http.StatusOK))
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

	registry = NewRestRegistry(WithResponseEnvelope(&nilNewEnvelopeSpec{}, http.StatusOK))
	if err := registry.ResolveEnvelopeResponse(http.StatusOK, &body, &meta); err == nil {
		t.Fatal("expected error when envelope New returns nil")
	}
}

func TestNewRequestSpecAndReleaseRequestSpec(t *testing.T) {
	registry := NewRestRegistry(WithRequest(&nilRequestSpec{}, MethodCreate))
	if _, err := registry.NewRequestSpec(MethodCreate); err == nil {
		t.Fatal("expected error when request New returns nil")
	}

	registry = NewRestRegistry()
	if _, err := registry.NewRequestSpec(MethodCreate); err != nil {
		t.Fatalf("unexpected error creating default request spec: %v", err)
	}

	registry.ReleaseRequestSpec(nil)
}

type marshalFailRequest struct {
	C chan int `json:"c"`
}

func (m *marshalFailRequest) New() RestRequestSpec         { return &marshalFailRequest{} }
func (m *marshalFailRequest) SetBody(body RestRequestSpec) {}

type nilEnvelopeRequest struct{}

func (n *nilEnvelopeRequest) New() RestRequestSpec         { return nil }
func (n *nilEnvelopeRequest) SetBody(body RestRequestSpec) {}

type marshalFailEnvelopeRequest struct {
	C chan int `json:"c"`
}

func (m *marshalFailEnvelopeRequest) New() RestRequestSpec         { return &marshalFailEnvelopeRequest{} }
func (m *marshalFailEnvelopeRequest) SetBody(body RestRequestSpec) {}

type nilItemSliceResponse struct{}

func (n *nilItemSliceResponse) New() RestResponseSpec { return &nilItemSliceResponse{} }
func (n *nilItemSliceResponse) NewSlice(nItems int) []RestResponseSpec {
	out := make([]RestResponseSpec, nItems)
	for i := range out {
		out[i] = nil
	}
	return out
}

type nilNewResponseSpec struct{}

func (n *nilNewResponseSpec) New() RestResponseSpec             { return nil }
func (n *nilNewResponseSpec) NewSlice(_ int) []RestResponseSpec { return nil }

type badLenSliceResponse struct{}

func (b *badLenSliceResponse) New() RestResponseSpec { return &badLenSliceResponse{} }
func (b *badLenSliceResponse) NewSlice(_ int) []RestResponseSpec {
	return []RestResponseSpec{}
}

type nilNewEnvelopeSpec struct{}

func (n *nilNewEnvelopeSpec) New() RestEnvelopeSpec            { return nil }
func (n *nilNewEnvelopeSpec) EnvelopeData() json.RawMessage    { return json.RawMessage(`null`) }
func (n *nilNewEnvelopeSpec) EnvelopeMeta() types.ResponseMeta { return types.ResponseMeta{} }

type nilRequestSpec struct{}

func (n *nilRequestSpec) New() RestRequestSpec         { return nil }
func (n *nilRequestSpec) SetBody(body RestRequestSpec) {}
