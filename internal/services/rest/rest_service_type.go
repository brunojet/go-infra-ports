package rest

import (
	repcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type (
	RestRepository          = repcts.RestRepository
	RestRequest             = repcts.RestRequest
	RestResponse            = repcts.RestResponse
	RestResponses           = repcts.RestResponses
	RestResponseSpec        = repcts.RestResponseSpec
	Service[C, R, U any]    = svccts.Service[C, R, U]
	ServiceCreate[C any]    = svccts.ServiceCreate[C]
	ServiceUpdate[U any]    = svccts.ServiceUpdate[U]
	ServiceSave[C any]      = svccts.ServiceSave[C]
	ServiceResponse[R any]  = svccts.ServiceResponse[R]
	ServiceResponses[R any] = svccts.ServiceResponses[R]
	ServiceMeta             = svccts.ServiceMeta
	RequestContext          = types.RequestContext
	ResponseContext         = types.ResponseContext
)
