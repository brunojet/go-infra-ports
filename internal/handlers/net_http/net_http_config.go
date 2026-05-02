package net_http

import (
	"net/http"
	"regexp"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
)

// RouteMethod selects which HTTP verbs are registered for a route entry.
type RouteMethod uint8

const (
	MethodCreate RouteMethod = 1 << iota // POST  — collection
	MethodList                           // GET   — collection
	MethodGet                            // GET   — instance
	MethodUpdate                         // PUT   — instance
	MethodSave                           // PATCH — instance
	MethodDelete                         // DELETE — instance

	// AllCollectionMethods registers Create (POST) and List (GET).
	AllCollectionMethods RouteMethod = MethodCreate | MethodList
	// AllInstanceMethods registers Get, Update, Save and Delete.
	AllInstanceMethods RouteMethod = MethodGet | MethodUpdate | MethodSave | MethodDelete
)

type paramFormat struct {
	name   string
	format *regexp.Regexp // nil = sem validação
}

type routeEntry struct {
	path    string
	params  []paramFormat
	methods RouteMethod
}

type Middleware func(http.Handler) http.Handler

type handlerOptions struct {
	collection            *routeEntry
	instance              *routeEntry
	middlewares           []Middleware
	informationBuilder    func(status int, meta svccts.ServiceMeta) any
	envelopeBuilder       func(status int, meta svccts.ServiceMeta, resp any) any
	redirectionBuilder    func(status int, meta svccts.ServiceMeta) (body any, location string)
	problemDetailsBuilder func(status int, instance string, meta svccts.ServiceMeta) any
}

type HandlerOption func(*handlerOptions)

func newHandlerOptions(options []HandlerOption) handlerOptions {
	opts := handlerOptions{
		informationBuilder: func(status int, meta svccts.ServiceMeta) any {
			return DefaultInformation{Code: status, Message: firstNonEmpty(meta.Message, titleFromStatus(status))}
		},
		envelopeBuilder: func(_ int, _ svccts.ServiceMeta, resp any) any {
			return resp
		},
		redirectionBuilder: func(status int, meta svccts.ServiceMeta) (any, string) {
			return DefaultRedirection{Code: status, Location: meta.Location}, meta.Location
		},
		problemDetailsBuilder: func(status int, instance string, meta svccts.ServiceMeta) any {
			return problemFromMeta(status, instance, meta)
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return opts
}

// httpHandlerFuncs is satisfied by *NetHttpHandler — used to break the import cycle
// between config and the handler struct.
type httpHandlerFuncs interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// WithCollection registra os métodos selecionados em methods para o pathFmt base.
// parentsValidation são regexes posicionais para os {params} extraídos do pathFmt.
func WithCollection(methods RouteMethod, pathFmt string, parentsValidation ...*regexp.Regexp) HandlerOption {
	p := sanitizePath(pathFmt)
	return func(o *handlerOptions) {
		o.collection = &routeEntry{path: p, params: extractParams(p, parentsValidation), methods: methods}
	}
}

// WithInstance registra os métodos selecionados em methods para o pathFmt de instância.
// parentsValidation são regexes posicionais para os {params} extraídos do pathFmt.
func WithInstance(methods RouteMethod, pathFmt string, parentsValidation ...*regexp.Regexp) HandlerOption {
	p := sanitizePath(pathFmt)
	return func(o *handlerOptions) {
		o.instance = &routeEntry{path: p, params: extractParams(p, parentsValidation), methods: methods}
	}
}

// WithMiddlewareRegistry appends the middleware chain currently registered in registry.
// The chain is copied at configuration time so later registry mutations do not affect registered routes.
func WithMiddlewareRegistry(registry *MiddlewareRegistry) HandlerOption {
	if registry == nil {
		panic("net_http: WithMiddlewareRegistry: registry must not be nil")
	}
	validated := registry.Middlewares()
	return func(o *handlerOptions) {
		o.middlewares = append(o.middlewares, validated...)
	}
}

// WithInformation substitui o DTO padrão para respostas 1xx.
func WithInformation(fn func(status int, meta svccts.ServiceMeta) any) HandlerOption {
	if fn == nil {
		panic("net_http: WithInformation: fn must not be nil")
	}
	return func(o *handlerOptions) { o.informationBuilder = fn }
}

// WithEnvelope envolve o payload de respostas 2xx num body customizado.
// Útil para wrappers como {"data": ...} ou metadados de paginação em 206.
func WithEnvelope(fn func(status int, meta svccts.ServiceMeta, resp any) any) HandlerOption {
	if fn == nil {
		panic("net_http: WithEnvelope: fn must not be nil")
	}
	return func(o *handlerOptions) { o.envelopeBuilder = fn }
}

// WithRedirection substitui o DTO e o header Location para respostas 3xx.
// fn retorna (body, location); location vazio omite o header.
func WithRedirection(fn func(status int, meta svccts.ServiceMeta) (body any, location string)) HandlerOption {
	if fn == nil {
		panic("net_http: WithRedirection: fn must not be nil")
	}
	return func(o *handlerOptions) { o.redirectionBuilder = fn }
}

// WithProblemDetails substitui o DTO padrão RFC 7807 para respostas 4xx/5xx.
func WithProblemDetails(fn func(status int, instance string, meta svccts.ServiceMeta) any) HandlerOption {
	if fn == nil {
		panic("net_http: WithProblemDetails: fn must not be nil")
	}
	return func(o *handlerOptions) { o.problemDetailsBuilder = fn }
}
