package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/brunojet/go-infra-ports/pkg/http_clients/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

// newMockClient creates a MockHttpClient that captures the last request and returns the given status/body.
func newMockClient(t *testing.T, statusCode int, body string) (*mocks.MockHttpClient, *http.Request) {
	t.Helper()
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var captured *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			captured = req
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		},
	)
	return mock, captured
}

func newTestRepo(t *testing.T, client HttpClient) RestRepository {
	t.Helper()
	reg := NewRestRegistry(
		WithResponseOf[DefaultRestResponse](http.StatusOK, http.StatusCreated, http.StatusNoContent),
	)
	return NewRestRepository(
		WithHttpClient(client),
		WithRegistry(reg),
		WithPath(MethodCreate, "/items"),
		WithPath(MethodList, "/items"),
		WithPath(MethodGet, "/items/{id}"),
		WithPath(MethodUpdate, "/items/{id}"),
		WithPath(MethodSave, "/items/{id}"),
		WithPath(MethodDelete, "/items/{id}"),
	)
}

func TestNewRestRepository_MissingClient_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing HttpClient")
		}
	}()
	NewRestRepository(WithRegistry(registryStub{}))
}

func TestNewRestRepository_MissingRegistry_DoesNotPanic(t *testing.T) {
	// registry defaults to NewRestRegistry() — omitting WithRegistry is valid
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	repo := NewRestRepository(WithHttpClient(newNoOpMockClient(t)))
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}

func TestCreate_SendsPOSTAndMapsResponse(t *testing.T) {
	mock, _ := newMockClient(t, http.StatusCreated, `{}`)
	repo := newTestRepo(t, mock)

	resp := &RestResponse{}
	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{"name":"x"}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}
	err := repo.Create(context.Background(), RestRequest{
		Context: types.RequestContext{},
		Data:    reqBody,
	}, resp)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if resp.Context.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Context.StatusCode)
	}
}

func TestList_SendsGETAndMapsResponses(t *testing.T) {
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var capturedReq *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[{}]`))}, nil
		},
	)
	repo := newTestRepo(t, mock)

	out := &RestResponses{}
	if err := repo.List(context.Background(), types.RequestContext{}, out); err != nil {
		t.Fatalf("List: %v", err)
	}
	if capturedReq.Method != http.MethodGet {
		t.Fatalf("expected GET, got %s", capturedReq.Method)
	}
	if out.Context.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", out.Context.StatusCode)
	}
}

func TestGet_SendsGETWithIdentifierInURL(t *testing.T) {
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var capturedReq *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		},
	)
	repo := newTestRepo(t, mock)

	out := &RestResponse{}
	if err := repo.Get(context.Background(), types.RequestContext{
		Identifiers: types.Identifiers{"id": "99"},
	}, out); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if capturedReq.Method != http.MethodGet {
		t.Fatalf("expected GET, got %s", capturedReq.Method)
	}
	if !strings.Contains(capturedReq.URL.Path, "99") {
		t.Fatalf("expected identifier in URL path, got %s", capturedReq.URL.Path)
	}
}

func TestUpdate_SendsPUT(t *testing.T) {
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var capturedReq *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		},
	)
	repo := newTestRepo(t, mock)

	out := &RestResponse{}
	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}
	if err := repo.Update(context.Background(), RestRequest{
		Context: types.RequestContext{Identifiers: types.Identifiers{"id": "1"}},
		Data:    reqBody,
	}, out); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if capturedReq.Method != http.MethodPut {
		t.Fatalf("expected PUT, got %s", capturedReq.Method)
	}
}

func TestSave_SendsPATCH(t *testing.T) {
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var capturedReq *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		},
	)
	repo := newTestRepo(t, mock)

	out := &RestResponse{}
	reqBody := NewDataSpecOf[DefaultRestRequest]()
	if err := reqBody.SetBody(json.RawMessage(`{}`)); err != nil {
		t.Fatalf("SetBody failed: %v", err)
	}
	if err := repo.Save(context.Background(), RestRequest{
		Context: types.RequestContext{Identifiers: types.Identifiers{"id": "2"}},
		Data:    reqBody,
	}, out); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if capturedReq.Method != http.MethodPatch {
		t.Fatalf("expected PATCH, got %s", capturedReq.Method)
	}
}

func TestDelete_SendsDELETE(t *testing.T) {
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	var capturedReq *http.Request
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, req *http.Request) (*http.Response, error) {
			capturedReq = req
			return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
		},
	)
	repo := newTestRepo(t, mock)

	out := &RestResponse{}
	if err := repo.Delete(context.Background(), types.RequestContext{
		Identifiers: types.Identifiers{"id": "3"},
	}, out); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if capturedReq.Method != http.MethodDelete {
		t.Fatalf("expected DELETE, got %s", capturedReq.Method)
	}
}
