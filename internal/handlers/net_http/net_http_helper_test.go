package net_http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/brunojet/go-infra-ports/internal/dtos"

	"errors"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	mocksvc "github.com/brunojet/go-infra-ports/pkg/services/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func TestBuildRequestContext_NilCtxWrites500(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/users/1", nil)
	rr := httptest.NewRecorder()
	entry := &routeEntry{params: []paramFormat{{name: "userId"}}}
	ok := buildRequestContext(rr, req, entry, nil)
	if ok {
		t.Fatalf("expected false when ctx is nil")
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected problem content type, got %s", ct)
	}
	var p DefaultProblemDetails
	_ = json.Unmarshal(rr.Body.Bytes(), &p)
	if p.Detail == "" {
		t.Fatalf("expected detail body")
	}
}

func TestParamValidationRegexInvalid(t *testing.T) {
	// Use a minimal mock for service; handler will validate params before calling service.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)

	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithInstance(AllInstanceMethods, "users/{userId}", rx))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/abc", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid param, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected problem content type, got %s", ct)
	}
	var p DefaultProblemDetails
	_ = json.Unmarshal(rr.Body.Bytes(), &p)
	if p.Detail == "" {
		t.Fatalf("expected detail body")
	}
}

func TestStatusOrAndWriteHelpers(t *testing.T) {
	if statusOr(0, 42) != 42 {
		t.Fatalf("expected fallback 42")
	}
	if statusOr(10, 42) != 10 {
		t.Fatalf("expected code 10")
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/1", nil)
	writeBadRequest(rr, req, "reason")
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400")
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected problem content type, got %s", ct)
	}
	var p DefaultProblemDetails
	_ = json.Unmarshal(rr.Body.Bytes(), &p)
	if p.Detail == "" {
		t.Fatalf("expected error msg")
	}
	if p.Instance != "/users/1" {
		t.Fatalf("expected instance /users/1, got %s", p.Instance)
	}
	rr2 := httptest.NewRecorder()
	writeServiceError(rr2, req, errors.New("boom"))
	if rr2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500")
	}
	if ct := rr2.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected problem content type, got %s", ct)
	}
}

func TestBuildRequestBody_Success(t *testing.T) {
	type p struct{ Name string }
	body := p{Name: "ok"}
	b, _ := json.Marshal(body)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/x", bytes.NewReader(b))
	var got p
	if !buildRequestBody(rr, req, &got) {
		t.Fatalf("expected success")
	}
	if got.Name != "ok" {
		t.Fatalf("unexpected body")
	}
}

func TestValidateParams_NoParams(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/x", nil)
	if !validateParams(rr, req, nil, &types.RequestContext{}) {
		t.Fatalf("expected true when entry is nil")
	}
}

func TestTitleFromStatus_UnknownStatus(t *testing.T) {
	if got := titleFromStatus(999); got != "HTTP Error" {
		t.Fatalf("expected 'HTTP Error' for unknown status, got %q", got)
	}
}

func TestRequestPath_NilRequest(t *testing.T) {
	if got := requestPath(nil); got != "" {
		t.Fatalf("expected empty string for nil request, got %q", got)
	}
}

func TestRequestPath_NilURL(t *testing.T) {
	req := &http.Request{URL: nil}
	if got := requestPath(req); got != "" {
		t.Fatalf("expected empty string for nil URL, got %q", got)
	}
}

func TestFirstNonEmpty_AllEmpty(t *testing.T) {
	if got := firstNonEmpty("", "", ""); got != "" {
		t.Fatalf("expected empty string when all values empty, got %q", got)
	}
}

func TestMetaDetailsString_NilMetadata(t *testing.T) {
	svcMeta := svccts.ServiceMeta{}
	if got := metaDetailsString(svcMeta, "key"); got != "" {
		t.Fatalf("expected empty string for nil metadata, got %q", got)
	}
}

func TestMetaDetailsString_MissingKey(t *testing.T) {
	svcMeta := svccts.ServiceMeta{Metadata: map[string]any{"other": "val"}}
	if got := metaDetailsString(svcMeta, "missing"); got != "" {
		t.Fatalf("expected empty string for missing key, got %q", got)
	}
}

func TestMetaDetailsString_NonStringValue(t *testing.T) {
	svcMeta := svccts.ServiceMeta{Metadata: map[string]any{"key": 42}}
	if got := metaDetailsString(svcMeta, "key"); got != "" {
		t.Fatalf("expected empty string for non-string value, got %q", got)
	}
}

func TestProblemFromMeta_WithProblemTypeAndTitle(t *testing.T) {
	svcMeta := svccts.ServiceMeta{
		Message: "something went wrong",
		Code:    "ERR_001",
		Metadata: map[string]any{
			"problem_type":  "urn:problem:custom",
			"problem_title": "Custom Title",
			"instance":      "/resources/1",
		},
	}
	p := problemFromMeta(422, "/resources/1", svcMeta)
	if p.Type != "urn:problem:custom" {
		t.Fatalf("expected custom problem type, got %q", p.Type)
	}
	if p.Title != "Custom Title" {
		t.Fatalf("expected custom title, got %q", p.Title)
	}
	if p.Extensions == nil {
		t.Fatalf("expected extensions to be set")
	}
	if p.Extensions["code"] != "ERR_001" {
		t.Fatalf("expected code extension, got %#v", p.Extensions["code"])
	}
	if _, ok := p.Extensions["details"]; !ok {
		t.Fatalf("expected details extension")
	}
}
