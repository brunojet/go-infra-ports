package contracts

import (
	"net/http"
	"net/url"
	"testing"
)

type testUpstreamMapper struct{}

func (testUpstreamMapper) ToUpstreamPost(string, Identifiers, *RestRequestSpec) error  { return nil }
func (testUpstreamMapper) ToUpstreamPut(string, Identifiers, *RestRequestSpec) error   { return nil }
func (testUpstreamMapper) ToUpstreamPatch(string, Identifiers, *RestRequestSpec) error { return nil }
func (testUpstreamMapper) ToUpstreamQuery(url.Values, url.Values) error                { return nil }
func (testUpstreamMapper) ToUpstreamHeaders(http.Header, http.Header) error            { return nil }

type testDownstreamMapper struct{}

func (testDownstreamMapper) ToDownstreamStatusCode(int, *int) error                    { return nil }
func (testDownstreamMapper) ToDownstreamHeaders(http.Header, http.Header) error        { return nil }
func (testDownstreamMapper) ToDownstreamResponse(any, *string) error                   { return nil }
func (testDownstreamMapper) ToDownstreamResponseMeta(ResponseMeta, *ServiceMeta) error { return nil }
func (testDownstreamMapper) ToDownstreamInformation(int, RestResponseSpec, *ServiceMeta) error {
	return nil
}
func (testDownstreamMapper) ToDownstreamRedirection(int, RestResponseSpec, *ServiceMeta) error {
	return nil
}
func (testDownstreamMapper) ToDownstreamProblem(int, RestResponseSpec, *ServiceMeta) error {
	return nil
}

var _ RestUpstreamMapper[string, string] = testUpstreamMapper{}
var _ RestDownstreamMapper[string] = testDownstreamMapper{}

func TestRestMapperContracts_InterfaceCompatibility(t *testing.T) {
	var up RestUpstreamMapper[string, string] = testUpstreamMapper{}
	var down RestDownstreamMapper[string] = testDownstreamMapper{}
	_, _ = up, down
}
