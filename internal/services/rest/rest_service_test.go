package rest

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	repcts "github.com/brunojet/go-infra-ports/pkg/repositories"
	repmocks "github.com/brunojet/go-infra-ports/pkg/repositories/rest/mocks"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"

	"go.uber.org/mock/gomock"
)

func TestNewRestService_WithNilRepository_ReturnsError(t *testing.T) {
	_, err := NewRestService[testCreatePayload, testResponse, testUpdatePayload](nil)

	if err == nil {
		t.Fatal("expected error for nil repository, got nil")
	}
	if err != errRestServiceRepositoryNil {
		t.Fatalf("expected errRestServiceRepositoryNil, got %v", err)
	}
}

func TestNewRestService_WithNilMappers_UsesDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)

	svc, err := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	if err != nil {
		t.Fatalf("expected no error with nil mappers, got %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestNewRestService_WithProvidedMappers_UsesProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)

	svc, err := NewRestService[testCreatePayload, testResponse, testUpdatePayload](
		repo,
		WithUpstreamMapper[testCreatePayload, testUpdatePayload](&DefaultRestUpstreamMapper[testCreatePayload, testUpdatePayload]{}),
		WithDownstreamMapper[testResponse](&DefaultRestDownstreamMapper[testResponse]{}),
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestRestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
			response.Context = types.ResponseContext{StatusCode: 201}
			response.Data = &testResponse{ID: "1", Name: "created"}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{
			Query:       url.Values{},
			Headers:     http.Header{},
			Identifiers: types.Identifiers{},
		},
		Body: testCreatePayload{Name: "test"},
	}

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Create(context.Background(), req, &resp); err != nil {
		t.Fatalf("expected no error in Create, got %v", err)
	}
	if resp.Context.StatusCode != 201 {
		t.Fatalf("expected status code 201, got %d", resp.Context.StatusCode)
	}
}

func TestRestService_Get_NilResponse_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	err := svc.Get(context.Background(), types.RequestContext{}, nil)
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestRestService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ types.RequestContext, response *repcts.RestResponses) error {
			response.Context = types.ResponseContext{StatusCode: 200}
			response.Data = []repcts.RestResponseSpec{
				&testResponse{ID: "1", Name: "item1"},
				&testResponse{ID: "2", Name: "item2"},
			}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	reqCtx := types.RequestContext{
		Query:       url.Values{},
		Headers:     http.Header{},
		Identifiers: types.Identifiers{},
	}

	var resp svccts.ServiceResponses[testResponse]
	if err := svc.List(context.Background(), reqCtx, &resp); err != nil {
		t.Fatalf("expected no error in List, got %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", resp.Context.StatusCode)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 items in response, got %d", len(resp.Data))
	}
}

func TestRestService_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ types.RequestContext, response *repcts.RestResponse) error {
			response.Context = types.ResponseContext{StatusCode: 200}
			response.Data = &testResponse{ID: "1", Name: "fetched"}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	reqCtx := types.RequestContext{
		Query:       url.Values{},
		Headers:     http.Header{},
		Identifiers: types.Identifiers{"id": "1"},
	}

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Get(context.Background(), reqCtx, &resp); err != nil {
		t.Fatalf("expected no error in Get, got %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", resp.Context.StatusCode)
	}
}

func TestRestService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
			response.Context = types.ResponseContext{StatusCode: 200}
			response.Data = &testResponse{ID: "1", Name: "updated"}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceUpdate[testUpdatePayload]{
		Context: types.RequestContext{
			Query:       url.Values{},
			Headers:     http.Header{},
			Identifiers: types.Identifiers{"id": "1"},
		},
		Body: testUpdatePayload{Name: "updated"},
	}

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Update(context.Background(), req, &resp); err != nil {
		t.Fatalf("expected no error in Update, got %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", resp.Context.StatusCode)
	}
}

func TestRestService_Update_NilResponse_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceUpdate[testUpdatePayload]{
		Context: types.RequestContext{},
		Body:    testUpdatePayload{},
	}

	err := svc.Update(context.Background(), req, nil)
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestRestService_Save_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ repcts.RestRequest, response *repcts.RestResponse) error {
			response.Context = types.ResponseContext{StatusCode: 200}
			response.Data = &testResponse{ID: "1", Name: "saved"}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceSave[testCreatePayload]{
		Context: types.RequestContext{
			Query:       url.Values{},
			Headers:     http.Header{},
			Identifiers: types.Identifiers{"id": "1"},
		},
		Body: testCreatePayload{Name: "saved"},
	}

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Save(context.Background(), req, &resp); err != nil {
		t.Fatalf("expected no error in Save, got %v", err)
	}
	if resp.Context.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", resp.Context.StatusCode)
	}
}

func TestRestService_Save_NilResponse_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceSave[testCreatePayload]{
		Context: types.RequestContext{},
		Body:    testCreatePayload{},
	}

	err := svc.Save(context.Background(), req, nil)
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestRestService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ types.RequestContext, response *repcts.RestResponse) error {
			response.Context = types.ResponseContext{StatusCode: 204}
			return nil
		},
	)

	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	reqCtx := types.RequestContext{
		Query:       url.Values{},
		Headers:     http.Header{},
		Identifiers: types.Identifiers{"id": "1"},
	}

	var resp svccts.ServiceResponse[testResponse]
	if err := svc.Delete(context.Background(), reqCtx, &resp); err != nil {
		t.Fatalf("expected no error in Delete, got %v", err)
	}
	if resp.Context.StatusCode != 204 {
		t.Fatalf("expected status code 204, got %d", resp.Context.StatusCode)
	}
}

func TestRestService_Delete_NilResponse_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repmocks.NewMockRestRepository(ctrl)
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	err := svc.Delete(context.Background(), types.RequestContext{}, nil)
	if err != errRestServiceResponseNil {
		t.Fatalf("expected errRestServiceResponseNil, got %v", err)
	}
}

func TestRestService_Create_UpstreamPostError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&postErrUpstreamMapper{}))

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]

	err := svc.Create(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected upstream post mapping error")
	}
}

func TestRestService_Update_UpstreamPatchError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&patchErrUpstreamMapper{}))

	req := svccts.ServiceUpdate[testUpdatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testUpdatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]

	err := svc.Update(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected upstream patch mapping error")
	}
}

func TestRestService_Save_UpstreamPutError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&putErrUpstreamMapper{}))

	req := svccts.ServiceSave[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]

	err := svc.Save(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected upstream put mapping error")
	}
}

func TestRestService_Create_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{createErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceCreate[testCreatePayload]{Context: types.RequestContext{}, Body: testCreatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Create(context.Background(), req, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_List_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{listErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_Get_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{getErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{}, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_Update_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{updateErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceUpdate[testUpdatePayload]{Context: types.RequestContext{}, Body: testUpdatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Update(context.Background(), req, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_Save_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{saveErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	req := svccts.ServiceSave[testCreatePayload]{Context: types.RequestContext{}, Body: testCreatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Save(context.Background(), req, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_Delete_RepositoryError_ReturnsError(t *testing.T) {
	want := context.Canceled
	repo := &errorRepository{deleteErr: want}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo)

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Delete(context.Background(), types.RequestContext{}, &resp)
	if err != want {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRestService_Create_ContextMappingError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	req := svccts.ServiceCreate[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Create(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected create context mapping error")
	}
}

func TestRestService_Create_DownstreamContextError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	req := svccts.ServiceCreate[testCreatePayload]{Context: types.RequestContext{}, Body: testCreatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Create(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected create downstream context mapping error")
	}
}

func TestRestService_List_ContextMappingError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, sliceData: []repcts.RestResponseSpec{&testResponse{ID: "1"}}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{Query: url.Values{}, Headers: http.Header{}}, &resp)
	if err == nil {
		t.Fatal("expected list context mapping error")
	}
}

func TestRestService_List_DownstreamContextError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, sliceData: []repcts.RestResponseSpec{&testResponse{ID: "1"}}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	var resp svccts.ServiceResponses[testResponse]
	err := svc.List(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected list downstream context mapping error")
	}
}

func TestRestService_Get_ContextMappingError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Get(context.Background(), types.RequestContext{Query: url.Values{}, Headers: http.Header{}}, &resp)
	if err == nil {
		t.Fatal("expected get context mapping error")
	}
}

func TestRestService_Update_ContextMappingError_ReturnsError(t *testing.T) {

	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	req := svccts.ServiceUpdate[testUpdatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testUpdatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Update(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected update context mapping error")
	}
}

func TestRestService_Update_DownstreamContextError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	req := svccts.ServiceUpdate[testUpdatePayload]{Context: types.RequestContext{}, Body: testUpdatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Update(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected update downstream context mapping error")
	}
}

func TestRestService_Save_ContextMappingError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	req := svccts.ServiceSave[testCreatePayload]{
		Context: types.RequestContext{Query: url.Values{}, Headers: http.Header{}},
		Body:    testCreatePayload{},
	}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Save(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected save context mapping error")
	}
}

func TestRestService_Save_DownstreamContextError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200, data: &testResponse{ID: "1"}}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	req := svccts.ServiceSave[testCreatePayload]{Context: types.RequestContext{}, Body: testCreatePayload{}}
	var resp svccts.ServiceResponse[testResponse]
	err := svc.Save(context.Background(), req, &resp)
	if err == nil {
		t.Fatal("expected save downstream context mapping error")
	}
}

func TestRestService_Delete_ContextMappingError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithUpstreamMapper[testCreatePayload, testUpdatePayload](&queryErrUpstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Delete(context.Background(), types.RequestContext{Query: url.Values{}, Headers: http.Header{}}, &resp)
	if err == nil {
		t.Fatal("expected delete context mapping error")
	}
}

func TestRestService_Delete_DownstreamContextError_ReturnsError(t *testing.T) {
	repo := &statusMockRepository{statusCode: 200}
	svc, _ := NewRestService[testCreatePayload, testResponse, testUpdatePayload](repo, WithDownstreamMapper[testResponse](&statusCodeErrDownstreamMapper{}))

	var resp svccts.ServiceResponse[testResponse]
	err := svc.Delete(context.Background(), types.RequestContext{}, &resp)
	if err == nil {
		t.Fatal("expected delete downstream context mapping error")
	}
}
