package http_helper

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestApplyQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		input   url.Values
		wantQ   map[string][]string
		wantErr error
	}{
		{
			name:    "nil url",
			wantErr: errNilURL,
		},
		{
			name:   "empty q does nothing",
			rawURL: "http://example.com/",
			input:  url.Values{},
		},
		{
			name:   "append and dedupe",
			rawURL: "http://example.com/?a=1&a=2",
			input:  url.Values{"a": {"2", "3"}, "b": {"x"}},
			wantQ: map[string][]string{
				"a": {"1", "2", "3"},
				"b": {"x"},
			},
		},
		{
			name:   "sort and dedupe",
			rawURL: "http://example.com/?a=3&a=1&a=3",
			input:  url.Values{"a": {"2", "1"}, "b": {"z", "y"}},
			wantQ: map[string][]string{
				"a": {"1", "2", "3"},
				"b": {"y", "z"},
			},
		},
		{
			name:   "empty url no existing query",
			rawURL: "http://example.com/",
			input:  url.Values{"a": {"1"}},
			wantQ:  map[string][]string{"a": {"1"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *url.URL
			if tt.rawURL != "" {
				u, _ = url.Parse(tt.rawURL)
			}
			err := ApplyURLQueryParams(u, tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantQ != nil {
				vals := u.Query()
				for k, want := range tt.wantQ {
					if !reflect.DeepEqual(vals[k], want) {
						t.Fatalf("key %q: expected %v, got %v", k, want, vals[k])
					}
				}
			}
		})
	}
}

func TestApplyHeaderParams(t *testing.T) {
	tests := []struct {
		name    string
		dst     http.Header
		src     http.Header
		wantDst http.Header
		wantErr error
	}{
		{
			name:    "nil dst",
			dst:     nil,
			src:     http.Header{"X-Test": {"v1"}},
			wantErr: errNilDstHeader,
		},
		{
			name:    "empty src leaves dst unchanged",
			dst:     http.Header{"K": {"v"}},
			src:     http.Header{},
			wantDst: http.Header{"K": {"v"}},
		},
		{
			name:    "copies src into dst",
			dst:     http.Header{},
			src:     http.Header{"X-Test": {"v1", "v2"}},
			wantDst: http.Header{"X-Test": {"v1", "v2"}},
		},
		{
			name:    "src overwrites dst on conflict",
			dst:     http.Header{"X-Test": {"old"}},
			src:     http.Header{"X-Test": {"new"}},
			wantDst: http.Header{"X-Test": {"new"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyHeaderParams(tt.dst, tt.src)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantDst != nil && !reflect.DeepEqual(tt.dst, tt.wantDst) {
				t.Fatalf("expected dst %v, got %v", tt.wantDst, tt.dst)
			}
		})
	}
}
