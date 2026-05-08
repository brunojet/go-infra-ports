package rest

import (
	"bytes"
	"testing"
)

type payloadT struct {
	Name string `json:"name"`
}

func TestDataSpec_NewAndSetBody(t *testing.T) {
	spec := NewDataSpecOf[payloadT]()
	if spec == nil {
		t.Fatal("NewDataSpecOf returned nil")
	}
	created := spec.New()
	if created == nil {
		t.Fatal("New() returned nil")
	}

	// Assert concrete type and then exercise SetBody via concrete pointer
	ds, ok := spec.(*dataSpecOf[payloadT])
	if !ok {
		t.Fatalf("unexpected type: %T", spec)
	}
	// Set typed value
	if err := ds.SetBody(payloadT{Name: "alpha"}); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}
	if ds.Prototype.Name != "alpha" {
		t.Fatalf("unexpected prototype: got %v", ds.Prototype)
	}
	// Validate Body() returns the current prototype value
	body := ds.Body()
	pb, ok := body.(payloadT)
	if !ok {
		t.Fatalf("unexpected Body() type: %T", body)
	}
	if pb.Name != "alpha" {
		t.Fatalf("unexpected Body() value: %+v", pb)
	}

	// Call SetBody(nil) to exercise nil path. Current implementation
	// returns an error when given a nil `any` value; assert that behavior
	// so tests remain consistent with the application code.
	var nilAny any
	if err := ds.SetBody(nilAny); err == nil {
		t.Fatalf("expected error for SetBody(nil), got nil")
	}
}

func TestDataSpec_UnmarshalJSONAndNewSlice(t *testing.T) {
	var d dataSpecOf[payloadT]
	if got := d.New(); got == nil {
		t.Fatal("New() returned nil")
	}

	slice := d.NewSlice(3)
	if len(slice) != 3 {
		t.Fatalf("unexpected slice length: got %d, want 3", len(slice))
	}
	if slice[0] == slice[1] {
		t.Fatal("NewSlice() returned duplicate elements")
	}

	// Unmarshal JSON into data spec
	resp := &dataSpecOf[payloadT]{}
	raw := []byte(`{"name":"jet"}`)
	if err := resp.UnmarshalJSON(raw); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if resp.Prototype.Name != "jet" {
		t.Fatalf("unexpected prototype after unmarshal: %v", resp.Prototype)
	}
}

func TestDataSpec_PrototypePreserved(t *testing.T) {
	proto := payloadT{Name: "proto"}
	spec := NewDataSpecOf[payloadT]()
	// constructor is no-arg; set the prototype via SetBody to preserve semantics
	if err := spec.SetBody(proto); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}
	if spec == nil {
		t.Fatal("NewDataSpecOf returned nil")
	}
	rs, ok := spec.(*dataSpecOf[payloadT])
	if !ok {
		t.Fatalf("unexpected type: %T", spec)
	}
	if rs.Prototype.Name != "proto" {
		t.Fatalf("unexpected prototype preserved: %v", rs.Prototype)
	}
}

func TestEnvelopeSpec_UnmarshalAndAccessors(t *testing.T) {
	env := NewEnvelopeSpec("data", "meta")
	re, ok := env.(*envelopeSpecOf)
	if !ok {
		t.Fatalf("unexpected type: %T", env)
	}

	raw := []byte(`{"data":{"x":1},"meta":{"k":"v"}}`)
	if err := re.UnmarshalJSON(raw); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if !bytes.Contains(re.Data, []byte(`"x":1`)) {
		t.Fatalf("unexpected data payload: %s", string(re.Data))
	}
	meta := re.EnvelopeMeta()
	if meta["k"] != "v" {
		t.Fatalf("unexpected meta value: %v", meta["k"])
	}

	// When meta field is not configured, Unmarshal should accept data-only envelopes
	env2 := NewEnvelopeSpec("data", "")
	re2 := env2.(*envelopeSpecOf)
	if err := re2.UnmarshalJSON([]byte(`{"data":null}`)); err != nil {
		t.Fatalf("unexpected error when meta not configured: %v", err)
	}
	if string(re2.EnvelopeData()) != "null" {
		t.Fatalf("unexpected envelope data fallback: %s", string(re2.EnvelopeData()))
	}
	if len(re2.EnvelopeMeta()) != 0 {
		t.Fatalf("expected empty meta fallback, got %v", re2.EnvelopeMeta())
	}

	// Missing data field should return an error
	re3 := NewEnvelopeSpec("body", "meta").(*envelopeSpecOf)
	if err := re3.UnmarshalJSON([]byte(`{"meta":{}}`)); err == nil {
		t.Fatalf("expected error for missing data field")
	}
}

func TestEnvelopeSpec_EnvelopeDataNilFallback(t *testing.T) {
	env := NewEnvelopeSpec("data", "meta")
	re := env.(*envelopeSpecOf)
	if string(re.EnvelopeData()) != "null" {
		t.Fatalf("expected null fallback for nil data, got %s", string(re.EnvelopeData()))
	}
}

func TestDataSpec_MarshalJSON(t *testing.T) {
	r := &dataSpecOf[payloadT]{Prototype: payloadT{Name: "pack"}}
	b, err := r.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if !bytes.Contains(b, []byte(`"name":"pack"`)) {
		t.Fatalf("unexpected marshal output: %s", string(b))
	}
}

func TestDataSpec_UnmarshalJSON_Error(t *testing.T) {
	var r dataSpecOf[payloadT]
	if err := r.UnmarshalJSON([]byte(`{"name":`)); err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestEnvelopeSpec_DefaultDataField(t *testing.T) {
	env := NewEnvelopeSpec("", "meta")
	re, ok := env.(*envelopeSpecOf)
	if !ok {
		t.Fatalf("unexpected type: %T", env)
	}
	if re.DataField != "data" {
		t.Fatalf("expected DataField 'data', got %q", re.DataField)
	}
}

func TestEnvelopeSpec_NewCopiesFields(t *testing.T) {
	env := NewEnvelopeSpec("d", "m")
	re, ok := env.(*envelopeSpecOf)
	if !ok {
		t.Fatalf("unexpected type: %T", env)
	}
	created := re.New()
	nr, ok := created.(*envelopeSpecOf)
	if !ok {
		t.Fatalf("unexpected new type: %T", created)
	}
	if nr.DataField != "d" || nr.MetaField != "m" {
		t.Fatalf("fields not copied: got %v/%v", nr.DataField, nr.MetaField)
	}
}

func TestEnvelopeSpec_UnmarshalJSON_Errors(t *testing.T) {
	// invalid JSON
	re1 := NewEnvelopeSpec("data", "meta").(*envelopeSpecOf)
	if err := re1.UnmarshalJSON([]byte(`{`)); err == nil {
		t.Fatal("expected error for invalid json")
	}

	// meta is not an object
	re2 := NewEnvelopeSpec("data", "meta").(*envelopeSpecOf)
	if err := re2.UnmarshalJSON([]byte(`{"data":{"x":1},"meta":"notobj"}`)); err == nil {
		t.Fatal("expected error for meta unmarshal failure")
	}

	// missing meta field
	re3 := NewEnvelopeSpec("data", "meta").(*envelopeSpecOf)
	if err := re3.UnmarshalJSON([]byte(`{"data":{"x":1}}`)); err == nil {
		t.Fatal("expected error for missing meta field")
	}
}
