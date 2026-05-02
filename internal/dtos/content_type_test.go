package dtos

import "testing"

func TestContentTypeConstants(t *testing.T) {
	if ContentTypeProblemJSON != "application/problem+json" {
		t.Fatalf("unexpected ContentTypeProblemJSON value: %s", ContentTypeProblemJSON)
	}
	if ContentTypeJSON != "application/json" {
		t.Fatalf("unexpected ContentTypeJSON value: %s", ContentTypeJSON)
	}
	if ContentTypeURIList != "text/uri-list" {
		t.Fatalf("unexpected ContentTypeURIList value: %s", ContentTypeURIList)
	}
}
