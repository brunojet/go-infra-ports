package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	repcts "github.com/brunojet/go-infra-ports/pkg/repositories"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type testCreatePayload struct {
	Name string
}

type testUpdatePayload struct {
	Name string
}

type testResponse struct {
	ID   string
	Name string
}

// Implement RestDataSpec (RestRequestSpec/RestResponseSpec) interface for testCreatePayload
func (t *testCreatePayload) New() repcts.RestRequestSpec { return &testCreatePayload{} }
func (t *testCreatePayload) NewSlice(n int) []repcts.RestRequestSpec {
	specs := make([]repcts.RestRequestSpec, n)
	for i := range specs {
		specs[i] = &testCreatePayload{}
	}
	return specs
}
func (t *testCreatePayload) SetBody(body any) error {
	if body == nil {
		*t = testCreatePayload{}
		return nil
	}
	switch v := body.(type) {
	case *testCreatePayload:
		*t = *v
		return nil
	case testCreatePayload:
		*t = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (t *testCreatePayload) Body() any                       { return *t }
func (t *testCreatePayload) MarshalJSON() ([]byte, error)    { return json.Marshal(*t) }
func (t *testCreatePayload) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, t) }

// Implement RestDataSpec (RestRequestSpec) interface for testUpdatePayload
func (t *testUpdatePayload) New() repcts.RestRequestSpec { return &testUpdatePayload{} }
func (t *testUpdatePayload) NewSlice(n int) []repcts.RestRequestSpec {
	specs := make([]repcts.RestRequestSpec, n)
	for i := range specs {
		specs[i] = &testUpdatePayload{}
	}
	return specs
}
func (t *testUpdatePayload) SetBody(body any) error {
	if body == nil {
		*t = testUpdatePayload{}
		return nil
	}
	switch v := body.(type) {
	case *testUpdatePayload:
		*t = *v
		return nil
	case testUpdatePayload:
		*t = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (t *testUpdatePayload) Body() any                       { return *t }
func (t *testUpdatePayload) MarshalJSON() ([]byte, error)    { return json.Marshal(*t) }
func (t *testUpdatePayload) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, t) }

// Implement RestDataSpec (RestResponseSpec) interface for testResponse
func (tr *testResponse) New() repcts.RestResponseSpec { return &testResponse{} }
func (tr *testResponse) NewSlice(n int) []repcts.RestResponseSpec {
	specs := make([]repcts.RestResponseSpec, n)
	for i := range specs {
		specs[i] = &testResponse{}
	}
	return specs
}
func (tr *testResponse) SetBody(body any) error {
	if body == nil {
		*tr = testResponse{}
		return nil
	}
	switch v := body.(type) {
	case *testResponse:
		*tr = *v
		return nil
	case testResponse:
		*tr = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (tr *testResponse) Body() any                       { return *tr }
func (tr *testResponse) MarshalJSON() ([]byte, error)    { return json.Marshal(*tr) }
func (tr *testResponse) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, tr) }

func TestDefaultRestUpstreamMapper_ToUpstreamPost_CopiesPayload(t *testing.T) {
	m := &DefaultRestUpstreamMapper[*testCreatePayload, *testUpdatePayload]{}
	payload := &testCreatePayload{Name: "test"}
	ids := types.Identifiers{}
	var upsPayload any

	err := m.ToUpstreamPost(payload, ids, &upsPayload)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if casted, ok := upsPayload.(*testCreatePayload); !ok || casted.Name != "test" {
		t.Fatalf("expected payload to be copied correctly, got %v", upsPayload)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamPut_CopiesUpdatePayload(t *testing.T) {
	m := &DefaultRestUpstreamMapper[*testCreatePayload, *testUpdatePayload]{}
	payload := &testCreatePayload{Name: "updated"}
	ids := types.Identifiers{}
	var upsPayload any

	err := m.ToUpstreamPut(payload, ids, &upsPayload)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if casted, ok := upsPayload.(*testCreatePayload); !ok || casted.Name != "updated" {
		t.Fatalf("expected update payload to be copied correctly, got %v", upsPayload)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamPatch_CopiesPatchPayload(t *testing.T) {
	m := &DefaultRestUpstreamMapper[*testCreatePayload, *testUpdatePayload]{}
	payload := &testUpdatePayload{Name: "patched"}
	ids := types.Identifiers{}
	var upsPayload any

	err := m.ToUpstreamPatch(payload, ids, &upsPayload)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if casted, ok := upsPayload.(*testUpdatePayload); !ok || casted.Name != "patched" {
		t.Fatalf("expected patch payload to be copied correctly, got %v", upsPayload)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamQuery_CopiesQueryValues(t *testing.T) {
	m := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	reqQuery := url.Values{"page": []string{"1"}, "limit": []string{"10"}}
	upsQuery := url.Values{}

	err := m.ToUpstreamQuery(reqQuery, upsQuery)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if upsQuery.Get("page") != "1" || upsQuery.Get("limit") != "10" {
		t.Fatalf("expected query values to be copied correctly, got %v", upsQuery)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamQuery_WithEmptyValues(t *testing.T) {
	m := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	reqQuery := url.Values{}
	upsQuery := url.Values{}

	err := m.ToUpstreamQuery(reqQuery, upsQuery)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(upsQuery) != 0 {
		t.Fatalf("expected empty query values, got %v", upsQuery)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamHeaders_CopiesHeaders(t *testing.T) {
	m := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	reqHeader := http.Header{"Content-Type": []string{"application/json"}, "X-Custom": []string{"value"}}
	upsHeader := http.Header{}

	err := m.ToUpstreamHeaders(reqHeader, upsHeader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if upsHeader.Get("Content-Type") != "application/json" || upsHeader.Get("X-Custom") != "value" {
		t.Fatalf("expected headers to be copied correctly, got %v", upsHeader)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamHeaders_WithEmptyHeaders(t *testing.T) {
	m := &DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}
	reqHeader := http.Header{}
	upsHeader := http.Header{}

	err := m.ToUpstreamHeaders(reqHeader, upsHeader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(upsHeader) != 0 {
		t.Fatalf("expected empty headers, got %v", upsHeader)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamStatusCode_CopiesStatusCode(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 200
	var downstreamStatusCode int

	err := m.ToDownstreamStatusCode(statusCode, &downstreamStatusCode)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if downstreamStatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", downstreamStatusCode)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamStatusCode_WithVariousStatuses(t *testing.T) {
	testCases := []int{201, 204, 301, 400, 404, 500, 503}
	m := &DefaultRestDownstreamMapper[testResponse]{}

	for _, status := range testCases {
		var downstreamStatusCode int
		err := m.ToDownstreamStatusCode(status, &downstreamStatusCode)

		if err != nil {
			t.Fatalf("expected no error for status %d, got %v", status, err)
		}
		if downstreamStatusCode != status {
			t.Fatalf("expected status code %d, got %d", status, downstreamStatusCode)
		}
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamHeaders_CopiesHeaders(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	upsHeader := http.Header{"X-Custom": []string{"value"}, "Content-Length": []string{"100"}}
	downstreamHeader := http.Header{}

	err := m.ToDownstreamHeaders(upsHeader, downstreamHeader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if downstreamHeader.Get("X-Custom") != "value" || downstreamHeader.Get("Content-Length") != "100" {
		t.Fatalf("expected headers to be copied correctly, got %v", downstreamHeader)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponse_SuccessfulCast(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	var upsPayload any = testResponse{ID: "1", Name: "test"}
	var payload testResponse

	err := m.ToDownstreamResponse(upsPayload, &payload)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if payload.ID != "1" || payload.Name != "test" {
		t.Fatalf("expected response to be cast correctly, got %v", payload)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponse_TypeAssertionFailure(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	var upsPayload any = "invalid string payload"
	var payload testResponse

	err := m.ToDownstreamResponse(upsPayload, &payload)

	if err == nil {
		t.Fatalf("expected type assertion error, got nil")
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponse_WithWrongType(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	var upsPayload any = 42
	var payload testResponse

	err := m.ToDownstreamResponse(upsPayload, &payload)

	if err == nil {
		t.Fatalf("expected type assertion error for int payload, got nil")
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamInformation_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 100
	var upsPayload any = &testResponse{ID: "1", Name: "info"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamInformation(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamRedirection_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 301
	var upsPayload any = &testResponse{ID: "1", Name: "redirect"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamRedirection(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamProblem_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 400
	var upsPayload any = &testResponse{ID: "1", Name: "problem"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamProblem(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamProblem_With500Error(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 500
	var upsPayload any = &testResponse{ID: "1", Name: "server error"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamProblem(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error for 500 status, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponseMeta_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	meta := svccts.ResponseMeta{"page": 1, "total": 100}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamResponseMeta(meta, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
