package net_http

import (
	"github.com/brunojet/go-infra-ports/internal/dtos"
)

// DefaultInformation é um alias para `dtos.Information`.
// Mantém compatibilidade com código existente que usa o tipo antigo.
type DefaultInformation = dtos.Information

// DefaultRedirection é um alias para `dtos.Redirection`.
// Mantém compatibilidade com código existente que usa o tipo antigo.
type DefaultRedirection = dtos.Redirection

// DefaultProblemDetails é um alias para `dtos.ProblemDetails`.
// Mantém compatibilidade com código existente que usa o tipo antigo.
type DefaultProblemDetails = dtos.ProblemDetails
