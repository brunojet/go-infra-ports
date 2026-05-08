package rest

import (
	"net/http"
	"testing"
)

type sampleT struct{ V int }

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

func TestRegistryOptionsDefaults(t *testing.T) {
	ro := newRegistryOptions()
	if ro == nil {
		t.Fatal("newRegistryOptions returned nil")
	}

	// Requests defaults
	rc := ro.requests[MethodCreate]
	ru := ro.requests[MethodUpdate]
	rs := ro.requests[MethodSave]
	if rc == nil || ru == nil || rs == nil {
		t.Fatalf("expected default request specs to be non-nil: %v %v %v", rc, ru, rs)
	}
	if rc != ru || rc != rs {
		t.Fatalf("expected default request specs to be the same instance")
	}

	// Responses/Informations/Redirections/Problems share the same default spec at DefaultStatusCode
	resp := ro.responses[defaultStatusCode]
	info := ro.informations[defaultStatusCode]
	red := ro.redirections[defaultStatusCode]
	prob := ro.problems[defaultStatusCode]
	if resp == nil || info == nil || red == nil || prob == nil {
		t.Fatalf("expected default response specs to be non-nil")
	}
	if resp != info || resp != red || resp != prob {
		t.Fatalf("expected default response specs to be the same instance")
	}
}

func TestRegisterRequest_SuccessAndPanics(t *testing.T) {
	ro := newRegistryOptions()
	spec := NewDataSpecOf[sampleT]()

	// success: register for Create and Update
	ro.registerRequest(spec, ro.requests, MethodCreate|MethodUpdate)
	if ro.requests[MethodCreate] != spec || ro.requests[MethodUpdate] != spec {
		t.Fatalf("expected request mapping populated")
	}

	// panics: nil spec
	mustPanic(t, func() { ro.registerRequest(nil, ro.requests, MethodCreate) })

	// panics: zero methods
	mustPanic(t, func() { ro.registerRequest(spec, ro.requests, 0) })

	// panics: methods that don't include a write bit
	mustPanic(t, func() { ro.registerRequest(spec, ro.requests, MethodList) })
}

func TestRegisterResponse_SuccessAndPanics(t *testing.T) {
	ro := newRegistryOptions()
	spec := NewDataSpecOf[sampleT]()

	// default behavior (no status codes) -> DefaultStatusCode used
	target := make(map[int]RestDataSpec)
	ro.registerResponse(spec, target, http.StatusOK, http.StatusOK+99)
	if target[defaultStatusCode] != spec {
		t.Fatalf("expected DefaultStatusCode mapping set")
	}

	// explicit valid code
	target2 := make(map[int]RestDataSpec)
	ro.registerResponse(spec, target2, http.StatusOK, http.StatusOK+99, http.StatusOK)
	if target2[http.StatusOK] != spec {
		t.Fatalf("expected explicit status mapping set")
	}

	// panics: nil spec
	mustPanic(t, func() { ro.registerResponse(nil, target, http.StatusOK, http.StatusOK+99) })

	// panics: out of range status code
	mustPanic(t, func() { ro.registerResponse(spec, target, http.StatusOK, http.StatusOK+99, http.StatusContinue) })
}

func TestRegisterResponseEnvelope_SuccessAndPanics(t *testing.T) {
	ro := newRegistryOptions()
	env := NewEnvelopeSpec("data", "meta")

	// default -> DefaultStatusCode
	ro.registerResponseEnvelope(env)
	if ro.responseEnvelopes[defaultStatusCode] == nil {
		t.Fatalf("expected default envelope mapping set")
	}

	// explicit valid code
	ro2 := newRegistryOptions()
	ro2.registerResponseEnvelope(env, http.StatusOK)
	if ro2.responseEnvelopes[http.StatusOK] == nil {
		t.Fatalf("expected explicit envelope mapping set")
	}

	// panics: nil spec
	mustPanic(t, func() { ro.registerResponseEnvelope(nil) })

	// panics: out of range code (not in 200..299 and not DefaultStatusCode)
	mustPanic(t, func() { ro.registerResponseEnvelope(env, http.StatusInternalServerError) })
}

func TestWithOptions_RegisterAndPanics(t *testing.T) {
	// WithRequestOf
	cfg := newRegistryConfig(WithRequestOf[sampleT](MethodCreate | MethodSave))
	if _, ok := cfg.requests[MethodCreate].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected dataSpecOf for MethodCreate")
	}
	if _, ok := cfg.requests[MethodSave].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected dataSpecOf for MethodSave")
	}

	// WithRequestEnvelope valid
	cfg2 := newRegistryConfig(WithRequestEnvelope("data", MethodCreate))
	if _, ok := cfg2.requestsEnvelopes[MethodCreate].(*envelopeSpecOf); !ok {
		t.Fatalf("expected envelopeSpecOf for MethodCreate")
	}

	// WithRequestEnvelope invalid methods -> panic on apply
	mustPanic(t, func() { newRegistryConfig(WithRequestEnvelope("data", MethodList)) })

	// WithResponseOf default and explicit
	cfg3 := newRegistryConfig(WithResponseOf[sampleT]())
	if _, ok := cfg3.responses[defaultStatusCode].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected dataSpecOf at DefaultStatusCode")
	}
	cfg4 := newRegistryConfig(WithResponseOf[sampleT](http.StatusOK))
	if _, ok := cfg4.responses[http.StatusOK].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected dataSpecOf at 200")
	}
	mustPanic(t, func() { newRegistryConfig(WithResponseOf[sampleT](http.StatusContinue)) })

	// WithResponseEnvelope default and explicit
	cfg5 := newRegistryConfig(WithResponseEnvelope("data", "meta"))
	if _, ok := cfg5.responseEnvelopes[defaultStatusCode].(*envelopeSpecOf); !ok {
		t.Fatalf("expected envelopeSpecOf at DefaultStatusCode")
	}
	cfg6 := newRegistryConfig(WithResponseEnvelope("d", "m", http.StatusOK))
	if _, ok := cfg6.responseEnvelopes[http.StatusOK].(*envelopeSpecOf); !ok {
		t.Fatalf("expected envelopeSpecOf at 200")
	}
	mustPanic(t, func() { newRegistryConfig(WithResponseEnvelope("d", "m", http.StatusInternalServerError)) })

	// WithInformationOf (100..199)
	cfg7 := newRegistryConfig(WithInformationOf[sampleT]())
	if _, ok := cfg7.informations[defaultStatusCode].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected DefaultStatusCode information spec")
	}
	cfg8 := newRegistryConfig(WithInformationOf[sampleT](http.StatusContinue))
	if _, ok := cfg8.informations[http.StatusContinue].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected information spec at 100")
	}
	mustPanic(t, func() { newRegistryConfig(WithInformationOf[sampleT](http.StatusOK)) })

	// WithRedirectionOf (300..399)
	cfg9 := newRegistryConfig(WithRedirectionOf[sampleT]())
	if _, ok := cfg9.redirections[defaultStatusCode].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected DefaultStatusCode redirection spec")
	}
	cfg10 := newRegistryConfig(WithRedirectionOf[sampleT](http.StatusMultipleChoices))
	if _, ok := cfg10.redirections[http.StatusMultipleChoices].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected redirection spec at 300")
	}
	mustPanic(t, func() { newRegistryConfig(WithRedirectionOf[sampleT](http.StatusOK)) })

	// WithProblemOf (400..598)
	cfg11 := newRegistryConfig(WithProblemOf[sampleT]())
	if _, ok := cfg11.problems[defaultStatusCode].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected DefaultStatusCode problem spec")
	}
	cfg12 := newRegistryConfig(WithProblemOf[sampleT](http.StatusInternalServerError))
	if _, ok := cfg12.problems[http.StatusInternalServerError].(*dataSpecOf[sampleT]); !ok {
		t.Fatalf("expected problem spec at 500")
	}
	mustPanic(t, func() { newRegistryConfig(WithProblemOf[sampleT](http.StatusOK)) })
}
