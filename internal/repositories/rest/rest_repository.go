package rest

import (
	"context"
	"net/http"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// NewRestRepository builds a RestRepository with the provided options.
// HttpClient is required; panics at startup if missing or if any option is invalid.
func NewRestRepository(opts ...RepositoryOption) RestRepository {
	o := newRepositoryOptions(opts...)
	return &restRepository{
		client:   o.client,
		registry: o.registry,
		opts:     o,
	}
}

func (r *restRepository) Create(ctx context.Context, request RestRequest, response *RestResponse) error {
	resp, err := r.executeBodyRequest(ctx, MethodCreate, http.MethodPost, request) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponse(resp, response)
}

func (r *restRepository) List(ctx context.Context, reqCtx types.RequestContext, response *RestResponses) error {
	resp, err := r.executeNoBodyRequest(ctx, MethodList, http.MethodGet, reqCtx) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponses(resp, response)
}

func (r *restRepository) Get(ctx context.Context, reqCtx types.RequestContext, response *RestResponse) error {
	resp, err := r.executeNoBodyRequest(ctx, MethodGet, http.MethodGet, reqCtx) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponse(resp, response)
}

func (r *restRepository) Update(ctx context.Context, request RestRequest, response *RestResponse) error {
	resp, err := r.executeBodyRequest(ctx, MethodUpdate, http.MethodPut, request) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponse(resp, response)
}

func (r *restRepository) Save(ctx context.Context, request RestRequest, response *RestResponse) error {
	resp, err := r.executeBodyRequest(ctx, MethodSave, http.MethodPatch, request) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponse(resp, response)
}

func (r *restRepository) Delete(ctx context.Context, reqCtx types.RequestContext, response *RestResponse) error {
	resp, err := r.executeNoBodyRequest(ctx, MethodDelete, http.MethodDelete, reqCtx) //nolint:bodyclose // closed inside readBody
	if err != nil {
		return err
	}
	return r.resolveResponse(resp, response)
}
