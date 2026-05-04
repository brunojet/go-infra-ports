package contracts

import (
	"context"
	"net/http"
	"testing"
)

type testHttpClient struct{}

func (testHttpClient) Do(_ context.Context, _ *http.Request) (*http.Response, error) {
	return nil, nil
}

var _ HttpClient = testHttpClient{}

func TestHttpClient_InterfaceSatisfied(t *testing.T) {
	var _ HttpClient = testHttpClient{}
}
