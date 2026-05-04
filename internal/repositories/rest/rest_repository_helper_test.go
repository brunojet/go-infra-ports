package rest

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestSanitizeAndValidatePath_Valid_Sanitizes(t *testing.T) {
	got, err := sanitizeAndValidatePath(" /api/v1/users/ ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/api/v1/users" {
		t.Fatalf("expected /api/v1/users, got %q", got)
	}
}

func TestSanitizeAndValidatePath_Empty_ReturnsBaseURLEmpty(t *testing.T) {
	_, err := sanitizeAndValidatePath("  ")
	if !errors.Is(err, errRepositoryBaseURLEmpty) {
		t.Fatalf("expected errRepositoryBaseURLEmpty, got %v", err)
	}
}

func TestSanitizeAndValidatePath_InvalidChars_ReturnsSentinel(t *testing.T) {
	_, err := sanitizeAndValidatePath("/users?bad=1")
	if !errors.Is(err, errRepositoryPathInvalidChars) {
		t.Fatalf("expected errRepositoryPathInvalidChars, got %v", err)
	}
}

func TestSanitizeAndValidatePath_InvalidStructure_ReturnsSentinel(t *testing.T) {
	_, err := sanitizeAndValidatePath("/users/:")
	if !errors.Is(err, errRepositoryPathInvalidStructure) {
		t.Fatalf("expected errRepositoryPathInvalidStructure, got %v", err)
	}
}

func TestExtractPathParams_WithParams_ReturnsFormatAndNames(t *testing.T) {
	format, names, err := extractPathParams("/users/{id}/orders/:orderId")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format != "/users/%s/orders/%s" {
		t.Fatalf("expected /users/%%s/orders/%%s, got %q", format)
	}
	if len(names) != 2 || names[0] != "id" || names[1] != "orderId" {
		t.Fatalf("expected names [id orderId], got %v", names)
	}
}

func TestExtractPathParams_WithoutParams_ReturnsSanitizedPath(t *testing.T) {
	format, names, err := extractPathParams("/users/list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if format != "/users/list" {
		t.Fatalf("expected /users/list, got %q", format)
	}
	if names != nil {
		t.Fatalf("expected nil names, got %v", names)
	}
}

func TestExtractPathParams_InvalidTemplate_WrapsUnderlyingError(t *testing.T) {
	_, _, err := extractPathParams("/users?bad=1")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
	if !errors.Is(err, errRepositoryPathInvalidChars) {
		t.Fatalf("expected wrapped errRepositoryPathInvalidChars, got %v", err)
	}
}

func TestApplyQueryParams_AppendsParams(t *testing.T) {
	u, _ := url.Parse("https://api.example.com/items?existing=1")
	q := url.Values{"page": {"2"}, "limit": {"10"}}
	applyQueryParams(u, q)

	got := u.Query()
	if got.Get("existing") != "1" {
		t.Fatalf("expected existing=1 to be preserved")
	}
	if got.Get("page") != "2" {
		t.Fatalf("expected page=2, got %q", got.Get("page"))
	}
	if got.Get("limit") != "10" {
		t.Fatalf("expected limit=10, got %q", got.Get("limit"))
	}
}

func TestApplyQueryParams_EmptyQueryIsNoOp(t *testing.T) {
	u, _ := url.Parse("https://api.example.com/items?existing=1")
	original := u.RawQuery
	applyQueryParams(u, nil)
	if u.RawQuery != original {
		t.Fatalf("expected RawQuery unchanged, got %q", u.RawQuery)
	}
}

func TestApplyQueryParams_RepeatedKeys_AppendsUniqueValues(t *testing.T) {
	u, _ := url.Parse("https://api.example.com/items?tag=go")
	q := url.Values{"tag": {"infra", "ports"}}
	applyQueryParams(u, q)

	got := u.Query()["tag"]
	if len(got) != 3 {
		t.Fatalf("expected 3 tag values, got %d (%v)", len(got), got)
	}
	if got[0] != "go" || got[1] != "infra" || got[2] != "ports" {
		t.Fatalf("expected [go infra ports], got %v", got)
	}
}

func TestApplyQueryParams_DuplicatePairs_AreDeduped(t *testing.T) {
	u, _ := url.Parse("https://api.example.com/items?page=1")
	q := url.Values{"page": {"1"}}
	applyQueryParams(u, q)

	got := u.Query()["page"]
	if len(got) != 1 {
		t.Fatalf("expected 1 page value after dedup, got %d (%v)", len(got), got)
	}
}

func TestApplyHeaderParams_OverridesByKey(t *testing.T) {
	dst := http.Header{"X-Config": {"cfg"}, "X-Keep": {"keep"}}
	src := http.Header{"X-Config": {"override"}, "X-Request": {"req"}}

	applyHeaderParams(dst, src)

	if dst.Get("X-Config") != "override" {
		t.Fatalf("expected X-Config=override, got %q", dst.Get("X-Config"))
	}
	if dst.Get("X-Keep") != "keep" {
		t.Fatalf("expected X-Keep=keep, got %q", dst.Get("X-Keep"))
	}
	if dst.Get("X-Request") != "req" {
		t.Fatalf("expected X-Request=req, got %q", dst.Get("X-Request"))
	}
}

func TestApplyHeaderParams_ClonesValues(t *testing.T) {
	dst := http.Header{}
	src := http.Header{"X-Config": {"cfg"}}

	applyHeaderParams(dst, src)
	src["X-Config"][0] = "mutated"

	if dst.Get("X-Config") != "cfg" {
		t.Fatalf("expected cloned value cfg, got %q", dst.Get("X-Config"))
	}
}
