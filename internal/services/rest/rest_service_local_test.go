package rest

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	repcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

// --- error mapper stubs for coverage of error paths ---

type queryErrUpstreamMapper struct {
	DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]
}

func (m *queryErrUpstreamMapper) ToUpstreamQuery(_, _ url.Values) error {
	return errors.New("upstream query error")
}

type headersErrUpstreamMapper struct {
	DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]
}

func (m *headersErrUpstreamMapper) ToUpstreamHeaders(_, _ http.Header) error {
	return errors.New("upstream headers error")
}

type statusCodeErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *statusCodeErrDownstreamMapper) ToDownstreamStatusCode(_ int, _ *int) error {
	return errors.New("downstream status error")
}

type headersErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *headersErrDownstreamMapper) ToDownstreamHeaders(_, _ http.Header) error {
	return errors.New("downstream headers error")
}

type responseErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *responseErrDownstreamMapper) ToDownstreamResponse(_ any, _ *testResponse) error {
	return errors.New("downstream response error")
}

type responseMetaErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *responseMetaErrDownstreamMapper) ToDownstreamResponseMeta(_ types.ResponseMeta, _ *ServiceMeta) error {
	return errors.New("downstream meta error")
}

type informationErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *informationErrDownstreamMapper) ToDownstreamInformation(_ int, _ RestResponseSpec, _ *ServiceMeta) error {
	return errors.New("downstream information error")
}

type redirectionErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *redirectionErrDownstreamMapper) ToDownstreamRedirection(_ int, _ RestResponseSpec, _ *ServiceMeta) error {
	return errors.New("downstream redirection error")
}

type problemErrDownstreamMapper struct {
	DefaultRestDownstreamMapper[testResponse]
}

func (m *problemErrDownstreamMapper) ToDownstreamProblem(_ int, _ RestResponseSpec, _ *ServiceMeta) error {
	return errors.New("downstream problem error")
}

type postErrUpstreamMapper struct {
	DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]
}

func (m *postErrUpstreamMapper) ToUpstreamPost(_ testCreatePayload, _ Identifiers, _ *repcts.RestRequestSpec) error {
	return errors.New("upstream post error")
}

type patchErrUpstreamMapper struct {
	DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]
}

func (m *patchErrUpstreamMapper) ToUpstreamPatch(_ testUpdatePayload, _ Identifiers, _ *repcts.RestRequestSpec) error {
	return errors.New("upstream patch error")
}

type putErrUpstreamMapper struct {
	DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]
}

func (m *putErrUpstreamMapper) ToUpstreamPut(_ testCreatePayload, _ Identifiers, _ *repcts.RestRequestSpec) error {
	return errors.New("upstream put error")
}

// errorRepository allows forcing repo errors on each operation.
type errorRepository struct {
	createErr error
	listErr   error
	getErr    error
	updateErr error
	saveErr   error
	deleteErr error
}

func (r *errorRepository) Create(_ context.Context, _ repcts.RestRequest, _ *repcts.RestResponse) error {
	return r.createErr
}

func (r *errorRepository) List(_ context.Context, _ types.RequestContext, _ *repcts.RestResponses) error {
	return r.listErr
}

func (r *errorRepository) Get(_ context.Context, _ types.RequestContext, _ *repcts.RestResponse) error {
	return r.getErr
}

func (r *errorRepository) Update(_ context.Context, _ repcts.RestRequest, _ *repcts.RestResponse) error {
	return r.updateErr
}

func (r *errorRepository) Save(_ context.Context, _ repcts.RestRequest, _ *repcts.RestResponse) error {
	return r.saveErr
}

func (r *errorRepository) Delete(_ context.Context, _ types.RequestContext, _ *repcts.RestResponse) error {
	return r.deleteErr
}

// statusMockRepository returns a response with the configured status code and optional data.
type statusMockRepository struct {
	statusCode int
	data       repcts.RestResponseSpec
	sliceData  []repcts.RestResponseSpec
}

func (m *statusMockRepository) Create(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	response.Data = m.data
	return nil
}

func (m *statusMockRepository) List(_ context.Context, _ types.RequestContext, response *repcts.RestResponses) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	response.Data = m.sliceData
	return nil
}

func (m *statusMockRepository) Get(_ context.Context, _ types.RequestContext, response *repcts.RestResponse) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	response.Data = m.data
	return nil
}

func (m *statusMockRepository) Update(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	response.Data = m.data
	return nil
}

func (m *statusMockRepository) Save(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	response.Data = m.data
	return nil
}

func (m *statusMockRepository) Delete(_ context.Context, _ types.RequestContext, response *repcts.RestResponse) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode}
	return nil
}

func newSvcForLocal() (svccts.Service[testCreatePayload, testResponse, testUpdatePayload], error) {
	return NewRestService[testCreatePayload, testResponse, testUpdatePayload](
		&statusMockRepository{},
	)
}

func TestMapRestResponse_WithNilServiceResponse_ReturnsError(t *testing.T) {
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](
		&statusMockRepository{statusCode: 200, data: &testResponse{ID: "1", Name: "ok"}},
	)

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{},
		Body:    testCreatePayload{},
	}

	err := svc.Create(context.Background(), req, nil)
	if err == nil {
		t.Fatal("expected error for nil response, got nil")
	}
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestMapRestResponse_2xx_MapsDataAndStatus(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "42", Name: "ok"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.Context.StatusCode)
	}
	if resp.Data.ID != "42" {
		t.Fatalf("expected ID 42, got %s", resp.Data.ID)
	}
}

func TestMapRestResponse_1xx_SetsStatusAndCallsInformation(t *testing.T) {
	repo := &statusMockRepository{statusCode: 100}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 100 {
		t.Fatalf("expected status 100, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponse_3xx_SetsStatusAndCallsRedirection(t *testing.T) {
	repo := &statusMockRepository{statusCode: 301}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 301 {
		t.Fatalf("expected status 301, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponse_4xx_SetsStatusAndCallsProblem(t *testing.T) {
	repo := &statusMockRepository{statusCode: 404}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponse_5xx_SetsStatusAndCallsProblem(t *testing.T) {
	repo := &statusMockRepository{statusCode: 500}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 500 {
		t.Fatalf("expected status 500, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponses_WithNilServiceResponse_ReturnsError(t *testing.T) {
	svc, _ := newSvcForLocal()

	err := svc.List(context.Background(), types.RequestContext{}, nil)
	if err == nil {
		t.Fatal("expected error for nil response, got nil")
	}
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestMapRestResponses_2xx_MapsSliceData(t *testing.T) {
	repo := &statusMockRepository{
		statusCode: 200,
		sliceData: []repcts.RestResponseSpec{
			&testResponse{ID: "1", Name: "first"},
			&testResponse{ID: "2", Name: "second"},
		},
	}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.Context.StatusCode)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != "1" || resp.Data[1].ID != "2" {
		t.Fatalf("unexpected data: %+v", resp.Data)
	}
}

func TestMapRestResponses_2xx_WithEmptySlice_MapsEmptyData(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, sliceData: []repcts.RestResponseSpec{}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Fatalf("expected empty data slice, got %d", len(resp.Data))
	}
}

func TestMapRestResponses_2xx_WithNilRestData_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, sliceData: nil}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err != errRestServiceNilResponseData {
		t.Fatalf("expected errRestServiceNilResponseData, got %v", err)
	}
}

func TestMapRestResponses_2xx_WithNilOutputSlicePointer_ReturnsError(t *testing.T) {
	svc := &restService[testCreatePayload, testResponse, testUpdatePayload]{
		downstream: &DefaultRestDownstreamMapper[testResponse]{},
	}

	err := svc.mapRestResponsesToServiceResponses(
		RestResponses{Context: types.ResponseContext{StatusCode: 200}, Data: []repcts.RestResponseSpec{}},
		nil,
		&ServiceMeta{},
	)
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestMapRestResponses_1xx_CallsInformation(t *testing.T) {
	repo := &statusMockRepository{statusCode: 100}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 100 {
		t.Fatalf("expected status 100, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponses_3xx_CallsRedirection(t *testing.T) {
	repo := &statusMockRepository{statusCode: 301}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 301 {
		t.Fatalf("expected status 301, got %d", resp.Context.StatusCode)
	}
}

func TestMapRestResponses_4xx_CallsProblem(t *testing.T) {
	repo := &statusMockRepository{statusCode: 404}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Context.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", resp.Context.StatusCode)
	}
}

// metaListRepository returns a 2xx RestResponses with Context.Meta populated.
type metaListRepository struct {
	statusMockRepository
	responseMeta types.ResponseMeta
}

func (m *metaListRepository) List(_ context.Context, _ types.RequestContext, response *repcts.RestResponses) error {
	response.Context = types.ResponseContext{StatusCode: m.statusCode, Meta: m.responseMeta}
	response.Data = m.sliceData
	return nil
}

// metaCapturingMapper embeds the default downstream mapper and records the meta passed to ToDownstreamResponseMeta.
type metaCapturingMapper struct {
	DefaultRestDownstreamMapper[testResponse]
	capturedMeta types.ResponseMeta
}

func (m *metaCapturingMapper) ToDownstreamResponseMeta(meta types.ResponseMeta, _ *ServiceMeta) error {
	m.capturedMeta = meta
	return nil
}

func TestMapRestResponses_2xx_CallsToDownstreamResponseMeta(t *testing.T) {
	wantMeta := types.ResponseMeta{"page": 2, "total": 50}
	repo := &metaListRepository{
		statusMockRepository: statusMockRepository{
			statusCode: 200,
			sliceData:  []repcts.RestResponseSpec{&testResponse{ID: "1", Name: "first"}},
		},
		responseMeta: wantMeta,
	}
	downstream := &metaCapturingMapper{}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](downstream))

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), types.RequestContext{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(downstream.capturedMeta) == 0 {
		t.Fatal("expected ToDownstreamResponseMeta to be called with meta, but capturedMeta is empty")
	}
	if downstream.capturedMeta["page"] != wantMeta["page"] || downstream.capturedMeta["total"] != wantMeta["total"] {
		t.Fatalf("unexpected captured meta: %v", downstream.capturedMeta)
	}
}

func TestMapUpstreamContext_QueryError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]

	err := svc.Create(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected query mapping error")
	}
}

func TestMapUpstreamContext_HeadersError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&headersErrUpstreamMapper{}))

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]

	err := svc.Create(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected headers mapping error")
	}
}

func TestMapDownstreamContext_StatusCodeError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream status mapping error")
	}
}

func TestMapDownstreamContext_HeadersError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&headersErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream headers mapping error")
	}
}

func TestMapRestResponse_2xx_WithNilData_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: nil}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err != errRestServiceNilResponseData {
		t.Fatalf("expected errRestServiceNilResponseData, got %v", err)
	}
}

func TestMapRestResponse_2xx_ToDownstreamResponseError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&responseErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream response mapping error")
	}
}

func TestMapRestResponse_InvalidNon2xxStatus_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 700}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err != errRestServiceInvalidNon2xxStatus {
		t.Fatalf("expected errRestServiceInvalidNon2xxStatus, got %v", err)
	}
}

func TestMapNon2xxResponse_MapperError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 100}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&informationErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream information mapping error")
	}
}

func TestMapNon2xxResponse_RedirectionMapperError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 301}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&redirectionErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream redirection mapping error")
	}
}

func TestMapNon2xxResponse_ProblemMapperError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 404}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&problemErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream problem mapping error")
	}
}

func TestMapRestResponses_2xx_ItemMapperError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, sliceData: []repcts.RestResponseSpec{&testResponse{ID: "1"}}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&responseErrDownstreamMapper{}))

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream slice item mapping error")
	}
}

func TestMapRestResponses_2xx_MetaMapperError_ReturnsWrappedError(t *testing.T) {
	repo := &metaListRepository{
		statusMockRepository: statusMockRepository{
			statusCode: 200,
			sliceData:  []repcts.RestResponseSpec{&testResponse{ID: "1"}},
		},
		responseMeta: types.ResponseMeta{"page": 1},
	}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&responseMetaErrDownstreamMapper{}))

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream response meta mapping error")
	}
}

func TestMapRestResponses_InvalidNon2xxStatus_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 700}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err != errRestServiceInvalidNon2xxStatus {
		t.Fatalf("expected errRestServiceInvalidNon2xxStatus, got %v", err)
	}
}

func TestMapRestResponses_Non2xxMapperError_ReturnsWrappedError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 100}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&informationErrDownstreamMapper{}))

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected downstream non-2xx mapper error for list")
	}
}

func TestMapNon2xxResponse_InvalidStatus_ReturnsError(t *testing.T) {
	svc := &restService[testCreatePayload, testResponse, testUpdatePayload]{
		downstream: &DefaultRestDownstreamMapper[testResponse]{},
	}

	err := svc.mapNon2xxResponse(700, nil, &ServiceMeta{})
	if err != errRestServiceInvalidNon2xxStatus {
		t.Fatalf("expected errRestServiceInvalidNon2xxStatus, got %v", err)
	}
}
