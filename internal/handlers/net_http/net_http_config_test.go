package net_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	mocksvc "github.com/brunojet/go-infra-ports/pkg/services/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func TestRegisterMethodSelection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mux := http.NewServeMux()
	// only Create on collection, only Get on instance
	NewNetHttpHandler(mux, mock, WithCollection(MethodCreate, "events"), WithInstance(MethodGet, "events/{id}"))
	// GET on collection should be 404
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/events", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected 404 for unregistered GET on collection")
	}
	// POST should work (no assert on body) - we expect the mock Create to be called once
	mock.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req svccts.ServiceCreate[userCreate], resp *svccts.ServiceResponse[userResp]) error {
			resp.Context.StatusCode = 0
			resp.Data = userResp{ID: "1", Name: req.Body.Name}
			return nil
		}).Times(1)

	rr = httptest.NewRecorder()
	cb, _ := json.Marshal(userCreate{Name: "X"})
	req = httptest.NewRequest("POST", "http://example.com/events", bytes.NewReader(cb))
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 for POST, got %d", rr.Code)
	}
}

func TestWithCollectionSanitizesPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req interface{}, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "A"}}
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	// leading + trailing slash should be stripped
	NewNetHttpHandler(mux, mock, WithCollection(MethodList, "/items/"))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/items", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 after path sanitization, got %d", rr.Code)
	}
}

func TestWithNilBuildersPanic(t *testing.T) {
	cases := []struct {
		name string
		fn   func()
	}{
		{"WithMiddlewareRegistry", func() { WithMiddlewareRegistry(nil) }},
		{"WithInformation", func() { WithInformation(nil) }},
		{"WithEnvelope", func() { WithEnvelope(nil) }},
		{"WithRedirection", func() { WithRedirection(nil) }},
		{"WithProblemDetails", func() { WithProblemDetails(nil) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("%s(nil) should panic", tc.name)
				}
			}()
			tc.fn()
		})
	}
}

func TestWithMiddlewareRegistry_WrapsRegisteredRoutes_OrderAndHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "A"}}
			return nil
		}).Times(1)

	var order []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			w.Header().Add("X-Middleware-Order", "mw1")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			w.Header().Add("X-Middleware-Order", "mw2")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	registry := NewMiddlewareRegistry()
	registry.Register(mw1, mw2)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithCollection(MethodList, "users"),
		WithMiddlewareRegistry(registry),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Header().Values("X-Middleware-Order"); fmt.Sprint(got) != "[mw1 mw2]" {
		t.Fatalf("expected middleware headers in registration order, got %v", got)
	}
	if got := fmt.Sprint(order); got != "[mw1-before mw2-before mw2-after mw1-after]" {
		t.Fatalf("unexpected middleware call order: %s", got)
	}
}

func TestWithMiddlewareRegistry_CopiesRegistrySnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "A"}}
			return nil
		}).Times(1)

	var order []string
	mws := []Middleware{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw1-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw1-after")
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw2-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw2-after")
			})
		},
	}
	registry := NewMiddlewareRegistry()
	registry.Register(mws...)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithCollection(MethodList, "users"),
		WithMiddlewareRegistry(registry),
	)

	mws[0] = nil
	registry.Reset()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := fmt.Sprint(order); got != "[mw1-before mw2-before mw2-after mw1-after]" {
		t.Fatalf("unexpected middleware call order: %s", got)
	}
}

func TestWithMiddlewareRegistry_UsesRegisteredMiddlewares(t *testing.T) {
	registry := NewMiddlewareRegistry()
	registry.Register(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-Registry", "mw1")
				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-Registry", "mw2")
				next.ServeHTTP(w, r)
			})
		},
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[userCreate, userResp, userUpdate](ctrl)
	mock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, req types.RequestContext, resp *svccts.ServiceResponses[userResp]) error {
			resp.Data = []userResp{{ID: "1", Name: "A"}}
			return nil
		}).Times(1)

	mux := http.NewServeMux()
	NewNetHttpHandler(mux, mock,
		WithCollection(MethodList, "users"),
		WithMiddlewareRegistry(registry),
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Header().Values("X-Registry"); fmt.Sprint(got) != "[mw1 mw2]" {
		t.Fatalf("expected registry middleware headers in registration order, got %v", got)
	}
}

func TestNewHandlerOptions_DefaultBuilders(t *testing.T) {
	opts := newHandlerOptions(nil)

	info := opts.informationBuilder(http.StatusContinue, svccts.ServiceMeta{})
	gotInfo, ok := info.(DefaultInformation)
	if !ok {
		t.Fatalf("expected DefaultInformation type, got %T", info)
	}
	if gotInfo.Code != http.StatusContinue {
		t.Fatalf("expected info code %d, got %d", http.StatusContinue, gotInfo.Code)
	}
	if gotInfo.Message != "Continue" {
		t.Fatalf("expected default info message Continue, got %q", gotInfo.Message)
	}

	redir, location := opts.redirectionBuilder(http.StatusFound, svccts.ServiceMeta{Location: "https://example.com/new"})
	gotRedir, ok := redir.(DefaultRedirection)
	if !ok {
		t.Fatalf("expected DefaultRedirection type, got %T", redir)
	}
	if gotRedir.Location != "https://example.com/new" || location != "https://example.com/new" {
		t.Fatalf("expected redirection location propagated, got body=%q header=%q", gotRedir.Location, location)
	}

	problemAny := opts.problemDetailsBuilder(http.StatusBadRequest, "/r/1", svccts.ServiceMeta{})
	problem, ok := problemAny.(DefaultProblemDetails)
	if !ok {
		t.Fatalf("expected DefaultProblemDetails type, got %T", problemAny)
	}
	if problem.Status != http.StatusBadRequest {
		t.Fatalf("expected problem status %d, got %d", http.StatusBadRequest, problem.Status)
	}
	if problem.Instance != "/r/1" {
		t.Fatalf("expected problem instance /r/1, got %q", problem.Instance)
	}
}
