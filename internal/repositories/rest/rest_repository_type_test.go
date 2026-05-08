package rest

import (
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/http_clients/mocks"
)

func TestRestRepositoryStruct_ImplementsRestRepositoryV2(t *testing.T) {
	var _ RestRepository = (*restRepository)(nil)
}

func TestHttpClientAlias_SatisfiedByConcreteType(t *testing.T) {
	var _ HttpClient = (*mocks.MockHttpClient)(nil)
}
