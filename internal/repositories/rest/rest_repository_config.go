package rest

import (
	"fmt"
	"net/http"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

// PathMethod selects which HTTP verbs are registered for a route entry.
type PathMethod uint8

const (
	MethodCreate PathMethod = 1 << iota // POST  — collection
	MethodList                          // GET   — collection
	MethodGet                           // GET   — instance
	MethodUpdate                        // PUT   — instance
	MethodSave                          // PATCH — instance
	MethodDelete                        // DELETE — instance

	AllCollectionMethods PathMethod = MethodCreate | MethodList
	AllInstanceMethods   PathMethod = MethodGet | MethodUpdate | MethodSave | MethodDelete
)

var (
	DefaultCollectionPathEntry = &pathEntry{templateFmt: "/"}
	DefaultInstancePathEntry   = &pathEntry{templateFmt: "/%s", paramNames: []string{"id"}}
)

// RepositoryOption configures restRepository construction.
type RepositoryOption func(*repositoryOptions) error

// pathEntry stores a URL path template and the identifier names it requires.
type pathEntry struct {
	templateFmt string
	paramNames  []string
}

// fullPath constructs the full URL path by substituting args into the templateFmt.
func (e *pathEntry) fullPath(basePath string, ids types.Identifiers) (string, error) {
	if len(ids) != len(e.paramNames) {
		return "", fmt.Errorf("expected %d identifiers for path template, got %d", len(e.paramNames), len(ids))
	}
	if len(e.paramNames) == 0 {
		return basePath + e.templateFmt, nil // path sem parâmetros
	}
	anyArgs := make([]any, len(ids))
	for i, name := range e.paramNames {
		if id, ok := ids[name]; ok {
			anyArgs[i] = id
		} else {
			return "", fmt.Errorf("missing identifier %q for path template", name)
		}
	}
	return basePath + fmt.Sprintf(e.templateFmt, anyArgs...), nil
}

type repositoryOptions struct {
	client   HttpClient
	registry RestRegistry
	basePath string
	paths    map[PathMethod]*pathEntry
	headers  http.Header
}

func newRepositoryOptions(opts ...RepositoryOption) *repositoryOptions {
	o := &repositoryOptions{
		paths: map[PathMethod]*pathEntry{
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
		if err := opt(o); err != nil {
			panic(err.Error())
		}
	}
	if o.client == nil {
		panic(errRepositoryMissingHttpClient.Error())
	}
	return o
}

// WithHttpClient sets the HttpClient used for HTTP transport. Required.
func WithHttpClient(c HttpClient) RepositoryOption {
	return func(o *repositoryOptions) error {
		if c == nil {
			panic(errRepositoryMissingHttpClient.Error())
		}
		o.client = c
		return nil
	}
}

// WithRegistry sets the RestRegistry used for request/response marshaling.
func WithRegistry(r RestRegistry) RepositoryOption {
	return func(o *repositoryOptions) error {
		if r == nil {
			panic(errRepositoryMissingRegistry.Error())
		}
		o.registry = r
		return nil
	}
}

// WithBasePath sets the path prefix common to all operations (e.g. "/api/v1/users").
// The HttpClient is responsible for schema+host+port.
func WithBasePath(basePath string) RepositoryOption {
	return func(o *repositoryOptions) error {
		p, err := sanitizeAndValidatePath(basePath)
		if err != nil {
			panic(errRepositoryInvalidBasePath(err).Error())
		}
		o.basePath = p
		return nil
	}
}

// WithPath registers pathTemplate for the given methods.
func WithPath(methods PathMethod, pathTemplate string) RepositoryOption {
	return func(o *repositoryOptions) error {
		templateFmt, paramNames, err := extractPathParams(pathTemplate)
		if err != nil {
			panic(err.Error())
		}
		entry := &pathEntry{templateFmt: templateFmt, paramNames: paramNames}
		for _, m := range []PathMethod{MethodCreate, MethodList, MethodGet, MethodUpdate, MethodSave, MethodDelete} {
			if methods&m != 0 {
				o.paths[m] = entry
			}
		}
		return nil
	}
}

// WithHeader adds a default request header sent on every operation.
func WithHeader(key, value string) RepositoryOption {
	return func(o *repositoryOptions) error {
		o.headers.Set(key, value)
		return nil
	}
}
