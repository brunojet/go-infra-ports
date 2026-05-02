package rest

import (
	"errors"
	"fmt"
)

var (
	errRestResolveRequestBodyNil             = errors.New("rest: resolve request: requestBody must not be nil")
	errRestResolveRequestSpecBodyNil         = errors.New("rest: resolve request: body must not be nil")
	errRestResolveEnvelopeRequestBodyNil     = errors.New("rest: resolve envelope request: dataBody must not be nil")
	errRestResolveEnvelopeRequestSpecNewNil  = errors.New("rest: resolve envelope request: envelope spec New returned nil")
	errRestResolveResponseBodyNil            = errors.New("rest: resolve response: body must not be nil")
	errRestResolveResponseSpecNil            = errors.New("rest: resolve response: resolved spec must not be nil")
	errRestResolveResponseSpecNewNil         = errors.New("rest: resolve response: response spec New returned nil")
	errRestResolveResponsesBodiesNil         = errors.New("rest: resolve responses: bodies must not be nil")
	errRestResolveResponsesSpecNewNil        = errors.New("rest: resolve responses: response spec New returned nil")
	errRestResolveEnvelopeResponseBodyNil    = errors.New("rest: resolve envelope response: dataBody must not be nil")
	errRestResolveEnvelopeResponseMetaNil    = errors.New("rest: resolve envelope response: meta must not be nil")
	errRestResolveEnvelopeResponseSpecNewNil = errors.New("rest: resolve envelope response: spec New returned nil")

	errRestRegisterRequestSpecNil  = errors.New("rest: registerRequest: spec must not be nil")
	errRestRegisterResponseSpecNil = errors.New("rest: registerResponse: spec must not be nil")
	errRestNewRequestSpecNil       = errors.New("rest: NewRequestSpec: spec New returned nil")
)

func errRestResolveRequestMarshal(err error) error {
	return fmt.Errorf("rest: resolve request marshal: %w", err)
}

func errRestResolveEnvelopeRequestMarshal(err error) error {
	return fmt.Errorf("rest: resolve envelope request marshal: %w", err)
}

func errRestResolveResponseUnmarshal(err error) error {
	return fmt.Errorf("rest: resolve response unmarshal: %w", err)
}

func errRestResolveResponsesUnmarshalRawList(err error) error {
	return fmt.Errorf("rest: resolve responses unmarshal raw list: %w", err)
}

func errRestResolveResponsesNewSliceLen(got, expected int) error {
	return fmt.Errorf("rest: resolve responses: spec NewSlice returned length %d, expected %d", got, expected)
}

func errRestResolveResponsesUnmarshalItem(index int, err error) error {
	return fmt.Errorf("rest: resolve responses unmarshal item %d: %w", index, err)
}

func errRestResolveEnvelopeResponseUnmarshal(err error) error {
	return fmt.Errorf("rest: resolve envelope response unmarshal: %w", err)
}

func errRestRegisterRequestInvalidMethod(method string) error {
	return fmt.Errorf("rest: registerRequest: invalid method %q, must be POST, PUT or PATCH", method)
}

func errRestRegisterResponseOutOfRange(code, firstStatusCode, lastStatusCode int) error {
	return fmt.Errorf("rest: registerResponse: status code %d out of range [%d, %d]", code, firstStatusCode, lastStatusCode)
}
