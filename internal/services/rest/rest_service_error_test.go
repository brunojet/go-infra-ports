package rest

import (
	"testing"
)

func TestErrRestServiceRepositoryNil_IsError(t *testing.T) {
	if errRestServiceRepositoryNil == nil {
		t.Fatalf("expected errRestServiceRepositoryNil to be non-nil")
	}
}

func TestErrRestServiceRepositoryNil_HasMessage(t *testing.T) {
	msg := errRestServiceRepositoryNil.Error()
	if msg == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestErrRestServiceResponseNil_IsError(t *testing.T) {
	if errRestServiceResponseNil == nil {
		t.Fatalf("expected errRestServiceResponseNil to be non-nil")
	}
}

func TestErrRestServiceResponseNil_HasMessage(t *testing.T) {
	msg := errRestServiceResponseNil.Error()
	if msg == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestErrRestServiceUpstreamMappingFailed_WrapsOriginalError(t *testing.T) {
	originalErr := errRestDownstreamMapperResponseTypeAssertion
	err := errRestServiceUpstreamMappingFailed("test operation", originalErr)

	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if err.Error() == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestErrRestServiceDownstreamMappingFailed_WrapsOriginalError(t *testing.T) {
	originalErr := errRestDownstreamMapperResponseTypeAssertion
	err := errRestServiceDownstreamMappingFailed("test operation", originalErr)

	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if err.Error() == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestErrRestServiceInvalidNon2xxStatus_IsError(t *testing.T) {
	if errRestServiceInvalidNon2xxStatus == nil {
		t.Fatalf("expected errRestServiceInvalidNon2xxStatus to be non-nil")
	}
}

func TestErrRestServiceInvalidNon2xxStatus_HasMessage(t *testing.T) {
	msg := errRestServiceInvalidNon2xxStatus.Error()
	if msg == "" {
		t.Fatalf("expected error message, got empty string")
	}
}

func TestErrRestServiceNilResponseData_IsError(t *testing.T) {
	if errRestServiceNilResponseData == nil {
		t.Fatalf("expected errRestServiceNilResponseData to be non-nil")
	}
}

func TestErrRestServiceNilResponseData_HasMessage(t *testing.T) {
	msg := errRestServiceNilResponseData.Error()
	if msg == "" {
		t.Fatalf("expected error message, got empty string")
	}
}
