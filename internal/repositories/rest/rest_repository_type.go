package rest

import (
	"github.com/brunojet/go-infra-ports/pkg/http_clients/contracts"
	rpcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
)

type HttpClient = contracts.HttpClient

type RestRepository = rpcts.RestRepository

type RestRequest = rpcts.RestRequest

type RestResponse = rpcts.RestResponse

type RestResponses = rpcts.RestResponses

type restRepository struct {
	client   HttpClient
	registry RestRegistry
	opts     *repositoryOptions
}
