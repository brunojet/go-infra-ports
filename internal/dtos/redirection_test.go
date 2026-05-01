package dtos

import "testing"

func TestRedirection_ContentType(t *testing.T) {
	var r Redirection
	if got := r.ContentType(); got != ContentTypeJSON {
		t.Fatalf("expected %q, got %q", ContentTypeJSON, got)
	}
}
