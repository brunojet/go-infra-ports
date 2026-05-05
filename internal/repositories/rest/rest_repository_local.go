package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	http_util "github.com/brunojet/go-infra-ports/internal/helpers/http_util"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func (r *restRepository) resolveURL(pathMethod RestMethod, ids types.Identifiers, query url.Values) (string, error) {
	entry, ok := r.opts.paths[pathMethod]
	if !ok {
		return "", errRepositoryPathMethodNotConfiguredf(pathMethod)
	}
	pathOnly, err := entry.expandPath(ids)
	if err != nil {
		return "", err
	}
	rawURL := r.opts.basePath + pathOnly
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errRepositoryBuildRequest(err)
	}
	http_util.ApplyQueryParams(u, query)
	return u.String(), nil
}

func (r *restRepository) buildHTTPRequest(ctx context.Context, method, rawURL string, body []byte, reqHeaders http.Header) (*http.Request, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, errRepositoryBuildRequest(err)
	}
	req.Header = r.applyHeaderParams(reqHeaders, len(body) > 0)
	return req, nil
}

// applyHeaderParams merges config-level and request-level headers, with request taking precedence.
// Sets Content-Type: application/json when hasBody is true and the caller has not already provided one.
func (r *restRepository) applyHeaderParams(reqHeaders http.Header, hasBody bool) http.Header {
	merged := make(http.Header, len(r.opts.headers)+len(reqHeaders))
	http_util.ApplyHeaderParams(merged, r.opts.headers)
	http_util.ApplyHeaderParams(merged, reqHeaders)
	if hasBody && merged.Get("Content-Type") == "" {
		merged.Set("Content-Type", "application/json")
	}
	return merged
}

func (r *restRepository) executeRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := r.client.Do(ctx, req)
	if err != nil {
		closeBody(resp)
		return nil, errRepositoryExecuteRequest(err)
	}
	if resp == nil {
		return nil, errRepositoryExecuteRequest(errRepositoryNilHTTPResponse)
	}
	return resp, nil
}

func (r *restRepository) readBody(resp *http.Response) ([]byte, error) {
	defer closeBody(resp)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errRepositoryReadResponseBody(err)
	}
	return data, nil
}

func (r *restRepository) mapResponse(status int, body []byte, out *RestResponse) error {
	out.Context.StatusCode = status
	if err := r.registry.ResolveEnvelopeResponse(status, &body, &out.Context.Meta); err != nil {
		return errRepositoryResolveEnvelopeResponse(err)
	}
	if err := r.registry.ResolveResponse(status, body, &out.Data); err != nil {
		return errRepositoryResolveResponse(err)
	}
	return nil
}

func (r *restRepository) mapResponses(status int, body []byte, out *RestResponses) error {
	out.Context.StatusCode = status
	if err := r.registry.ResolveEnvelopeResponse(status, &body, &out.Context.Meta); err != nil {
		return errRepositoryResolveEnvelopeResponse(err)
	}
	if err := r.registry.ResolveResponses(status, body, &out.Data); err != nil {
		return errRepositoryResolveResponse(err)
	}
	return nil
}

// executeBodyRequest is the shared flow for Create, Update, and Save operations.
func (r *restRepository) executeBodyRequest(ctx context.Context, restMethod RestMethod, method string, request RestRequest) (*http.Response, error) {
	if request.Body == nil {
		return nil, errRepositoryBuildRequest(errRepositoryRequestBodyNilf(restMethod))
	}
	var body []byte
	if err := r.registry.ResolveRequest(request.Body, &body); err != nil {
		return nil, errRepositoryResolveRequest(err)
	}
	if err := r.registry.ResolveEnvelopeRequest(restMethod, &body); err != nil {
		return nil, errRepositoryResolveEnvelopeRequest(err)
	}
	rawURL, err := r.resolveURL(restMethod, request.Context.Identifiers, request.Context.Query)
	if err != nil {
		return nil, err
	}
	req, err := r.buildHTTPRequest(ctx, method, rawURL, body, request.Context.Headers)
	if err != nil {
		return nil, err
	}
	return r.executeRequest(ctx, req)
}

func (r *restRepository) executeNoBodyRequest(ctx context.Context, restMethod RestMethod, method string, reqCtx types.RequestContext) (*http.Response, error) {
	rawURL, err := r.resolveURL(restMethod, reqCtx.Identifiers, reqCtx.Query)
	if err != nil {
		return nil, err
	}
	req, err := r.buildHTTPRequest(ctx, method, rawURL, nil, reqCtx.Headers)
	if err != nil {
		return nil, err
	}
	return r.executeRequest(ctx, req) //nolint:bodyclose // closed inside readBody

}

func (r *restRepository) resolveResponse(httpResponse *http.Response, out *RestResponse) error {
	out.Context.StatusCode = httpResponse.StatusCode
	out.Context.Headers = httpResponse.Header
	respBody, err := r.readBody(httpResponse)
	if err != nil {
		return err
	}
	return r.mapResponse(httpResponse.StatusCode, respBody, out)
}

func (r *restRepository) resolveResponses(httpResponse *http.Response, out *RestResponses) error {
	out.Context.StatusCode = httpResponse.StatusCode
	out.Context.Headers = httpResponse.Header
	respBody, err := r.readBody(httpResponse)
	if err != nil {
		return err
	}
	return r.mapResponses(httpResponse.StatusCode, respBody, out)
}
