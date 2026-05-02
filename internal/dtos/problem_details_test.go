package dtos

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalRawField_StatusSuccessDeletesKey(t *testing.T) {
	raw := map[string]json.RawMessage{
		"status": json.RawMessage("422"),
	}
	var status int

	if err := unmarshalRawField(raw, "status", &status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 422 {
		t.Fatalf("expected status 422, got %d", status)
	}
	if _, ok := raw["status"]; ok {
		t.Fatalf("expected status key removed from raw map")
	}
}

func TestUnmarshalRawField_StatusInvalidTypeKeepsKey(t *testing.T) {
	raw := map[string]json.RawMessage{
		"status": json.RawMessage(`"bad"`),
	}
	var status int

	if err := unmarshalRawField(raw, "status", &status); err == nil {
		t.Fatalf("expected error for invalid status type")
	}
	if _, ok := raw["status"]; !ok {
		t.Fatalf("expected status key to remain when decode fails")
	}
}

func TestUnmarshalRawField_MissingKeyNoop(t *testing.T) {
	raw := map[string]json.RawMessage{
		"title": json.RawMessage(`"x"`),
	}
	var status int

	if err := unmarshalRawField(raw, "status", &status); err != nil {
		t.Fatalf("unexpected error for missing key: %v", err)
	}
	if status != 0 {
		t.Fatalf("expected unchanged zero status, got %d", status)
	}
	if _, ok := raw["title"]; !ok {
		t.Fatalf("expected unrelated keys to remain untouched")
	}
}

func TestUnmarshalExtensions_Empty(t *testing.T) {
	extensions, err := unmarshalExtensions(map[string]json.RawMessage{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if extensions != nil {
		t.Fatalf("expected nil extensions for empty raw map, got %#v", extensions)
	}
}

func TestUnmarshalExtensions_Success(t *testing.T) {
	raw := map[string]json.RawMessage{
		"trace_id": json.RawMessage(`"abc"`),
		"retries":  json.RawMessage(`2`),
	}
	extensions, err := unmarshalExtensions(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if extensions["trace_id"] != "abc" {
		t.Fatalf("expected trace_id abc, got %#v", extensions["trace_id"])
	}
	if extensions["retries"] != float64(2) {
		t.Fatalf("expected retries 2, got %#v", extensions["retries"])
	}
}

func TestUnmarshalExtensions_InvalidValue(t *testing.T) {
	raw := map[string]json.RawMessage{
		"broken": json.RawMessage(`{"x":`),
	}
	if _, err := unmarshalExtensions(raw); err == nil {
		t.Fatalf("expected error for invalid extension JSON")
	}
}

func TestProblemDetails_MarshalJSON_EmitsRFC7807FieldsAndExtensions(t *testing.T) {
	p := ProblemDetails{
		Type:     "urn:problem:test",
		Title:    "Invalid",
		Status:   422,
		Detail:   "bad payload",
		Instance: "/users/1",
		Extensions: map[string]any{
			"trace_id": "abc",
			"status":   999, // must not override base field
		},
	}

	b, err := json.Marshal(&p)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if got["type"] != "urn:problem:test" || got["title"] != "Invalid" || got["detail"] != "bad payload" || got["instance"] != "/users/1" {
		t.Fatalf("unexpected base fields: %#v", got)
	}
	if got["status"] != float64(422) {
		t.Fatalf("expected status 422, got %#v", got["status"])
	}
	if got["trace_id"] != "abc" {
		t.Fatalf("expected extension trace_id, got %#v", got["trace_id"])
	}
}

func TestProblemDetails_UnmarshalJSON_ParsesExtensions(t *testing.T) {
	data := []byte(`{"type":"urn:problem:test","title":"Invalid","status":400,"detail":"bad request","instance":"/users","trace_id":"abc","retries":2}`)
	var p ProblemDetails
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if p.Type != "urn:problem:test" || p.Title != "Invalid" || p.Status != 400 || p.Detail != "bad request" || p.Instance != "/users" {
		t.Fatalf("unexpected problem details: %#v", p)
	}
	if p.Extensions["trace_id"] != "abc" || p.Extensions["retries"] != float64(2) {
		t.Fatalf("unexpected extensions: %#v", p.Extensions)
	}
}

func TestProblemDetails_UnmarshalJSON_NoExtensions(t *testing.T) {
	data := []byte(`{"type":"urn:test","title":"T","status":200,"detail":"d","instance":"/x"}`)
	var p ProblemDetails
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if p.Extensions != nil {
		t.Fatalf("expected nil extensions, got %#v", p.Extensions)
	}
}

func TestProblemDetails_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var p ProblemDetails
	if err := p.UnmarshalJSON([]byte("not-json")); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func TestProblemDetails_UnmarshalJSON_InvalidFieldType(t *testing.T) {
	var p ProblemDetails
	if err := json.Unmarshal([]byte(`{"status":"not-a-number"}`), &p); err == nil {
		t.Fatalf("expected error for wrong status type")
	}
}

func TestProblemDetails_UnmarshalJSON_InvalidKnownFieldTypes(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "type", data: `{"type":123}`},
		{name: "title", data: `{"title":123}`},
		{name: "detail", data: `{"detail":123}`},
		{name: "instance", data: `{"instance":123}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p ProblemDetails
			if err := json.Unmarshal([]byte(tt.data), &p); err == nil {
				t.Fatalf("expected error for wrong %s type", tt.name)
			}
		})
	}
}

func TestProblemDetails_ContentType(t *testing.T) {
	var p ProblemDetails
	if got := p.ContentType(); got != ContentTypeProblemJSON {
		t.Fatalf("expected %q, got %q", ContentTypeProblemJSON, got)
	}
}
