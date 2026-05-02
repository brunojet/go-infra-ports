package dtos

import "testing"

func TestInformation_ContentType(t *testing.T) {
	var info Information
	if got := info.ContentType(); got != ContentTypeJSON {
		t.Fatalf("expected %q, got %q", ContentTypeJSON, got)
	}
}
