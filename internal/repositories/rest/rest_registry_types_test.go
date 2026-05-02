package rest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

func TestDefaultRestRequest_NewAndSetBody(t *testing.T) {
	req := &DefaultRestRequest{}
	spec := req.New()
	if spec == nil {
		t.Fatal("New() returned nil")
	}
	payload := json.RawMessage(`{"key":"value"}`)
	spec.SetBody(payload)

	cast, ok := spec.(*DefaultRestRequest)
	if !ok {
		t.Fatalf("unexpected type: %T", spec)
	}
	if !bytes.Equal(cast.Body, payload) {
		t.Fatalf("unexpected body: got %s, want %s", cast.Body, payload)
	}
}

func TestDefaultRestResponse_NewReturnsDistinctInstance(t *testing.T) {
	resp := &DefaultRestResponse{}
	a := resp.New()
	b := resp.New()
	if a == b {
		t.Fatal("New() returned the same instance twice")
	}
}

func TestDefaultRestResponse_UnmarshalJSON_StoresRawBytes(t *testing.T) {
	raw := []byte(`{"name":"jet"}`)
	resp := &DefaultRestResponse{}
	if err := json.Unmarshal(raw, resp); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if !bytes.Equal(resp.Body, raw) {
		t.Fatalf("unexpected body: got %s, want %s", resp.Body, raw)
	}
}

func TestDefaultRestResponse_UnmarshalJSON_IsolatesBytes(t *testing.T) {
	raw := []byte(`{"x":1}`)
	resp := &DefaultRestResponse{}
	if err := json.Unmarshal(raw, resp); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	raw[0] = 'Z'
	if resp.Body[0] == 'Z' {
		t.Fatal("UnmarshalJSON did not copy bytes; mutation affected Body")
	}
}

func TestDefaultRestResponse_NewSlice_Length(t *testing.T) {
	resp := &DefaultRestResponse{}
	slice := resp.NewSlice(3)
	if len(slice) != 3 {
		t.Fatalf("unexpected slice length: got %d, want 3", len(slice))
	}
}

func TestDefaultRestResponse_NewSlice_ElementsAreDistinct(t *testing.T) {
	resp := &DefaultRestResponse{}
	slice := resp.NewSlice(2)
	if slice[0] == slice[1] {
		t.Fatal("NewSlice() returned duplicate pointers")
	}
}

func TestResponseEnvelope_NewReturnsEnvelope(t *testing.T) {
	env := &responseEnvelope{}
	created := env.New()
	if _, ok := created.(*responseEnvelope); !ok {
		t.Fatalf("unexpected envelope type: %T", created)
	}
	if _, ok := any(created).(RestEnvelopeSpec); !ok {
		t.Fatalf("expected created envelope to satisfy RestEnvelopeSpec, got %T", created)
	}
}

func TestResponseEnvelope_EnvelopeData(t *testing.T) {
	env := &responseEnvelope{Data: json.RawMessage(`{"x":1}`)}
	if !bytes.Equal(env.EnvelopeData(), json.RawMessage(`{"x":1}`)) {
		t.Fatalf("unexpected envelope data: %s", env.EnvelopeData())
	}

	env.Data = nil
	if string(env.EnvelopeData()) != "null" {
		t.Fatalf("unexpected nil envelope data fallback: %s", env.EnvelopeData())
	}
}

func TestResponseEnvelope_EnvelopeMeta(t *testing.T) {
	env := &responseEnvelope{Meta: types.ResponseMeta{"k": "v"}}
	if env.EnvelopeMeta()["k"] != "v" {
		t.Fatalf("unexpected meta value: %v", env.EnvelopeMeta()["k"])
	}

	env.Meta = nil
	meta := env.EnvelopeMeta()
	if len(meta) != 0 {
		t.Fatalf("expected empty meta fallback, got %v", meta)
	}
}
