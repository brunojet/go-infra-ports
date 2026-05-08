package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	restimpl "github.com/brunojet/go-infra-ports/internal/repositories/rest"
)

type testHttpClientStub struct{}

func (testHttpClientStub) Do(_ context.Context, _ *http.Request) (*http.Response, error) {
	return nil, nil
}

// test payload types used with V2 data specs.
type testRequestPayload struct {
	Body json.RawMessage `json:"body"`
}

type testResponsePayload struct {
	Body json.RawMessage `json:"body"`
}

func TestNewRestRegistry_ReturnsNonNil(t *testing.T) {
	if got := NewRestRegistry(); got == nil {
		t.Fatalf("expected non-nil registry")
	}
}

func TestRegistryOptions_ReturnOptions(t *testing.T) {
	// register using V2 helpers and concrete payload types
	opts := []RegistryOption{
		WithRequestOf[testRequestPayload](MethodCreate),
		WithRequestEnvelope("data", MethodCreate),
		WithResponseOf[testResponsePayload](http.StatusOK),
		WithResponseEnvelope("data", "meta", http.StatusOK),
		WithInformationOf[testResponsePayload](http.StatusContinue),
		WithRedirectionOf[testResponsePayload](http.StatusMovedPermanently),
		WithProblemOf[testResponsePayload](http.StatusBadRequest),
	}

	for i, opt := range opts {
		if opt == nil {
			t.Fatalf("expected non-nil option at index %d", i)
		}
	}
}

func TestAliases_AreUsable(t *testing.T) {
	request := RestRequest{Data: restimpl.NewDataSpecOf[testRequestPayload]()}
	response := RestResponse{Data: restimpl.NewDataSpecOf[testResponsePayload]()}
	responses := RestResponses{Data: []restimpl.RestDataSpec{restimpl.NewDataSpecOf[testResponsePayload]()}}

	if request.Data == nil {
		t.Fatalf("expected request data alias to be usable")
	}
	if response.Data == nil {
		t.Fatalf("expected response data alias to be usable")
	}
	if len(responses.Data) != 1 {
		t.Fatalf("expected responses data alias to be usable")
	}
}

func TestNewRestRepository_ReturnsNonNil(t *testing.T) {
	reg := NewRestRegistry()
	repo := NewRestRepository(
		WithHttpClient(&testHttpClientStub{}),
		WithRegistryOpt(reg),
	)
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}

func TestNewRestRepository_MissingClient_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing HttpClient")
		}
	}()
	NewRestRepository(WithRegistryOpt(NewRestRegistry()))
}

func TestNewRestRepository_WithValidClient_ReturnsNonNil(t *testing.T) {
	repo := NewRestRepository(
		WithHttpClient(&testHttpClientStub{}),
		WithRegistryOpt(NewRestRegistry()),
	)
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}
