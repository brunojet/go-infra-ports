package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gomock "go.uber.org/mock/gomock"

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
	r.opts.paths = map[RestMethod]*pathEntry{}

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

type localTestResponse struct {
	Name string `json:"name"`
}

func TestMapResponse_DefaultRawBody(t *testing.T) {
	reg := NewRestRegistry(WithResponseOf[DefaultRestResponse](http.StatusOK))
	r := newTestRepository(newNoOpMockClient(t), reg)

	out := &RestResponse{}
	payload := []byte(`{"hello":"world"}`)
	if err := r.mapResponse(http.StatusOK, payload, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, out.Context.StatusCode)
	}
	if rb, ok := out.Data.Body().(DefaultRestResponse); ok {
		if !bytes.Equal([]byte(rb), payload) {
			t.Fatalf("unexpected raw body: got %s, want %s", []byte(rb), payload)
		}
	} else {
		t.Fatalf("unexpected data body type: %T", out.Data.Body())
	}
}

func TestMapResponses_DefaultRawBodySlice(t *testing.T) {
	reg := NewRestRegistry(WithResponseOf[DefaultRestResponse](http.StatusOK))
	r := newTestRepository(newNoOpMockClient(t), reg)

	out := &RestResponses{}
	payload := []byte(`[{"name":"a"},{"name":"b"}]`)
	if err := r.mapResponses(http.StatusOK, payload, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, out.Context.StatusCode)
	}
	if len(out.Data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out.Data))
	}
	if rb, ok := out.Data[0].Body().(DefaultRestResponse); ok {
		if !bytes.Contains([]byte(rb), []byte(`"name":"a"`)) {
			t.Fatalf("unexpected first item body: %s", []byte(rb))
		}
	} else {
		t.Fatalf("unexpected data body type: %T", out.Data[0].Body())
	}
}

func TestMapResponse_WithEnvelope(t *testing.T) {
	reg := NewRestRegistry(
		WithResponseOf[localTestResponse](http.StatusOK),
		WithResponseEnvelope("data", "meta", http.StatusOK),
	)
	r := newTestRepository(newNoOpMockClient(t), reg)

	out := &RestResponse{}
	payload := []byte(`{"data":{"name":"joe"},"meta":{"page":2}}`)
	if err := r.mapResponse(http.StatusOK, payload, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, out.Context.StatusCode)
	}
	if page, ok := out.Context.Meta["page"]; !ok || page.(float64) != 2 {
		t.Fatalf("unexpected meta page: %#v", out.Context.Meta)
	}
	if data, ok := out.Data.Body().(localTestResponse); ok {
		if data.Name != "joe" {
			t.Fatalf("unexpected name: %s", data.Name)
		}
	} else {
		t.Fatalf("unexpected data type: %T", out.Data.Body())
	}
}

// Tests for executeRequest / executeBodyRequest / executeNoBodyRequest error and success paths.

func TestExecuteRequest_DoError_ReturnsWrappedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockHttpClient(ctrl)
	orig := errors.New("network fail")
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, orig)

	r := newTestRepository(mock, registryStub{})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	resp, err := r.executeRequest(context.Background(), req) //nolint:bodyclose // test reads body when present
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected error from executeRequest when client.Do returns an error")
	}
	if !errors.Is(err, orig) {
		t.Fatalf("expected wrapped original error, got %v", err)
	}
}

func TestExecuteBodyRequest_BodyNil_ReturnsError(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})
	resp, err := r.executeBodyRequest(context.Background(), MethodCreate, http.MethodPost, RestRequest{
		Context: types.RequestContext{},
		Data:    nil,
	}) //nolint:bodyclose // expected error path, but close if present
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected error for nil request body")
	}
	if !errors.Is(err, errRepositoryRequestBodyNil) {
		t.Fatalf("expected errRepositoryRequestBodyNil, got %v", err)
	}
}

type failingResolveRegistry struct{ Err error }

func (f failingResolveRegistry) Merge(other RestRegistry) RestRegistry                { return other }
func (f failingResolveRegistry) ResolveRequest(_ RestDataSpec, _ *[]byte) error       { return f.Err }
func (f failingResolveRegistry) ResolveEnvelopeRequest(_ RestMethod, _ *[]byte) error { return nil }
func (f failingResolveRegistry) ResolveResponse(int, []byte, *RestDataSpec) error     { return nil }
func (f failingResolveRegistry) ResolveResponses(int, []byte, *[]RestDataSpec) error  { return nil }
func (f failingResolveRegistry) ResolveEnvelopeResponse(int, *[]byte, *types.ResponseMeta) error {
	return nil
}
func (f failingResolveRegistry) NewRequestSpec(_ RestMethod) (RestDataSpec, error) { return nil, nil }
func (f failingResolveRegistry) ReleaseRequestSpec(RestDataSpec)                   {}

func TestExecuteBodyRequest_ResolveRequestError_ReturnsWrappedError(t *testing.T) {
	orig := errors.New("marshal fail")
	r := newTestRepository(newNoOpMockClient(t), failingResolveRegistry{Err: orig})

	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{"x":1}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}

	resp, err := r.executeBodyRequest(context.Background(), MethodCreate, http.MethodPost, RestRequest{Data: reqBody}) //nolint:bodyclose // test reads body when present
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected error from ResolveRequest")
	}
	if !errors.Is(err, orig) {
		t.Fatalf("expected original resolve error wrapped, got %v", err)
	}
}

type failingEnvelopeRegistry struct{ Err error }

func (f failingEnvelopeRegistry) Merge(other RestRegistry) RestRegistry { return other }
func (f failingEnvelopeRegistry) ResolveRequest(_ RestDataSpec, out *[]byte) error {
	*out = []byte(`{"ok":true}`)
	return nil
}
func (f failingEnvelopeRegistry) ResolveEnvelopeRequest(_ RestMethod, _ *[]byte) error { return f.Err }
func (f failingEnvelopeRegistry) ResolveResponse(int, []byte, *RestDataSpec) error     { return nil }
func (f failingEnvelopeRegistry) ResolveResponses(int, []byte, *[]RestDataSpec) error  { return nil }
func (f failingEnvelopeRegistry) ResolveEnvelopeResponse(int, *[]byte, *types.ResponseMeta) error {
	return nil
}
func (f failingEnvelopeRegistry) NewRequestSpec(_ RestMethod) (RestDataSpec, error) { return nil, nil }
func (f failingEnvelopeRegistry) ReleaseRequestSpec(RestDataSpec)                   {}

func TestExecuteBodyRequest_ResolveEnvelopeRequestError_ReturnsWrappedError(t *testing.T) {
	orig := errors.New("envelope fail")
	r := newTestRepository(newNoOpMockClient(t), failingEnvelopeRegistry{Err: orig})

	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{"x":1}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}

	resp, err := r.executeBodyRequest(context.Background(), MethodCreate, http.MethodPost, RestRequest{Data: reqBody}) //nolint:bodyclose // test reads body when present
	if resp != nil {
		_ = resp.Body.Close()
	}
	if err == nil {
		t.Fatal("expected error from ResolveEnvelopeRequest")
	}
	if !errors.Is(err, orig) {
		t.Fatalf("expected original envelope error wrapped, got %v", err)
	}
}

func TestExecuteBodyRequest_Success_ReturnsResponse(t *testing.T) {
	mock, _ := newMockClient(t, http.StatusOK, `{"ok":true}`)
	repoIface := newTestRepo(t, mock)
	r, ok := repoIface.(*restRepository)
	if !ok {
		t.Fatal("expected *restRepository from newTestRepo")
	}

	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{"name":"x"}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}

	resp, err := r.executeBodyRequest(context.Background(), MethodCreate, http.MethodPost, RestRequest{Data: reqBody}) //nolint:bodyclose // test reads body when present
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected response: %#v", resp)
	}
	// consume body to ensure no leaks in helpers
	_, _ = io.ReadAll(resp.Body)
}

func TestExecuteNoBodyRequest_Success_ReturnsResponse(t *testing.T) {
	mock, _ := newMockClient(t, http.StatusOK, `{"ok":true}`)
	repoIface := newTestRepo(t, mock)
	r, ok := repoIface.(*restRepository)
	if !ok {
		t.Fatal("expected *restRepository from newTestRepo")
	}

	resp, err := r.executeNoBodyRequest(context.Background(), MethodList, http.MethodGet, types.RequestContext{}) //nolint:bodyclose // test reads body when present
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected response: %#v", resp)
	}
	_ = resp.Body.Close()
}

func TestExecuteNoBodyRequest_UnconfiguredPath_ReturnsError(t *testing.T) {
	r := newTestRepository(newNoOpMockClient(t), registryStub{})
	r.opts.paths = map[RestMethod]*pathEntry{}

	_, err := r.executeNoBodyRequest(context.Background(), MethodGet, http.MethodGet, types.RequestContext{}) //nolint:bodyclose // expected error path
	if err == nil {
		t.Fatal("expected error for missing path method")
	}
	if !errors.Is(err, errRepositoryPathMethodNotConfigured) {
		t.Fatalf("expected errRepositoryPathMethodNotConfigured, got %v", err)
	}
}
