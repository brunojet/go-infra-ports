package rest

import (
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/http_clients/mocks"
)

func TestRestRepositoryStruct_ImplementsRestRepository(t *testing.T) {
	var _ RestRepository = (*restRepository)(nil)
}

func TestHttpClientAlias_SatisfiedByConcreteType(t *testing.T) {
	var _ HttpClient = (*mocks.MockHttpClient)(nil)
}
