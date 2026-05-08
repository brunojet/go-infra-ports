package rest

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorWrappers_WrapOriginalError(t *testing.T) {
	base := errors.New("base")
	cases := []struct {
		name string
		fn   func(error) error
		want string
	}{
		{"resolve request marshal", errRestResolveRequestMarshal, "resolve request marshal"},
		{"resolve envelope request marshal", errRestResolveEnvelopeRequestMarshal, "resolve envelope request marshal"},
		{"resolve response unmarshal", errRestResolveResponseUnmarshal, "resolve response unmarshal"},
		{"resolve responses unmarshal raw list", errRestResolveResponsesUnmarshalRawList, "resolve responses unmarshal raw list"},
		{"resolve envelope response unmarshal", errRestResolveEnvelopeResponseUnmarshal, "resolve envelope response unmarshal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(base)
			if !errors.Is(err, base) {
				t.Fatalf("expected wrapped error to contain base: %v", err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q in error %q", tc.want, err.Error())
			}
		})
	}
}

func TestResolveResponsesErrorBuilders(t *testing.T) {
	errLen := errRestResolveResponsesNewSliceLen(1, 3)
	if !strings.Contains(errLen.Error(), "length 1, expected 3") {
		t.Fatalf("unexpected slice len error: %v", errLen)
	}

	base := errors.New("bad-item")
	errItem := errRestResolveResponsesUnmarshalItem(2, base)
	if !errors.Is(errItem, base) {
		t.Fatalf("expected wrapped item error to contain base: %v", errItem)
	}
	if !strings.Contains(errItem.Error(), "item 2") {
		t.Fatalf("unexpected item error: %v", errItem)
	}
}

func TestRegisterErrorBuilders(t *testing.T) {
	// Pass a zero-value RestMethod to avoid depending on external constants
	errMethod := errRestRegisterRequestInvalidMethod(RestMethod(0))
	if !strings.Contains(errMethod.Error(), "invalid method") {
		t.Fatalf("unexpected method error: %v", errMethod)
	}

	errRange := errRestRegisterResponseOutOfRange(99, 200, 299)
	if !strings.Contains(errRange.Error(), "99") || !strings.Contains(errRange.Error(), "[200, 299]") {
		t.Fatalf("unexpected range error: %v", errRange)
	}
}
