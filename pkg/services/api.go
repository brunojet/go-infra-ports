// Package services exposes public service contracts and aliases.
package services

import contracts "github.com/brunojet/go-infra-ports/pkg/services/contracts"

// ServiceCreate aliases the public create input contract.
type ServiceCreate[C any] = contracts.ServiceCreate[C]

// ServiceUpdate aliases the public update input contract.
type ServiceUpdate[U any] = contracts.ServiceUpdate[U]

// ServiceSave aliases the public save input contract.
type ServiceSave[C any] = contracts.ServiceSave[C]

// ServiceMeta aliases public service metadata.
type ServiceMeta = contracts.ServiceMeta

// ServiceResponse aliases the single-entity service output.
type ServiceResponse[R any] = contracts.ServiceResponse[R]

// ServiceResponses aliases the multi-entity service output.
type ServiceResponses[R any] = contracts.ServiceResponses[R]

// Service aliases the public service port contract.
type Service[C, R, U any] = contracts.Service[C, R, U]
