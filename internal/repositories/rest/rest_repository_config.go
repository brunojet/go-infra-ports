package rest

import (
	"fmt"
	"net/http"

	"github.com/brunojet/go-infra-ports/internal/helpers/http_helper"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

var (
	DefaultCollectionPathEntry = &pathEntry{templateFmt: "/"}
	DefaultInstancePathEntry   = &pathEntry{templateFmt: "/%s", paramNames: []string{"id"}}
)

// RepositoryOption configures restRepository construction.
type RepositoryOption func(*repositoryOptions)

// pathEntry stores a URL path template and the identifier names it requires.
type pathEntry struct {
	templateFmt string
	paramNames  []string
}

// expandPath constructs the URL path by substituting identifiers into the templateFmt.
// Note: it returns only the path portion (caller should prepend `basePath`).
func (e *pathEntry) expandPath(ids types.Identifiers) (string, error) {
	if len(e.paramNames) == 0 {
		return e.templateFmt, nil // path sem parâmetros
	}
	anyArgs := make([]any, len(e.paramNames))
	for i, name := range e.paramNames {
		if id, ok := ids[name]; ok {
			anyArgs[i] = id
		} else {
			return "", fmt.Errorf("identifiers: missing identifier %q for path template", name)
		}
	}
	return fmt.Sprintf(e.templateFmt, anyArgs...), nil
}

type repositoryOptions struct {
	client   HttpClient
	registry RestRegistry
	basePath string
	paths    map[RestMethod]*pathEntry
	headers  http.Header
}

func newRepositoryOptions(opts ...RepositoryOption) *repositoryOptions {
	o := &repositoryOptions{
		paths: map[RestMethod]*pathEntry{
			MethodCreate: DefaultCollectionPathEntry,
			MethodList:   DefaultCollectionPathEntry,
			MethodGet:    DefaultInstancePathEntry,
			MethodUpdate: DefaultInstancePathEntry,
			MethodSave:   DefaultInstancePathEntry,
			MethodDelete: DefaultInstancePathEntry,
		},
		registry: NewRestRegistry(),
		headers:  make(http.Header),
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.client == nil {
		panic(errRepositoryMissingHttpClient.Error())
	}
	return o
}

// WithHttpClient sets the HttpClient used for HTTP transport. Required.
func WithHttpClient(c HttpClient) RepositoryOption {
	return func(o *repositoryOptions) {
		if c == nil {
			panic(errRepositoryMissingHttpClient.Error())
		}
		o.client = c
	}
}

// WithRegistry sets the RestRegistry used for request/response marshaling.
func WithRegistry(r RestRegistry) RepositoryOption {
	return func(o *repositoryOptions) {
		if r == nil {
			panic(errRepositoryMissingRegistry.Error())
		}
		o.registry = r
	}
}

// WithBasePath sets the path prefix common to all operations (e.g. "/api/v1/users").
// The HttpClient is responsible for schema+host+port.
func WithBasePath(basePath string) RepositoryOption {
	return func(o *repositoryOptions) {
		p, err := http_helper.SanitizeAndValidatePath(basePath)
		if err != nil {
			panic(errRepositoryInvalidBasePath(err).Error())
		}
		o.basePath = p
	}
}

// WithPath registers pathTemplate for the given methods.
func WithPath(methods RestMethod, pathTemplate string) RepositoryOption {
	return func(o *repositoryOptions) {
		// Validate template and extract parameter names
		paramNames, err := http_helper.ExtractPathParams(pathTemplate)
		if err != nil {
			panic(err.Error())
		}
		// Build format string with %s placeholders for interpolation
		templateFmt, err := http_helper.PathParamsFmt(pathTemplate)
		if err != nil {
			panic(err.Error())
		}
		entry := &pathEntry{templateFmt: templateFmt, paramNames: paramNames}
		for _, m := range []RestMethod{MethodCreate, MethodList, MethodGet, MethodUpdate, MethodSave, MethodDelete} {
			if methods&m != 0 {
				o.paths[m] = entry
			}
		}
	}
}

// WithHeader adds a default request header sent on every operation.
func WithHeader(key, value string) RepositoryOption {
	return func(o *repositoryOptions) {
		o.headers.Set(key, value)
	}
}
