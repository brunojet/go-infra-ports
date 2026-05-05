package net_http

import (
	"regexp"

	http_helper "github.com/brunojet/go-infra-ports/internal/helpers/http_helper"
)

// sanitizePath normaliza um path de rota: remove barras duplas, leading e trailing slash.
// Usa o helper central `http_util.SanitizeAndValidatePath` para validação.
// OBS: a validação completa é delegada ao helper — entradas vazias ou inválidas
// provocarão panic (o helper retorna erro para entradas vazias).
func sanitizePath(p string) string {
	s, err := http_helper.SanitizeAndValidatePath(p)
	if err != nil {
		panic("net_http: invalid path: " + err.Error())
	}
	return s
}

// extractParams extracts the {params} names from pathFmt and associates
// positional regexes. Uses `http_util.ExtractPathParams` for validation and
// extraction; `formats` are assigned in positional order.
func extractParams(pathFmt string, formats []*regexp.Regexp) []paramFormat {
	// Delegate template validation/extraction to `http_util.ExtractPathParams`.
	// It returns a sanitized format (with %s placeholders) and the parameter
	// names in positional order. We then map those names to the package-local
	// `paramFormat` type and apply positional regexes from `formats`. This
	// keeps `http_util` decoupled from `net_http` types while centralizing
	// template validation/normalization.
	names, err := http_helper.ExtractPathParams(pathFmt)
	if err != nil {
		panic("net_http: invalid path template: " + err.Error())
	}
	lenFmts := len(formats)
	lenNames := len(names)
	if lenFmts > 0 && lenFmts != lenNames {
		panic("net_http: number of formats does not match number of parameters in path template")
	}
	params := make([]paramFormat, 0, lenNames)
	for i, name := range names {
		pf := paramFormat{name: name}
		if i < lenFmts {
			pf.format = formats[i]
		}
		params = append(params, pf)
	}
	return params
}
