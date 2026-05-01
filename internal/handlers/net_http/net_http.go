package net_http

import (
	"net/http"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type netHttpHandler[C, R, U any] struct {
	svc  Service[C, R, U]
	opts handlerOptions
}

func NewNetHttpHandler[C, R, U any](
	mux *http.ServeMux,
	svc Service[C, R, U],
	options ...HandlerOption,
) NetHttpHandler {
	opts := newHandlerOptions(options)
	h := &netHttpHandler[C, R, U]{svc: svc, opts: opts}
	register(mux, opts, h)
	return h
}

func (h *netHttpHandler[C, R, U]) Create(w http.ResponseWriter, r *http.Request) {
	req := ServiceCreate[C]{}
	if ok := buildRequestContext(w, r, h.opts.collection, &req.Context); !ok {
		return
	}
	if ok := buildRequestBody(w, r, &req.Body); !ok {
		return
	}
	var resp ServiceResponse[R]
	if err := h.svc.Create(r.Context(), req, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusCreated), resp.Meta, resp)
}

func (h *netHttpHandler[C, R, U]) List(w http.ResponseWriter, r *http.Request) {
	reqCtx := types.RequestContext{}
	if ok := buildRequestContext(w, r, h.opts.collection, &reqCtx); !ok {
		return
	}
	var resp ServiceResponses[R]
	if err := h.svc.List(r.Context(), reqCtx, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusOK), resp.Meta, resp)
}

func (h *netHttpHandler[C, R, U]) Get(w http.ResponseWriter, r *http.Request) {
	reqCtx := types.RequestContext{}
	if ok := buildRequestContext(w, r, h.opts.instance, &reqCtx); !ok {
		return
	}
	var resp ServiceResponse[R]
	if err := h.svc.Get(r.Context(), reqCtx, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusOK), resp.Meta, resp)
}

func (h *netHttpHandler[C, R, U]) Update(w http.ResponseWriter, r *http.Request) {
	req := ServiceUpdate[U]{}
	if ok := buildRequestContext(w, r, h.opts.instance, &req.Context); !ok {
		return
	}
	if ok := buildRequestBody(w, r, &req.Body); !ok {
		return
	}
	var resp ServiceResponse[R]
	if err := h.svc.Update(r.Context(), req, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusOK), resp.Meta, resp)
}

func (h *netHttpHandler[C, R, U]) Save(w http.ResponseWriter, r *http.Request) {
	req := ServiceSave[C]{}
	if ok := buildRequestContext(w, r, h.opts.instance, &req.Context); !ok {
		return
	}
	if ok := buildRequestBody(w, r, &req.Body); !ok {
		return
	}
	var resp ServiceResponse[R]
	if err := h.svc.Save(r.Context(), req, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusOK), resp.Meta, resp)
}

func (h *netHttpHandler[C, R, U]) Delete(w http.ResponseWriter, r *http.Request) {
	reqCtx := types.RequestContext{}
	if ok := buildRequestContext(w, r, h.opts.instance, &reqCtx); !ok {
		return
	}
	var resp ServiceResponse[R]
	if err := h.svc.Delete(r.Context(), reqCtx, &resp); err != nil {
		writeServiceError(w, r, err)
		return
	}
	h.writeServiceResponse(w, r, statusOr(resp.Context.StatusCode, http.StatusNoContent), resp.Meta, resp)
}
