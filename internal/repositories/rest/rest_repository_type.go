package rest

import (
	"github.com/brunojet/go-infra-ports/pkg/http_clients/contracts"
)

type HttpClient = contracts.HttpClient

type restRepository struct {
	client   HttpClient
	registry RestRegistry
	opts     *repositoryOptions
}
