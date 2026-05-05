package http_util

import (
	"errors"
	"reflect"
	"testing"
)

func TestSanitizeAndValidatePath_ErrorsAndNormalization(t *testing.T) {
	_, err := SanitizeAndValidatePath("")
	if !errors.Is(err, errRepositoryBaseURLEmpty) {
		t.Fatalf("expected errRepositoryBaseURLEmpty, got %v", err)
	}

	_, err = SanitizeAndValidatePath("/users?bad=1")
	if !errors.Is(err, errRepositoryPathInvalidChars) {
		t.Fatalf("expected errRepositoryPathInvalidChars, got %v", err)
	}

	_, err = SanitizeAndValidatePath("/users/{id")
	if !errors.Is(err, errRepositoryPathInvalidStructure) {
		t.Fatalf("expected errRepositoryPathInvalidStructure, got %v", err)
	}

	got, err := SanitizeAndValidatePath(" /api/v1/users/ ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/api/v1/users" {
		t.Fatalf("unexpected sanitized path: %q", got)
	}
	// root path should be accepted and normalized to "/"
	gotRoot, err := SanitizeAndValidatePath("/")
	if err != nil {
		t.Fatalf("unexpected error for root: %v", err)
	}
	if gotRoot != "/" {
		t.Fatalf("expected '/', got %q", gotRoot)
	}
}

func TestExtractPathParams_Behavior(t *testing.T) {
	// no params
	fmtNo, names, err := ExtractPathParams("/api/v1/static")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmtNo != "/api/v1/static" {
		t.Fatalf("expected same format, got %q", fmtNo)
	}
	if len(names) != 0 {
		t.Fatalf("expected no names, got %v", names)
	}

	// single param
	fmtOne, names, err := ExtractPathParams("/users/{id}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmtOne != "/users/%s" {
		t.Fatalf("unexpected format: %q", fmtOne)
	}
	if !reflect.DeepEqual(names, []string{"id"}) {
		t.Fatalf("unexpected names: %v", names)
	}

	// multiple params
	fmtMany, names, err := ExtractPathParams("/users/{id}/orders/{orderId}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmtMany != "/users/%s/orders/%s" {
		t.Fatalf("unexpected format: %q", fmtMany)
	}
	if !reflect.DeepEqual(names, []string{"id", "orderId"}) {
		t.Fatalf("unexpected names: %v", names)
	}

	// invalid param style (colon) should return an invalid-template wrapper
	_, _, err = ExtractPathParams("/users/:id")
	if err == nil || !errors.Is(err, errInvalidPathTemplate) {
		t.Fatalf("expected errInvalidPathTemplate, got %v", err)
	}
}
