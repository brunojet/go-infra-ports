package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	rpocts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

// --- test doubles ---

type apiTestResponse struct{ ID string }

func (r *apiTestResponse) New() rpocts.RestDataSpec { return &apiTestResponse{} }
func (r *apiTestResponse) NewSlice(n int) []rpocts.RestDataSpec {
	out := make([]rpocts.RestDataSpec, n)
	for i := range out {
		out[i] = &apiTestResponse{}
	}
	return out
}
func (r *apiTestResponse) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*struct{ ID string })(r))
}
func (r *apiTestResponse) MarshalJSON() ([]byte, error) { return json.Marshal(*r) }
func (r *apiTestResponse) SetBody(body any) error {
	if body == nil {
		*r = apiTestResponse{}
		return nil
	}
	switch v := body.(type) {
	case *apiTestResponse:
		*r = *v
		return nil
	case apiTestResponse:
		*r = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (r *apiTestResponse) Body() any { return *r }

type apiTestCreate struct{ Name string }

func (c *apiTestCreate) New() rpocts.RestDataSpec { return &apiTestCreate{} }
func (c *apiTestCreate) NewSlice(n int) []rpocts.RestDataSpec {
	out := make([]rpocts.RestDataSpec, n)
	for i := range out {
		out[i] = &apiTestCreate{}
	}
	return out
}
func (c *apiTestCreate) SetBody(body any) error {
	if body == nil {
		*c = apiTestCreate{}
		return nil
	}
	switch v := body.(type) {
	case *apiTestCreate:
		*c = *v
		return nil
	case apiTestCreate:
		*c = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (c *apiTestCreate) Body() any                       { return *c }
func (c *apiTestCreate) MarshalJSON() ([]byte, error)    { return json.Marshal(*c) }
func (c *apiTestCreate) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, c) }

type apiTestUpdate struct{ Name string }

func (u *apiTestUpdate) New() rpocts.RestDataSpec { return &apiTestUpdate{} }
func (u *apiTestUpdate) NewSlice(n int) []rpocts.RestDataSpec {
	out := make([]rpocts.RestDataSpec, n)
	for i := range out {
		out[i] = &apiTestUpdate{}
	}
	return out
}
func (u *apiTestUpdate) SetBody(body any) error {
	if body == nil {
		*u = apiTestUpdate{}
		return nil
	}
	switch v := body.(type) {
	case *apiTestUpdate:
		*u = *v
		return nil
	case apiTestUpdate:
		*u = v
		return nil
	default:
		return fmt.Errorf("unsupported body type: %T", body)
	}
}
func (u *apiTestUpdate) Body() any                       { return *u }
func (u *apiTestUpdate) MarshalJSON() ([]byte, error)    { return json.Marshal(*u) }
func (u *apiTestUpdate) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, u) }

// apiStubRepo is a minimal RestRepository that returns a configured 2xx response.
type apiStubRepo struct {
	statusCode int
	data       rpocts.RestDataSpec
}

func (r *apiStubRepo) Create(_ context.Context, _ rpocts.RestRequest, resp *rpocts.RestResponse) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	resp.Data = r.data
	return nil
}

func (r *apiStubRepo) List(_ context.Context, _ types.RequestContext, resp *rpocts.RestResponses) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	if r.data != nil {
		resp.Data = []rpocts.RestDataSpec{r.data}
	}
	return nil
}

func (r *apiStubRepo) Get(_ context.Context, _ types.RequestContext, resp *rpocts.RestResponse) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	resp.Data = r.data
	return nil
}

func (r *apiStubRepo) Update(_ context.Context, _ rpocts.RestRequest, resp *rpocts.RestResponse) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	resp.Data = r.data
	return nil
}

func (r *apiStubRepo) Save(_ context.Context, _ rpocts.RestRequest, resp *rpocts.RestResponse) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	resp.Data = r.data
	return nil
}

func (r *apiStubRepo) Delete(_ context.Context, _ types.RequestContext, resp *rpocts.RestResponse) error {
	resp.Context = types.ResponseContext{StatusCode: r.statusCode}
	return nil
}

func (r *apiStubRepo) NewRequest(method rpocts.RestMethod) (*rpocts.RestRequest, error) {
	req := &rpocts.RestRequest{Context: types.RequestContext{}}
	switch method {
	case rpocts.MethodCreate, rpocts.MethodSave:
		req.Data = &apiTestCreate{}
	case rpocts.MethodUpdate:
		req.Data = &apiTestUpdate{}
	default:
		req.Data = &apiTestCreate{}
	}
	return req, nil
}

func (r *apiStubRepo) NewResponse() *rpocts.RestResponse {
	return &rpocts.RestResponse{
		Context:     types.ResponseContext{},
		Information: &apiTestResponse{},
		Redirection: &apiTestResponse{},
		Problem:     &apiTestResponse{},
		Data:        &apiTestResponse{},
	}
}

// --- NewRestService ---

func TestNewRestService_NilRepo_ReturnsError(t *testing.T) {
	_, err := NewRestService[apiTestCreate, apiTestResponse, apiTestUpdate](nil)
	if err == nil {
		t.Fatal("expected error for nil repository")
	}
}

func TestNewRestService_ValidRepo_ReturnsService(t *testing.T) {
	repo := &apiStubRepo{statusCode: 200, data: &apiTestResponse{ID: "1"}}
	svc, err := NewRestService[apiTestCreate, apiTestResponse, apiTestUpdate](repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

// --- WithUpstreamMapper ---

func TestWithUpstreamMapper_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream mapper")
		}
	}()
	WithUpstreamMapper[apiTestCreate, apiTestUpdate](nil)
}

func TestWithUpstreamMapper_AppliedToService(t *testing.T) {
	repo := &apiStubRepo{statusCode: 200, data: &apiTestResponse{ID: "42"}}
	custom := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	opt := WithUpstreamMapper[apiTestCreate, apiTestUpdate](custom)

	svc, err := NewRestService[apiTestCreate, apiTestResponse, apiTestUpdate](repo, opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqCtx := types.RequestContext{Query: url.Values{}, Headers: http.Header{}}
	var resp svccts.ServiceResponse[apiTestResponse]
	if err := svc.Get(context.Background(), reqCtx, &resp); err != nil {
		t.Fatalf("unexpected error calling Get: %v", err)
	}
	if resp.Data.ID != "42" {
		t.Fatalf("expected ID 42, got %s", resp.Data.ID)
	}
}

// --- WithDownstreamMapper ---

func TestWithDownstreamMapper_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil downstream mapper")
		}
	}()
	WithDownstreamMapper[apiTestResponse](nil)
}

func TestWithDownstreamMapper_AppliedToService(t *testing.T) {
	repo := &apiStubRepo{statusCode: 200, data: &apiTestResponse{ID: "7"}}
	custom := &DefaultRestDownstreamMapper[apiTestResponse]{}
	opt := WithDownstreamMapper[apiTestResponse](custom)

	svc, err := NewRestService[apiTestCreate, apiTestResponse, apiTestUpdate](repo, opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqCtx := types.RequestContext{Query: url.Values{}, Headers: http.Header{}}
	var resp svccts.ServiceResponse[apiTestResponse]
	if err := svc.Get(context.Background(), reqCtx, &resp); err != nil {
		t.Fatalf("unexpected error calling Get: %v", err)
	}
	if resp.Data.ID != "7" {
		t.Fatalf("expected ID 7, got %s", resp.Data.ID)
	}
}

// --- DefaultRestUpstreamMapper ---

func TestDefaultRestUpstreamMapper_ToUpstreamPost(t *testing.T) {
	m := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	payload := apiTestCreate{Name: "x"}
	var spec any
	if err := m.ToUpstreamPost(payload, nil, &spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamPut(t *testing.T) {
	m := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	payload := apiTestCreate{Name: "y"}
	var spec any
	if err := m.ToUpstreamPut(payload, nil, &spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamPatch(t *testing.T) {
	m := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	payload := apiTestUpdate{Name: "z"}
	var spec any
	if err := m.ToUpstreamPatch(payload, nil, &spec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamQuery(t *testing.T) {
	m := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	src := url.Values{"k": {"v"}}
	dst := url.Values{}
	if err := m.ToUpstreamQuery(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Get("k") != "v" {
		t.Fatal("expected query values to be copied")
	}
}

func TestDefaultRestUpstreamMapper_ToUpstreamHeaders(t *testing.T) {
	m := &DefaultRestUpstreamMapper[apiTestCreate, apiTestUpdate]{}
	src := http.Header{"X-Test": {"1"}}
	dst := http.Header{}
	if err := m.ToUpstreamHeaders(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Get("X-Test") != "1" {
		t.Fatal("expected headers to be copied")
	}
}

// --- DefaultRestDownstreamMapper ---

func TestDefaultRestDownstreamMapper_ToDownstreamStatusCode(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	var code int
	if err := m.ToDownstreamStatusCode(200, &code); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamHeaders(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	src := http.Header{"X-Out": {"y"}}
	dst := http.Header{}
	if err := m.ToDownstreamHeaders(src, dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Get("X-Out") != "y" {
		t.Fatal("expected headers to be copied")
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponse(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	src := apiTestResponse{ID: "99"}
	var dst apiTestResponse
	if err := m.ToDownstreamResponse(src, &dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.ID != "99" {
		t.Fatalf("expected ID 99, got %s", dst.ID)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamResponseMeta(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	if err := m.ToDownstreamResponseMeta(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamInformation(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	if err := m.ToDownstreamInformation(100, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamRedirection(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	if err := m.ToDownstreamRedirection(301, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultRestDownstreamMapper_ToDownstreamProblem(t *testing.T) {
	m := &DefaultRestDownstreamMapper[apiTestResponse]{}
	if err := m.ToDownstreamProblem(404, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
