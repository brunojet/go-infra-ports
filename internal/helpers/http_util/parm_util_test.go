package http_util

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestApplyQueryParams_Append(t *testing.T) {
	u, _ := url.Parse("http://example.com/?a=1&a=2")
	q := url.Values{"a": {"2", "3"}, "b": {"x"}}
	ApplyQueryParams(u, q)

	vals := u.Query()
	if !reflect.DeepEqual(vals["a"], []string{"1", "2", "3"}) {
		t.Fatalf("unexpected a values: %v", vals["a"])
	}
	if !reflect.DeepEqual(vals["b"], []string{"x"}) {
		t.Fatalf("unexpected b values: %v", vals["b"])
	}

	// empty q does nothing
	u2, _ := url.Parse("http://example.com/")
	ApplyQueryParams(u2, url.Values{})
	if u2.RawQuery != "" {
		t.Fatalf("expected empty RawQuery, got %q", u2.RawQuery)
	}
}

func TestApplyHeaderParams_CopyBehavior(t *testing.T) {
	dst := make(http.Header)
	src := http.Header{"X-Test": {"v1", "v2"}}
	ApplyHeaderParams(dst, src)
	if !reflect.DeepEqual(dst["X-Test"], src["X-Test"]) {
		t.Fatalf("expected headers copied, got %v", dst)
	}

	// empty src leaves dst unchanged
	dst2 := http.Header{"K": {"v"}}
	ApplyHeaderParams(dst2, http.Header{})
	if !reflect.DeepEqual(dst2["K"], []string{"v"}) {
		t.Fatalf("expected dst unchanged, got %v", dst2)
	}
}

func TestApplyQueryParams_SortAndDedupe(t *testing.T) {
	u, _ := url.Parse("http://example.com/?a=3&a=1&a=3")
	q := url.Values{"a": {"2", "1"}, "b": {"z", "y"}}
	ApplyQueryParams(u, q)

	vals := u.Query()
	if !reflect.DeepEqual(vals["a"], []string{"1", "2", "3"}) {
		t.Fatalf("unexpected a values: %v", vals["a"])
	}
	if !reflect.DeepEqual(vals["b"], []string{"y", "z"}) {
		t.Fatalf("unexpected b values: %v", vals["b"])
	}
}
