package rest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"

	"github.com/brunojet/go-infra-ports/pkg/http_clients/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func newTestRepository(client HttpClient, registry RestRegistry, opts ...RepositoryOption) *restRepository {
	all := append([]RepositoryOption{
		WithHttpClient(client),
		WithRegistry(registry),
	}, opts...)
	o := newRepositoryOptions(all...)
	return &restRepository{client: o.client, registry: o.registry, opts: o}
}

func TestResolveURL_InterpolatesAndAppendsQuery(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{},
		WithBasePath("/api"),
		WithPath(MethodGet, "/users/{id}"),
	)

	got, err := r.resolveURL(MethodGet, types.Identifiers{"id": "42"}, url.Values{"v": {"1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, "/api/users/42") {
		t.Fatalf("unexpected URL: %s", got)
	}
	if !strings.Contains(got, "v=1") {
		t.Fatalf("expected query param v=1 in URL: %s", got)
	}
}

func TestResolveURL_MissingIdentifiers_ReturnsError(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{},
		WithBasePath("/api"),
		WithPath(MethodGet, "/users/{id}"),
	)
	_, err := r.resolveURL(MethodGet, types.Identifiers{}, nil)
	if err == nil {
		t.Fatal("expected error for missing required identifiers")
	}
	if !strings.Contains(err.Error(), "identifiers") {
		t.Fatalf("expected identifier-related error, got %q", err.Error())
	}
}

func TestResolveURL_UnconfiguredPathMethod_ReturnsError(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{},
		WithBasePath("/api"),
		WithPath(MethodGet, "/users/{id}"),
	)
	r.opts.paths = map[PathMethod]*pathEntry{}

	_, err := r.resolveURL(MethodCreate, nil, nil)
	if err == nil {
		t.Fatal("expected error for unconfigured path method")
	}
	if !errors.Is(err, errRepositoryPathMethodNotConfigured) {
		t.Fatalf("expected errRepositoryPathMethodNotConfigured, got %v", err)
	}
}

func TestResolveURL_InvalidBaseURL_ReturnsError(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})
	r.opts.basePath = "://bad url"

	_, err := r.resolveURL(MethodGet, nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid base URL")
	}
}

func TestBuildHTTPRequest_SetsMethodURLAndBody(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})

	req, err := r.buildHTTPRequest(context.Background(), http.MethodPost, "https://api.example.com/items", []byte(`{"x":1}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPost {
		t.Fatalf("expected POST, got %s", req.Method)
	}
	if req.URL.String() != "https://api.example.com/items" {
		t.Fatalf("unexpected URL: %s", req.URL)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type application/json")
	}
}

func TestExecuteRequest_NilResponse_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockHttpClient(ctrl)
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, nil)
	r := newTestRepository(mock, registryStub{})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	_, err = r.executeRequest(context.Background(), req) //nolint:bodyclose // closed inside readBody
	if err == nil {
		t.Fatal("expected error for nil response")
	}
	if !errors.Is(err, errRepositoryNilHTTPResponse) {
		t.Fatalf("expected errRepositoryNilHTTPResponse, got %v", err)
	}
}

func TestBuildHTTPRequest_MergesConfigAndRequestHeaders(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{}, WithHeader("X-Config", "cfg"))

	reqHeaders := http.Header{"X-Request": {"req"}, "X-Config": {"override"}}
	req, err := r.buildHTTPRequest(context.Background(), http.MethodGet, "https://api.example.com/", nil, reqHeaders)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Header.Get("X-Request") != "req" {
		t.Fatalf("expected X-Request=req")
	}
	if req.Header.Get("X-Config") != "override" {
		t.Fatalf("expected request-level header to override config-level header")
	}
}

func TestBuildHTTPRequest_ContentType_NotOverriddenWhenAlreadySet(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})

	reqHeaders := http.Header{"Content-Type": {"application/xml"}}
	req, err := r.buildHTTPRequest(context.Background(), http.MethodPost, "https://api.example.com/items", []byte(`<x/>`), reqHeaders)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Header.Get("Content-Type") != "application/xml" {
		t.Fatalf("expected caller Content-Type to be preserved, got %s", req.Header.Get("Content-Type"))
	}
}

func TestBuildHTTPRequest_ContentType_NotSetWhenNoBody(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})

	req, err := r.buildHTTPRequest(context.Background(), http.MethodGet, "https://api.example.com/items", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Header.Get("Content-Type") != "" {
		t.Fatalf("expected no Content-Type for bodyless request, got %s", req.Header.Get("Content-Type"))
	}
}

func TestReadBody_ReadsAndClosesBody(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(`{"ok":true}`)),
	}
	data, err := r.readBody(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", data)
	}
}

func TestMapResponse_SetsStatusCode(t *testing.T) {
	reg := NewRestRegistry()
	r := newTestRepository(newNoOpMockClient(t), reg)

	out := &RestResponse{}
	if err := r.mapResponse(http.StatusOK, []byte(`{}`), out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", out.Context.StatusCode)
	}
}

func TestMapResponses_SetsStatusCode(t *testing.T) {
	reg := NewRestRegistry()
	r := newTestRepository(newNoOpMockClient(t), reg)

	out := &RestResponses{}
	if err := r.mapResponses(http.StatusOK, []byte(`[{}]`), out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", out.Context.StatusCode)
	}
}
