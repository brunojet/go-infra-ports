package handlers

import (
	"net/http"
	"regexp"
	"testing"

	"go.uber.org/mock/gomock"

	internalnethttp "github.com/brunojet/go-infra-ports/internal/handlers/net_http"
	mocksvc "github.com/brunojet/go-infra-ports/pkg/services/mocks"
)

type apiCreate struct{ Name string }
type apiResp struct {
	ID   string
	Name string
}
type apiUpdate struct{ Name string }

func TestNewMiddlewareRegistry_ReturnsNonNil(t *testing.T) {
	if got := NewMiddlewareRegistry(); got == nil {
		t.Fatalf("expected non-nil middleware registry")
	}
}

func TestWithCollection_ReturnsOption(t *testing.T) {
	if opt := WithCollection(MethodList, "users"); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}

func TestWithMiddlewareRegistry_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for nil registry")
		}
	}()
	_ = WithMiddlewareRegistry(nil)
}

func TestNewNetHttpHandler_ReturnsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocksvc.NewMockService[apiCreate, apiResp, apiUpdate](ctrl)
	mux := http.NewServeMux()

	h := NewNetHttpHandler(mux, mock,
		WithCollection(MethodCreate, "users"),
	)

	if h == nil {
		t.Fatalf("expected non-nil handler")
	}
}

func TestWithInstance_ReturnsOption(t *testing.T) {
	rx := regexp.MustCompile(`^\d+$`)
	if opt := WithInstance(MethodGet, "users/{id}", rx); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}

func TestWithInformation_ReturnsOption(t *testing.T) {
	if opt := WithInformation(func(status int, meta internalnethttp.ServiceMeta) any {
		return map[string]any{"status": status}
	}); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}

func TestWithEnvelope_ReturnsOption(t *testing.T) {
	if opt := WithEnvelope(func(status int, meta internalnethttp.ServiceMeta, resp any) any {
		return map[string]any{"data": resp}
	}); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}

func TestWithRedirection_ReturnsOption(t *testing.T) {
	if opt := WithRedirection(func(status int, meta internalnethttp.ServiceMeta) (body any, location string) {
		return map[string]any{"to": meta.Location}, meta.Location
	}); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}

func TestWithProblemDetails_ReturnsOption(t *testing.T) {
	if opt := WithProblemDetails(func(status int, instance string, meta internalnethttp.ServiceMeta) any {
		return map[string]any{"status": status, "instance": instance}
	}); opt == nil {
		t.Fatalf("expected non-nil handler option")
	}
}
