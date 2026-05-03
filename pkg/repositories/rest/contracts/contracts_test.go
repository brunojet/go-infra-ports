package contracts

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type testReqSpec struct{}

func (testReqSpec) New() RestRequestSpec    { return testReqSpec{} }
func (testReqSpec) SetBody(RestRequestSpec) {}

type testRespSpec struct{}

func (testRespSpec) New() RestResponseSpec { return testRespSpec{} }
func (testRespSpec) NewSlice(n int) []RestResponseSpec {
	out := make([]RestResponseSpec, n)
	for i := range out {
		out[i] = testRespSpec{}
	}
	return out
}

type testEnvSpec struct{}

func (testEnvSpec) New() RestEnvelopeSpec            { return testEnvSpec{} }
func (testEnvSpec) EnvelopeData() json.RawMessage    { return json.RawMessage(`{"ok":true}`) }
func (testEnvSpec) EnvelopeMeta() types.ResponseMeta { return types.ResponseMeta{"page": 1} }

type testRegistry struct{}

func (testRegistry) Merge(other RestRegistry) RestRegistry                           { return other }
func (testRegistry) ResolveRequest(RestRequestSpec, *[]byte) error                   { return nil }
func (testRegistry) ResolveEnvelopeRequest(string, *[]byte) error                    { return nil }
func (testRegistry) ResolveResponse(int, []byte, *RestResponseSpec) error            { return nil }
func (testRegistry) ResolveResponses(int, []byte, *[]RestResponseSpec) error         { return nil }
func (testRegistry) ResolveEnvelopeResponse(int, *[]byte, *types.ResponseMeta) error { return nil }
func (testRegistry) NewRequestSpec(string) (RestRequestSpec, error)                  { return testReqSpec{}, nil }
func (testRegistry) ReleaseRequestSpec(RestRequestSpec)                              {}

type testRepository struct{}

func (testRepository) Create(context.Context, RestRequest, *RestResponse) error          { return nil }
func (testRepository) List(context.Context, types.RequestContext, *RestResponses) error  { return nil }
func (testRepository) Get(context.Context, types.RequestContext, *RestResponse) error    { return nil }
func (testRepository) Update(context.Context, RestRequest, *RestResponse) error          { return nil }
func (testRepository) Save(context.Context, RestRequest, *RestResponse) error            { return nil }
func (testRepository) Delete(context.Context, types.RequestContext, *RestResponse) error { return nil }

var _ RestRequestSpec = testReqSpec{}
var _ RestResponseSpec = testRespSpec{}
var _ RestEnvelopeSpec = testEnvSpec{}
var _ RestRegistry = testRegistry{}
var _ RestRepository = testRepository{}

func TestContracts_RequestResponseEnvelopeAndInterfaces(t *testing.T) {
	var req RestRequestSpec = testReqSpec{}
	var resp RestResponseSpec = testRespSpec{}
	var env RestEnvelopeSpec = testEnvSpec{}
	_, _, _ = req, resp, env

	if len(resp.NewSlice(2)) != 2 {
		t.Fatal("expected response slice with len=2")
	}
	if env.EnvelopeMeta()["page"] != 1 {
		t.Fatal("expected envelope meta page=1")
	}

	var registry RestRegistry = testRegistry{}
	var repo RestRepository = testRepository{}
	_, _ = registry, repo
}
