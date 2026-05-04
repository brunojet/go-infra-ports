package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	restcontracts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type testHttpClientStub struct{}

func (testHttpClientStub) Do(_ context.Context, _ *http.Request) (*http.Response, error) {
	return nil, nil
}

type testRequestSpec struct {
	Body json.RawMessage
}

func (t *testRequestSpec) New() restcontracts.RestRequestSpec {
	return &testRequestSpec{}
}

func (t *testRequestSpec) SetBody(body restcontracts.RestRequestSpec) {}

type testResponseSpec struct {
	Body json.RawMessage
}

func (t *testResponseSpec) New() restcontracts.RestResponseSpec {
	return &testResponseSpec{}
}

func (t *testResponseSpec) NewSlice(n int) []restcontracts.RestResponseSpec {
	out := make([]restcontracts.RestResponseSpec, n)
	for i := 0; i < n; i++ {
		out[i] = &testResponseSpec{}
	}
	return out
}

func (t *testResponseSpec) UnmarshalJSON(data []byte) error {
	t.Body = append([]byte(nil), data...)
	return nil
}

type testEnvelopeSpec struct {
	Data json.RawMessage
	Meta types.ResponseMeta
}

func (t *testEnvelopeSpec) New() restcontracts.RestEnvelopeSpec {
	return &testEnvelopeSpec{}
}

func (t *testEnvelopeSpec) EnvelopeData() json.RawMessage {
	return t.Data
}

func (t *testEnvelopeSpec) EnvelopeMeta() types.ResponseMeta {
	return t.Meta
}

func TestNewRestRegistry_ReturnsNonNil(t *testing.T) {
	if got := NewRestRegistry(); got == nil {
		t.Fatalf("expected non-nil registry")
	}
}

func TestRegistryOptions_ReturnOptions(t *testing.T) {
	req := &testRequestSpec{}
	resp := &testResponseSpec{}
	env := &testEnvelopeSpec{}

	opts := []RegistryOption{
		WithRequest(req, restcontracts.MethodCreate),
		WithRequestEnvelope(req, restcontracts.MethodCreate),
		WithResponse(resp, http.StatusOK),
		WithResponseEnvelope(env, http.StatusOK),
		WithInformation(resp, http.StatusContinue),
		WithRedirection(resp, http.StatusMovedPermanently),
		WithProblem(resp, http.StatusBadRequest),
	}

	for i, opt := range opts {
		if opt == nil {
			t.Fatalf("expected non-nil option at index %d", i)
		}
	}
}

func TestAliases_AreUsable(t *testing.T) {
	request := RestRequest{Body: &testRequestSpec{}}
	response := RestResponse{Data: &testResponseSpec{}}
	responses := RestResponses{Data: []RestResponseSpec{&testResponseSpec{}}}

	if request.Body == nil {
		t.Fatalf("expected request body alias to be usable")
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
