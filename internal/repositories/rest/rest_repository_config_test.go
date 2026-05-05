package rest

import (
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/brunojet/go-infra-ports/pkg/http_clients/mocks"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

// newNoOpMockClient returns a MockHttpClient that accepts any Do call (AnyTimes).
func newNoOpMockClient(t *testing.T) *mocks.MockHttpClient {
	t.Helper()
	mock := mocks.NewMockHttpClient(gomock.NewController(t))
	mock.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	return mock
}

type registryStub struct{}

func (registryStub) Merge(other RestRegistry) RestRegistry                           { return other }
func (registryStub) ResolveRequest(RestRequestSpec, *[]byte) error                   { return nil }
func (registryStub) ResolveEnvelopeRequest(_ RestMethod, dataBody *[]byte) error     { return nil }
func (registryStub) ResolveResponse(int, []byte, *RestResponseSpec) error            { return nil }
func (registryStub) ResolveResponses(int, []byte, *[]RestResponseSpec) error         { return nil }
func (registryStub) ResolveEnvelopeResponse(int, *[]byte, *types.ResponseMeta) error { return nil }
func (registryStub) NewRequestSpec(_ RestMethod) (RestRequestSpec, error)            { return nil, nil }
func (registryStub) ReleaseRequestSpec(RestRequestSpec)                              {}

func validOpts(t *testing.T) []RepositoryOption {
	t.Helper()
	return []RepositoryOption{WithHttpClient(newNoOpMockClient(t)), WithRegistry(registryStub{})}
}

func mustOpts(t *testing.T, opts ...RepositoryOption) *repositoryOptions {
	t.Helper()
	return newRepositoryOptions(append(validOpts(t), opts...)...)
}

// --- client ---

func TestWithHttpClient_Nil_Panics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil HttpClient")
		}
	}()
	newRepositoryOptions(WithHttpClient(nil), WithRegistry(registryStub{}))
}

func TestNewRepositoryOptions_MissingClient_Panics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for missing HttpClient")
		}
	}()
	newRepositoryOptions(WithRegistry(registryStub{}))
}

// --- registry ---

func TestWithRegistry_Nil_Panics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil RestRegistry")
		}
	}()
	newRepositoryOptions(WithHttpClient(newNoOpMockClient(t)), WithRegistry(nil))
}

// --- basePath ---

func TestWithBasePath_Empty_Panics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for empty basePath")
		}
	}()
	newRepositoryOptions(append(validOpts(t), WithBasePath(""))...)
}

func TestWithBasePath_Valid_Sets(t *testing.T) {
	o := mustOpts(t, WithBasePath("/api/v1/users"))
	if o.basePath != "/api/v1/users" {
		t.Fatalf("expected basePath=/api/v1/users, got %q", o.basePath)
	}
}

// --- WithPath ---

func TestWithPath_CollectionMethods_NoParams(t *testing.T) {
	o := mustOpts(t, WithPath(MethodCreate|MethodList, "/users"))
	for _, m := range []RestMethod{MethodCreate, MethodList} {
		e := o.paths[m]
		if e == nil || e.templateFmt != "/users" {
			t.Fatalf("expected templateFmt=/users for method %v, got %v", m, e)
		}
		if len(e.paramNames) != 0 {
			t.Fatalf("expected no paramNames for method %v, got %v", m, e.paramNames)
		}
	}
}

func TestWithPath_InstanceMethods_WithParams(t *testing.T) {
	o := mustOpts(t, WithPath(MethodGet|MethodUpdate|MethodSave|MethodDelete, "/users/{id}"))
	for _, m := range []RestMethod{MethodGet, MethodUpdate, MethodSave, MethodDelete} {
		e := o.paths[m]
		if e == nil || e.templateFmt != "/users/%s" {
			t.Fatalf("expected templateFmt=/users/%%s for method %v, got %v", m, e)
		}
		if len(e.paramNames) != 1 || e.paramNames[0] != "id" {
			t.Fatalf("expected paramNames=[id] for method %v, got %v", m, e.paramNames)
		}
	}
}

func TestWithPath_AllMethods_WithParams(t *testing.T) {
	o := mustOpts(t, WithPath(AllCollectionMethods|AllInstanceMethods, "/resources/{id}"))
	for _, m := range []RestMethod{MethodCreate, MethodList, MethodGet, MethodUpdate, MethodSave, MethodDelete} {
		e := o.paths[m]
		if e == nil || e.templateFmt != "/resources/%s" {
			t.Fatalf("expected templateFmt=/resources/%%s for method %v, got %v", m, e)
		}
	}
}

func TestWithPath_InvalidTemplate_Panics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invalid path template")
		}
	}()
	newRepositoryOptions(append(validOpts(t), WithPath(MethodCreate, "/users/{}"))...)
}

// --- WithHeader ---

func TestWithHeader_Sets(t *testing.T) {
	o := mustOpts(t, WithHeader("X-Tenant", "acme"))
	if o.headers.Get("X-Tenant") != "acme" {
		t.Fatalf("expected X-Tenant=acme, got %q", o.headers.Get("X-Tenant"))
	}
}

// --- defaults ---

func TestNewRepositoryOptions_Defaults(t *testing.T) {
	o := mustOpts(t)
	for _, m := range []RestMethod{MethodCreate, MethodList} {
		if o.paths[m] != DefaultCollectionPathEntry {
			t.Fatalf("expected DefaultCollectionPathEntry for method %v", m)
		}
	}
	for _, m := range []RestMethod{MethodGet, MethodUpdate, MethodSave, MethodDelete} {
		if o.paths[m] != DefaultInstancePathEntry {
			t.Fatalf("expected DefaultInstancePathEntry for method %v", m)
		}
	}
}

// --- pathEntry.fullPath ---

func TestFullPath_NoParams_ReturnsBasePathPlusTemplate(t *testing.T) {
	e := &pathEntry{templateFmt: "/users"}
	got, err := e.expandPath(types.Identifiers{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/users" {
		t.Fatalf("expected /users, got %q", got)
	}
}

func TestFullPath_WithParams_ReturnsInterpolatedPath(t *testing.T) {
	e := &pathEntry{templateFmt: "/users/%s", paramNames: []string{"id"}}
	got, err := e.expandPath(types.Identifiers{"id": "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/users/42" {
		t.Fatalf("expected /users/42, got %q", got)
	}
}

func TestFullPath_MultipleParams_ReturnsInterpolatedPath(t *testing.T) {
	e := &pathEntry{templateFmt: "/orgs/%s/users/%s", paramNames: []string{"org", "id"}}
	got, err := e.expandPath(types.Identifiers{"org": "acme", "id": "7"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/orgs/acme/users/7" {
		t.Fatalf("expected /orgs/acme/users/7, got %q", got)
	}
}

func TestFullPath_WrongIDCount_ReturnsError(t *testing.T) {
	e := &pathEntry{templateFmt: "/users/%s", paramNames: []string{"id"}}
	_, err := e.expandPath(types.Identifiers{})
	if err == nil {
		t.Fatal("expected error for wrong identifier count")
	}
}

func TestFullPath_MissingIDKey_ReturnsError(t *testing.T) {
	e := &pathEntry{templateFmt: "/users/%s", paramNames: []string{"id"}}
	_, err := e.expandPath(types.Identifiers{"other": "42"})
	if err == nil {
		t.Fatal("expected error for missing identifier key")
	}
}
