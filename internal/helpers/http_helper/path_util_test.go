package http_helper

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestSanitizeAndValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "empty",
			input:   "",
			wantErr: errRepositoryBasePathEmpty,
		},
		{
			name:    "path too long",
			input:   "/" + strings.Repeat("a", maxPathLen),
			wantErr: errPathTooLong,
		},
		{
			name:    "invalid chars query string",
			input:   "/users?bad=1",
			wantErr: errPathInvalidChars,
		},
		{
			name:    "invalid chars space in middle",
			input:   "/users /list",
			wantErr: errPathInvalidChars,
		},
		{
			name:    "invalid chars colon style",
			input:   "/users/:id",
			wantErr: errPathInvalidChars,
		},
		{
			name:    "invalid structure unclosed brace",
			input:   "/users/{id",
			wantErr: errPathInvalidStructure,
		},
		{
			name:  "normalizes trailing slash",
			input: "/api/v1/users/",
			want:  "/api/v1/users",
		},
		{
			name:  "normalizes spaces via TrimSpace",
			input: " /api/v1/users ",
			want:  "/api/v1/users",
		},
		{
			name:  "root path",
			input: "/",
			want:  "/",
		},
		{
			name:  "with param",
			input: "/users/{id}",
			want:  "/users/{id}",
		},
		{
			name:  "multiple params",
			input: "/users/{id}/orders/{orderId}",
			want:  "/users/{id}/orders/{orderId}",
		},
		{
			name:  "without leading slash",
			input: "users/{id}/orders/{orderId}",
			want:  "/users/{id}/orders/{orderId}",
		},
		{
			name:  "with trailing slash",
			input: "/users/{id}/orders/{orderId}/",
			want:  "/users/{id}/orders/{orderId}",
		},
		{
			name:  "messy slashes",
			input: "//users/{id}////orders///{orderId}////",
			want:  "/users/{id}/orders/{orderId}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeAndValidatePath(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestPathParamsFmt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "invalid template propagates",
			input:   "/users/{id",
			wantErr: errInvalidPathTemplateErr,
		},
		{
			name:  "no params",
			input: "/api/v1/static",
			want:  "/api/v1/static",
		},
		{
			name:  "single param",
			input: "/users/{id}",
			want:  "/users/%s",
		},
		{
			name:  "multiple params",
			input: "/users/{id}/orders/{orderId}",
			want:  "/users/%s/orders/%s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathParamsFmt(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestExtractPathParams(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNames []string
		wantErr   error
	}{
		{
			name:    "invalid template propagates",
			input:   "/users/{id",
			wantErr: errInvalidPathTemplateErr,
		},
		{
			name:    "too many params",
			input:   "/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}/{i}",
			wantErr: errPathParametersExceedMax,
		},
		{
			name:  "no params",
			input: "/api/v1/static",
		},
		{
			name:      "single param",
			input:     "/users/{id}",
			wantNames: []string{"id"},
		},
		{
			name:      "multiple params",
			input:     "/users/{id}/orders/{orderId}",
			wantNames: []string{"id", "orderId"},
		},
		{
			name:      "max params exactly",
			input:     "/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}",
			wantNames: []string{"a", "b", "c", "d", "e", "f", "g", "h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNames, err := ExtractPathParams(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Fatalf("expected names %v, got %v", tt.wantNames, gotNames)
			}
		})
	}
}
