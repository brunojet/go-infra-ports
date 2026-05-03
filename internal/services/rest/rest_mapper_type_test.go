package rest

import (
	"net/http"
	"net/url"
	"testing"

	repcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
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

// Implement RestRequestSpec interface for testCreatePayload
func (t *testCreatePayload) New() repcts.RestRequestSpec      { return &testCreatePayload{} }
func (t *testCreatePayload) SetBody(_ repcts.RestRequestSpec) {}

// Implement RestRequestSpec interface for testUpdatePayload
func (t *testUpdatePayload) New() repcts.RestRequestSpec      { return &testUpdatePayload{} }
func (t *testUpdatePayload) SetBody(_ repcts.RestRequestSpec) {}

// Implement RestResponseSpec interface for testResponse
func (tr *testResponse) New() repcts.RestResponseSpec {
	return &testResponse{}
}

func (tr *testResponse) NewSlice(n int) []repcts.RestResponseSpec {
	specs := make([]repcts.RestResponseSpec, n)
	for i := range specs {
		specs[i] = &testResponse{}
	}
	return specs
}

func TestDefaultRestUpstreamMapper_ToUpstreamPost_CopiesPayload(t *testing.T) {
	m := &DefaultRestUpstreamMapper[*testCreatePayload, *testUpdatePayload]{}
	payload := &testCreatePayload{Name: "test"}
	ids := types.Identifiers{}
	var upsPayload svccts.RestRequestSpec

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
	var upsPayload svccts.RestRequestSpec

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
	var upsPayload svccts.RestRequestSpec

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
	if err != errRestDownstreamMapperResponseTypeAssertion {
		t.Fatalf("expected errRestDownstreamMapperResponseTypeAssertion, got %v", err)
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
	var upsPayload svccts.RestResponseSpec = &testResponse{ID: "1", Name: "info"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamInformation(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamRedirection_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 301
	var upsPayload svccts.RestResponseSpec = &testResponse{ID: "1", Name: "redirect"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamRedirection(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamProblem_ReturnsNil(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 400
	var upsPayload svccts.RestResponseSpec = &testResponse{ID: "1", Name: "problem"}
	serviceMeta := &ServiceMeta{}

	err := m.ToDownstreamProblem(statusCode, upsPayload, serviceMeta)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamProblem_With500Error(t *testing.T) {
	m := &DefaultRestDownstreamMapper[testResponse]{}
	statusCode := 500
	var upsPayload svccts.RestResponseSpec = &testResponse{ID: "1", Name: "server error"}
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
