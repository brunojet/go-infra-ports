// Package repositories exposes public repository contracts and aliases.
package repositories

import (
	contracts "github.com/brunojet/go-infra-ports/pkg/repositories/contracts"
	restcontracts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
)

// RepositoryCreate aliases the public create input contract.
type RepositoryCreate[C any] = contracts.RepositoryCreate[C]

// RepositoryUpdate aliases the public update input contract.
type RepositoryUpdate[U any] = contracts.RepositoryUpdate[U]

// RepositorySave aliases the public save input contract.
type RepositorySave[C any] = contracts.RepositorySave[C]

// RepositoryResponse aliases the single-entity repository output.
type RepositoryResponse[R any] = contracts.RepositoryResponse[R]

// RepositoryResponses aliases the multi-entity repository output.
type RepositoryResponses[R any] = contracts.RepositoryResponses[R]

// Repository aliases the public storage-agnostic repository contract.
type Repository[C, R, U any] = contracts.Repository[C, R, U]

// RestCreate aliases REST create input.
type RestCreate[C any] = restcontracts.RestCreate[C]

// RestUpdate aliases REST update input.
type RestUpdate[U any] = restcontracts.RestUpdate[U]

// RestSave aliases REST save input.
type RestSave[C any] = restcontracts.RestSave[C]

// RestResponse aliases REST single-entity output.
type RestResponse[R any] = restcontracts.RestResponse[R]

// RestResponses aliases REST multi-entity output.
type RestResponses[R any] = restcontracts.RestResponses[R]

// RestRepository aliases the public REST repository contract.
type RestRepository[C, R, U any] = restcontracts.RestRepository[C, R, U]
