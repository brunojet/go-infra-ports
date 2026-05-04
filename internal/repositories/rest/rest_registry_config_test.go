package rest

import (
	"net/http"
	"testing"
)

func TestWithResponseRegistersInResponsesMap(t *testing.T) {
	cfg := newRegistryConfig(WithResponse(&DefaultRestResponse{}, http.StatusCreated))

	if cfg.Responses[http.StatusCreated] == nil {
		t.Fatalf("expected response spec registered for status %d", http.StatusCreated)
	}
	if cfg.ResponseEnvelopes[http.StatusCreated] != nil {
		t.Fatalf("did not expect response envelope spec registered for status %d", http.StatusCreated)
	}
}

func TestWithInformationDefaultStatusDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got %v", r)
		}
	}()

	cfg := newRegistryConfig(WithInformation(&DefaultRestResponse{}))
	if cfg.Informations[DefaultStatusCode] == nil {
		t.Fatalf("expected default information spec to be registered")
	}
}

func TestWithRedirectionAndProblemRegisterInProperMaps(t *testing.T) {
	cfg := newRegistryConfig(
		WithRedirection(&DefaultRestResponse{}, http.StatusMovedPermanently),
		WithProblem(&DefaultRestResponse{}, http.StatusBadRequest),
	)

	if cfg.Redirections[http.StatusMovedPermanently] == nil {
		t.Fatalf("expected redirection registered for status %d", http.StatusMovedPermanently)
	}
	if cfg.Problems[http.StatusBadRequest] == nil {
		t.Fatalf("expected problem registered for status %d", http.StatusBadRequest)
	}
}

func TestRegisterRequest_DefaultMethodWhenEmpty(t *testing.T) {
	ro := newRegistryOptions()
	ro.registerRequest(&DefaultRestRequest{}, ro.Requests, MethodCreate)
	if ro.Requests[MethodCreate] == nil {
		t.Fatal("expected method registration for MethodCreate")
	}
}

func TestRegisterRequest_PanicsOnInvalidMethod(t *testing.T) {
	ro := newRegistryOptions()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invalid method")
		}
	}()
	ro.registerRequest(&DefaultRestRequest{}, ro.Requests, MethodList)
}

func TestRegisterRequest_PanicsOnNilSpec(t *testing.T) {
	ro := newRegistryOptions()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil request spec")
		}
	}()
	ro.registerRequest(nil, ro.Requests, MethodCreate)
}

func TestRegisterResponse_DefaultStatusWhenEmpty(t *testing.T) {
	ro := newRegistryOptions()
	ro.registerResponse(&DefaultRestResponse{}, ro.Responses, http.StatusOK, http.StatusOK+99)
	if ro.Responses[DefaultStatusCode] == nil {
		t.Fatal("expected default status registration")
	}
}

func TestRegisterResponse_PanicsOnOutOfRange(t *testing.T) {
	ro := newRegistryOptions()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for out-of-range status")
		}
	}()
	ro.registerResponse(&DefaultRestResponse{}, ro.Responses, http.StatusOK, http.StatusOK+99, http.StatusContinue)
}

func TestRegisterResponse_DefaultCodeExplicitAndNilSpecPanic(t *testing.T) {
	ro := newRegistryOptions()
	ro.registerResponse(&DefaultRestResponse{}, ro.Responses, http.StatusOK, http.StatusOK+99, DefaultStatusCode)
	if ro.Responses[DefaultStatusCode] == nil {
		t.Fatal("expected default status to be explicitly registered")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil response spec")
		}
	}()
	ro.registerResponse(nil, ro.Responses, http.StatusOK, http.StatusOK+99, http.StatusOK)
}

func TestRegisterResponseEnvelope_DefaultAndStatusRegistration(t *testing.T) {
	ro := newRegistryOptions()
	ro.registerResponseEnvelope(&responseEnvelope{})
	ro.registerResponseEnvelope(&responseEnvelope{}, http.StatusOK)

	if ro.ResponseEnvelopes[DefaultStatusCode] == nil {
		t.Fatal("expected default envelope status registration")
	}
	if ro.ResponseEnvelopes[http.StatusOK] == nil {
		t.Fatal("expected envelope status registration for 200")
	}
}

func TestRegisterResponseEnvelope_PanicsOnOutOfRange(t *testing.T) {
	ro := newRegistryOptions()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for out-of-range envelope status")
		}
	}()
	ro.registerResponseEnvelope(&responseEnvelope{}, http.StatusMultipleChoices)
}

func TestRegisterResponseEnvelope_PanicsOnNilSpec(t *testing.T) {
	ro := newRegistryOptions()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil envelope spec")
		}
	}()
	ro.registerResponseEnvelope(nil)
}
