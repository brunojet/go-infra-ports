package rest

import (
	"errors"
	"strings"
	"testing"
)

func TestRepositoryErrorSentinels_AreDistinct(t *testing.T) {
	errs := []error{
		errRepositoryMissingHttpClient,
		errRepositoryMissingRegistry,
		errRepositoryBaseURLEmpty,
		errRepositoryPathInvalidChars,
		errRepositoryPathInvalidStructure,
		errRepositoryPathMethodNotConfigured,
		errRepositoryRequestBodyNil,
		errRepositoryNilHTTPResponse,
		errRepositoryCollectionPathHasIDs,
		errRepositoryInstancePathMissingID,
	}
	for i, a := range errs {
		for j, b := range errs {
			if i != j && errors.Is(a, b) {
				t.Fatalf("sentinel errors %d and %d should be distinct", i, j)
			}
		}
	}
}

func TestRepositoryErrorWrappers_WrapOriginalError(t *testing.T) {
	base := errors.New("base-error")
	cases := []struct {
		name string
		fn   func(error) error
		want string
	}{
		{"build request", errRepositoryBuildRequest, "build request"},
		{"execute request", errRepositoryExecuteRequest, "execute request"},
		{"read response body", errRepositoryReadResponseBody, "read response body"},
		{"resolve request", errRepositoryResolveRequest, "resolve request"},
		{"resolve envelope request", errRepositoryResolveEnvelopeRequest, "resolve envelope request"},
		{"resolve envelope response", errRepositoryResolveEnvelopeResponse, "resolve envelope response"},
		{"resolve response", errRepositoryResolveResponse, "resolve response"},
		{"invalid basePath", errRepositoryInvalidBasePath, "invalid basePath"},
		{"invalid path template", errRepositoryInvalidPathTemplate, "invalid path template"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(base)
			if !errors.Is(err, base) {
				t.Fatalf("expected wrapped error to contain base: %v", err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q in error %q", tc.want, err.Error())
			}
		})
	}
}

func TestRepositoryPathFormatErrors_WrapSentinels(t *testing.T) {
	errChars := errRepositoryPathInvalidCharsf("/users/{id}?bad")
	if !errors.Is(errChars, errRepositoryPathInvalidChars) {
		t.Fatalf("expected errRepositoryPathInvalidChars sentinel, got: %v", errChars)
	}
	if !strings.Contains(errChars.Error(), "/users/{id}?bad") {
		t.Fatalf("expected offending path in error message, got: %q", errChars.Error())
	}

	errStructure := errRepositoryPathInvalidStructuref("/users//id")
	if !errors.Is(errStructure, errRepositoryPathInvalidStructure) {
		t.Fatalf("expected errRepositoryPathInvalidStructure sentinel, got: %v", errStructure)
	}
	if !strings.Contains(errStructure.Error(), "/users//id") {
		t.Fatalf("expected offending path in error message, got: %q", errStructure.Error())
	}
}

func TestRepositoryLocalFormatErrors_WrapSentinels(t *testing.T) {
	errPathMethod := errRepositoryPathMethodNotConfiguredf(MethodGet)
	if !errors.Is(errPathMethod, errRepositoryPathMethodNotConfigured) {
		t.Fatalf("expected errRepositoryPathMethodNotConfigured sentinel, got: %v", errPathMethod)
	}

	errBodyNil := errRepositoryRequestBodyNilf(MethodCreate)
	if !errors.Is(errBodyNil, errRepositoryRequestBodyNil) {
		t.Fatalf("expected errRepositoryRequestBodyNil sentinel, got: %v", errBodyNil)
	}
}
