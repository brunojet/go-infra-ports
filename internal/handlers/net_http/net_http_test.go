package net_http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/brunojet/go-infra-ports/internal/dtos"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	mocksvc "github.com/brunojet/go-infra-ports/pkg/services/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func TestCreate_List_Get_Update_Delete_Save_and_Errors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)

	// Create
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceCreate[userCreate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = 0
			resp.Data = userResp{ID: "1", Name: req.Body.Name}
			return nil
		}).Times(1)

	// List
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "Alice"}}
			return nil
		}).Times(1)

	var gotGetIdentifiers map[string]string
	// Get
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			gotGetIdentifiers = req.Identifiers
			resp.Data = userResp{ID: req.Identifiers["userId"], Name: "FromSvc"}
			return nil
		}).Times(1)

	// Update
	mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceUpdate[userUpdate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Data = userResp{ID: "1", Name: req.Body.Name}
			return nil
		}).Times(1)

	// Save
	mock.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceSave[userCreate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Data = userResp{ID: "2", Name: req.Body.Name}
			return nil
		}).Times(1)

	// Delete
	mock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(AllCollectionMethods, "users"), WithInstance(AllInstanceMethods, "users/{userId}"))

	// Create
	createBody := userCreate{Name: "Alice"}
	b, _ := json.Marshal(createBody)
	req := httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	var creResp svccts.ServiceResponse[userResp]
	if err := json.Unmarshal(rr.Body.Bytes(), &creResp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if creResp.Data.Name != "Alice" {
		t.Fatalf("unexpected create name")
	}

	// List
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var listResp svccts.ServiceResponses[userResp]
	if err := json.Unmarshal(rr.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(listResp.Data) != 1 || listResp.Data[0].Name != "Alice" {
		t.Fatalf("unexpected list data")
	}

	// Get with path param extraction
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users/123", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var getResp svccts.ServiceResponse[userResp]
	if err := json.Unmarshal(rr.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("unmarshal get: %v", err)
	}
	if getResp.Data.ID != "123" {
		t.Fatalf("expected id 123, got %s", getResp.Data.ID)
	}
	if gotGetIdentifiers["userId"] != "123" {
		t.Fatalf("identifier not propagated")
	}

	// Update with valid JSON
	rr = httptest.NewRecorder()
	updateBody := userUpdate{Name: "Bob"}
	ub, _ := json.Marshal(updateBody)
	req = httptest.NewRequest("PUT", "http://example.com/users/123", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 update, got %d", rr.Code)
	}

	// Save (PATCH)
	rr = httptest.NewRecorder()
	saveBody := userCreate{Name: "Saved"}
	sb, _ := json.Marshal(saveBody)
	req = httptest.NewRequest("PATCH", "http://example.com/users/123", bytes.NewReader(sb))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 save, got %d", rr.Code)
	}

	// Delete
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "http://example.com/users/123", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 delete, got %d", rr.Code)
	}

	// Invalid body -> 400
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "http://example.com/users/123", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 invalid body, got %d", rr.Code)
	}

	// Service error -> 500
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	svcErr := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl2)
	svcErr.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	muxErr := http.NewServeMux()
	NewNetHttpHandler(muxErr, svcErr, WithCollection(MethodCreate, "users"))
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	muxErr.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 on service error, got %d", rr.Code)
	}
}

func TestHandlers_ServiceErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)

	// all methods return error
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	mock.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)
	mock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(AllCollectionMethods, "users"), WithInstance(AllInstanceMethods, "users/{userId}"))

	// Create
	b, _ := json.Marshal(userCreate{Name: "X"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 create, got %d", rr.Code)
	}

	// List
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 list, got %d", rr.Code)
	}

	// Get
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users/1", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 get, got %d", rr.Code)
	}

	// Update
	rr = httptest.NewRecorder()
	ub, _ := json.Marshal(userUpdate{Name: "u"})
	req = httptest.NewRequest("PUT", "http://example.com/users/1", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 update, got %d", rr.Code)
	}

	// Save
	rr = httptest.NewRecorder()
	sb, _ := json.Marshal(userCreate{Name: "s"})
	req = httptest.NewRequest("PATCH", "http://example.com/users/1", bytes.NewReader(sb))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 save, got %d", rr.Code)
	}

	// Delete
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "http://example.com/users/1", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 delete, got %d", rr.Code)
	}
}

func TestHandlers_CustomStatusCodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)

	// Create -> 202
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceCreate[userCreate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusAccepted
			resp.Data = userResp{ID: "10", Name: req.Body.Name}
			return nil
		}).Times(1)

	// List -> 201
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Context.StatusCode = http.StatusCreated
			resp.Data = []userResp{{ID: "10", Name: "L"}}
			return nil
		}).Times(1)

	// Get -> 202
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusAccepted
			resp.Data = userResp{ID: req.Identifiers["userId"], Name: "G"}
			return nil
		}).Times(1)

	// Update -> 203
	mock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceUpdate[userUpdate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = 203
			resp.Data = userResp{ID: "u", Name: req.Body.Name}
			return nil
		}).Times(1)

	// Save -> 206
	mock.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceSave[userCreate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = 206
			resp.Data = userResp{ID: "s", Name: req.Body.Name}
			return nil
		}).Times(1)

	// Delete -> 200 (override default 204)
	mock.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusOK
			resp.Data = userResp{ID: req.Identifiers["userId"], Name: "D"}
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(AllCollectionMethods, "users"), WithInstance(AllInstanceMethods, "users/{userId}"))

	// Create -> 202
	b, _ := json.Marshal(userCreate{Name: "A"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202 create, got %d", rr.Code)
	}

	// List -> 201
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 list, got %d", rr.Code)
	}

	// Get -> 202
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users/7", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202 get, got %d", rr.Code)
	}

	// Update -> 203
	rr = httptest.NewRecorder()
	ub, _ := json.Marshal(userUpdate{Name: "U"})
	req = httptest.NewRequest("PUT", "http://example.com/users/7", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != 203 {
		t.Fatalf("expected 203 update, got %d", rr.Code)
	}

	// Save -> 206
	rr = httptest.NewRecorder()
	sb, _ := json.Marshal(userCreate{Name: "S"})
	req = httptest.NewRequest("PATCH", "http://example.com/users/7", bytes.NewReader(sb))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != 206 {
		t.Fatalf("expected 206 save, got %d", rr.Code)
	}

	// Delete -> 200
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "http://example.com/users/7", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 delete, got %d", rr.Code)
	}
}

func TestCreate_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(MethodCreate, "users"))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 create invalid body, got %d", rr.Code)
	}
}

func TestSave_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithInstance(MethodSave, "users/{userId}"))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "http://example.com/users/1", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 save invalid body, got %d", rr.Code)
	}
}

func TestCreate_InvalidParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithCollection(MethodCreate, "users/{userId}", rx))

	b, _ := json.Marshal(userCreate{Name: "x"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/users/abc", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 create param invalid, got %d", rr.Code)
	}
}

func TestSave_InvalidParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithInstance(MethodSave, "users/{userId}", rx))

	b, _ := json.Marshal(userCreate{Name: "x"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "http://example.com/users/abc", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 save param invalid, got %d", rr.Code)
	}
}

func TestList_ParamsValidationOnCollection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithCollection(AllCollectionMethods, "users/{userId}", rx))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/abc", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 list param invalid, got %d", rr.Code)
	}
}

func TestUpdate_InvalidParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithInstance(AllInstanceMethods, "users/{userId}", rx))

	rr := httptest.NewRecorder()
	ub, _ := json.Marshal(userUpdate{Name: "x"})
	req := httptest.NewRequest("PUT", "http://example.com/users/abc", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 update param invalid, got %d", rr.Code)
	}
}

func TestDelete_InvalidParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	rx := regexp.MustCompile(`^\d+$`)
	NewNetHttpHandler(mux, mock, WithInstance(AllInstanceMethods, "users/{userId}", rx))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "http://example.com/users/abc", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 delete param invalid, got %d", rr.Code)
	}
}

func TestHandler_ServiceErrorReturnsProblemDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("boom")).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(MethodCreate, "users"))

	body, _ := json.Marshal(userCreate{Name: "A"})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://example.com/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected %s, got %s", dtos.ContentTypeProblemJSON, ct)
	}
	var p DefaultProblemDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Detail != "boom" {
		t.Fatalf("expected detail boom, got %s", p.Detail)
	}
	if p.Instance != "/users" {
		t.Fatalf("expected instance /users, got %s", p.Instance)
	}
}

func TestHandler_StatusGte300UsesDefaultDtos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)

	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusFound
			resp.Meta.Location = "https://example.com/new-resource"
			resp.Meta.Code = "MOVED_TEMPORARILY"
			return nil
		}).Times(1)

	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Context.StatusCode = http.StatusNotFound
			resp.Meta.Message = "resource not found"
			resp.Meta.Code = "RESOURCE_NOT_FOUND"
			resp.Meta.Metadata = map[string]any{"entity": "user", "id": "42"}
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock, WithCollection(MethodList, "users"), WithInstance(MethodGet, "users/{userId}"))

	// 302 -> Redirection body + Location header.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/42", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "https://example.com/new-resource" {
		t.Fatalf("expected location header, got %s", got)
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeJSON {
		t.Fatalf("expected %s, got %s", dtos.ContentTypeJSON, ct)
	}
	var redir DefaultRedirection
	if err := json.Unmarshal(rr.Body.Bytes(), &redir); err != nil {
		t.Fatalf("unmarshal redirection: %v", err)
	}
	if redir.Location == "" {
		t.Fatalf("expected redirection location in body")
	}

	// 404 -> ProblemDetails from ServiceMeta.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected %s, got %s", dtos.ContentTypeProblemJSON, ct)
	}
	var p DefaultProblemDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &p); err != nil {
		t.Fatalf("unmarshal problem details: %v", err)
	}
	if p.Title != "resource not found" {
		t.Fatalf("expected title from meta message, got %s", p.Title)
	}
	if p.Status != http.StatusNotFound {
		t.Fatalf("expected status 404 in body, got %d", p.Status)
	}
}

func TestWithProblemDetails_CustomBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusUnprocessableEntity
			resp.Meta.Message = "validation failed"
			return nil
		}).Times(1)

	type customProblem struct {
		Error  string `json:"error"`
		Status int    `json:"status"`
	}
	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithInstance(MethodGet, "users/{id}"),
		WithProblemDetails(func(status int, instance string, meta svccts.ServiceMeta) any {
			return customProblem{Error: meta.Message, Status: status}
		}),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/1", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
	var got customProblem
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Error != "validation failed" {
		t.Fatalf("expected custom error field, got %q", got.Error)
	}
	if got.Status != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422 in body, got %d", got.Status)
	}
}

func TestWithEnvelope_CustomBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "Alice"}}
			resp.Meta.Metadata = map[string]any{"total": 1}
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithCollection(MethodList, "users"),
		WithEnvelope(func(status int, meta svccts.ServiceMeta, resp any) any {
			return map[string]any{"data": resp, "meta": meta.Metadata}
		}),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got["data"]; !ok {
		t.Fatalf("expected 'data' key in envelope")
	}
	if _, ok := got["meta"]; !ok {
		t.Fatalf("expected 'meta' key in envelope")
	}
}

func TestWithRedirection_CustomBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusMovedPermanently
			resp.Meta.Location = "https://example.com/new"
			return nil
		}).Times(1)

	type customRedirect struct {
		To string `json:"to"`
	}
	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithInstance(MethodGet, "users/{id}"),
		WithRedirection(func(status int, meta svccts.ServiceMeta) (any, string) {
			return customRedirect{To: meta.Location}, meta.Location
		}),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/1", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rr.Code)
	}
	if got := rr.Header().Get("Location"); got != "https://example.com/new" {
		t.Fatalf("expected Location header, got %q", got)
	}
	var got customRedirect
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.To != "https://example.com/new" {
		t.Fatalf("expected custom 'to' field, got %q", got.To)
	}
}

func TestWithInformation_CustomBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = http.StatusContinue // 100
			resp.Meta.Message = "keep going"
			return nil
		}).Times(1)

	type customInfo struct {
		Hint string `json:"hint"`
	}
	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithInstance(MethodGet, "users/{id}"),
		WithInformation(func(status int, meta svccts.ServiceMeta) any {
			return customInfo{Hint: meta.Message}
		}),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users/1", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusContinue {
		t.Fatalf("expected 100, got %d", rr.Code)
	}
	var got customInfo
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Hint != "keep going" {
		t.Fatalf("expected custom hint, got %q", got.Hint)
	}
}
