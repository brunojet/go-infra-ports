package net_http

import (
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"net/url"

	"github.com/brunojet/go-infra-ports/internal/dtos"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

func register(mux *http.ServeMux, opts handlerOptions, h httpHandlerFuncs) {
	if c := opts.collection; c != nil {
		if c.methods&MethodCreate != 0 {
			registerRoute(mux, "POST /"+c.path, http.HandlerFunc(h.Create), opts)
		}
		if c.methods&MethodList != 0 {
			registerRoute(mux, "GET /"+c.path, http.HandlerFunc(h.List), opts)
		}
	}
	if i := opts.instance; i != nil {
		if i.methods&MethodGet != 0 {
			registerRoute(mux, "GET /"+i.path, http.HandlerFunc(h.Get), opts)
		}
		if i.methods&MethodUpdate != 0 {
			registerRoute(mux, "PUT /"+i.path, http.HandlerFunc(h.Update), opts)
		}
		if i.methods&MethodSave != 0 {
			registerRoute(mux, "PATCH /"+i.path, http.HandlerFunc(h.Save), opts)
		}
		if i.methods&MethodDelete != 0 {
			registerRoute(mux, "DELETE /"+i.path, http.HandlerFunc(h.Delete), opts)
		}
	}
}

func registerRoute(mux *http.ServeMux, pattern string, next http.Handler, opts handlerOptions) {
	mux.Handle(pattern, chainMiddlewares(next, opts.middlewares...))
}

func chainMiddlewares(next http.Handler, middlewares ...Middleware) http.Handler {
	wrapped := next
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}

// validateParams valida os path params da routeEntry contra seus regexes e
// armazena os valores extraídos em ctx.Identifiers. Escreve 400 e retorna
// false se algum param for inválido.
func validateParams(w http.ResponseWriter, r *http.Request, entry *routeEntry, ctx *types.RequestContext) bool {
	if entry == nil || len(entry.params) == 0 {
		return true
	}
	ctx.Identifiers = make(map[string]string, len(entry.params))
	for _, p := range entry.params {
		val := r.PathValue(p.name)
		if p.format != nil && !p.format.MatchString(val) {
			writeBadRequest(w, r, "invalid path parameter: "+p.name)
			return false
		}
		ctx.Identifiers[p.name] = val
	}
	return true
}

func buildRequestContext(w http.ResponseWriter, r *http.Request, entry *routeEntry, ctx *types.RequestContext) bool {
	if ctx == nil {
		writeServiceError(w, r, errors.New("internal error: missing request context"))
		return false
	}
	if !validateParams(w, r, entry, ctx) {
		return false
	}
	// Avoid modifying the original maps in ctx, which may be shared between requests.
	// Ensure the destination maps are allocated with the correct named types
	// before copying values from the request.
	if ctx.Headers == nil {
		ctx.Headers = make(http.Header, len(r.Header))
	}
	if ctx.Query == nil {
		ctx.Query = make(url.Values, len(r.URL.Query()))
	}
	maps.Copy(ctx.Headers, r.Header)
	maps.Copy(ctx.Query, r.URL.Query())
	return true
}

func buildRequestBody[B any](w http.ResponseWriter, r *http.Request, body *B) bool {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeBadRequest(w, r, "invalid body: "+err.Error())
		return false
	}
	return true
}

func writeBody(w http.ResponseWriter, status int, body any) {
	contentType := dtos.ContentTypeJSON
	if ct, ok := body.(dtos.ContentTyper); ok {
		contentType = ct.ContentType()
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	p := DefaultProblemDetails{
		Status:   http.StatusBadRequest,
		Title:    titleFromStatus(http.StatusBadRequest),
		Detail:   message,
		Instance: requestPath(r),
	}
	writeBody(w, http.StatusBadRequest, p)
}

func writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	p := DefaultProblemDetails{
		Status:   http.StatusInternalServerError,
		Title:    titleFromStatus(http.StatusInternalServerError),
		Detail:   err.Error(),
		Instance: requestPath(r),
	}
	errChain := errorChain(err)
	if len(errChain) > 0 {
		p.Extensions = map[string]any{"errors": errChain}
	}
	writeBody(w, http.StatusInternalServerError, p)
}

func (h *netHttpHandler[C, R, U]) writeServiceResponse(w http.ResponseWriter, r *http.Request, status int, meta svccts.ServiceMeta, resp any) {
	switch {
	case status < http.StatusOK: // 1xx Informational
		writeBody(w, status, h.opts.informationBuilder(status, meta))
	case status < http.StatusMultipleChoices: // 2xx Success
		writeBody(w, status, h.opts.envelopeBuilder(status, meta, resp))
	case status < http.StatusBadRequest: // 3xx Redirection
		body, location := h.opts.redirectionBuilder(status, meta)
		if location != "" {
			w.Header().Set("Location", location)
		}
		writeBody(w, status, body)
	default: // 4xx/5xx
		writeBody(w, status, h.opts.problemDetailsBuilder(status, requestPath(r), meta))
	}
}

func problemFromMeta(status int, instance string, meta svccts.ServiceMeta) DefaultProblemDetails {
	p := DefaultProblemDetails{
		Status:   status,
		Title:    firstNonEmpty(meta.Message, titleFromStatus(status)),
		Detail:   meta.Message,
		Instance: firstNonEmpty(metaDetailsString(meta, "instance"), instance),
	}
	if v := metaDetailsString(meta, "problem_type"); v != "" {
		p.Type = v
	}
	if v := metaDetailsString(meta, "problem_title"); v != "" {
		p.Title = v
	}
	ext := make(map[string]any)
	if meta.Code != "" {
		ext["code"] = meta.Code
	}
	if len(meta.Metadata) > 0 {
		ext["details"] = meta.Metadata
	}
	if len(ext) > 0 {
		p.Extensions = ext
	}
	return p
}

func titleFromStatus(status int) string {
	if text := http.StatusText(status); text != "" {
		return text
	}
	return "HTTP Error"
}

func requestPath(r *http.Request) string {
	if r == nil || r.URL == nil {
		return ""
	}
	return r.URL.Path
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func errorChain(err error) []string {
	chain := make([]string, 0, 2)
	for e := err; e != nil; e = errors.Unwrap(e) {
		if msg := e.Error(); msg != "" {
			chain = append(chain, msg)
		}
	}
	return chain
}

func statusOr(code, fallback int) int {
	if code == 0 {
		return fallback
	}
	return code
}

func metaDetailsString(meta svccts.ServiceMeta, key string) string {
	if meta.Metadata == nil {
		return ""
	}
	v, ok := meta.Metadata[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
