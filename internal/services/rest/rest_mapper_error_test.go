package rest

import (
	"testing"
)

func TestErrRestDownstreamMapperResponseTypeAssertion_IsError(t *testing.T) {
	if errRestDownstreamMapperResponseTypeAssertion == nil {
		t.Fatalf("expected errRestDownstreamMapperResponseTypeAssertion to be non-nil")
	}
}

func TestErrRestDownstreamMapperResponseTypeAssertion_HasMessage(t *testing.T) {
	msg := errRestDownstreamMapperResponseTypeAssertion.Error()
	if msg == "" {
		t.Fatalf("expected error message, got empty string")
	}
	if msg != "rest downstream mapper response type assertion failed" {
		t.Fatalf("expected specific error message, got %q", msg)
	}
}
