package net_http

import (
	"github.com/brunojet/go-infra-ports/internal/dtos"
)

// DefaultProblemDetails é um alias para `dtos.ProblemDetails`.
// Mantém compatibilidade com código existente que usa o tipo antigo.
type DefaultInformation = dtos.Information
type DefaultRedirection = dtos.Redirection
type DefaultProblemDetails = dtos.ProblemDetails
