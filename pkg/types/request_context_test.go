package types

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRequestAndResponseContext_BasicFields(t *testing.T) {
	req := RequestContext{
		Query:       url.Values{"q": {"search"}},
		Headers:     http.Header{"X-Test": {"1"}},
		Identifiers: Identifiers{"id": "42"},
	}

	resp := ResponseContext{
		StatusCode: 201,
		Headers:    http.Header{"Location": {"/items/42"}},
		Meta:       ResponseMeta{"page": 1},
	}

	if req.Query.Get("q") != "search" || req.Headers.Get("X-Test") != "1" || req.Identifiers["id"] != "42" {
		t.Fatal("unexpected request context values")
	}
	if resp.StatusCode != 201 || resp.Headers.Get("Location") != "/items/42" {
		t.Fatal("unexpected response context values")
	}
	if resp.Meta["page"] != 1 {
		t.Fatal("unexpected response meta value")
	}
}
