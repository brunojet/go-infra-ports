package net_http

import (
	hndcts "github.com/brunojet/go-infra-ports/pkg/handlers/net_http/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
)

type (
	NetHttpHandler          = hndcts.NetHttpHandler
	Service[C, R, U any]    = svccts.Service[C, R, U]
	ServiceCreate[C any]    = svccts.ServiceCreate[C]
	ServiceUpdate[U any]    = svccts.ServiceUpdate[U]
	ServiceSave[C any]      = svccts.ServiceSave[C]
	ServiceResponse[R any]  = svccts.ServiceResponse[R]
	ServiceResponses[R any] = svccts.ServiceResponses[R]
	ServiceMeta             = svccts.ServiceMeta
)
