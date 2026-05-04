package net_http

import (
	"testing"

	"github.com/brunojet/go-infra-ports/internal/dtos"
)

func TestDefaultAliasTypes_Behavior(t *testing.T) {
	info := DefaultInformation{Message: "ok", Code: 100}
	redir := DefaultRedirection{Location: "/items/1"}
	problem := DefaultProblemDetails{Title: "bad request", Status: 400}

	if info.ContentType() != dtos.ContentTypeJSON {
		t.Fatalf("expected JSON content type, got %s", info.ContentType())
	}
	if redir.ContentType() != dtos.ContentTypeJSON {
		t.Fatalf("expected JSON content type, got %s", redir.ContentType())
	}
	if problem.ContentType() != dtos.ContentTypeProblemJSON {
		t.Fatalf("expected problem+json content type, got %s", problem.ContentType())
	}
}
