package net_http

import (
	"path"
	"regexp"
	"strings"
)

// sanitizePath normaliza um path de rota: remove barras duplas, leading e trailing slash.
// Usa path.Clean (semântica URL) para garantir compatibilidade com http.ServeMux.
func sanitizePath(p string) string {
	return strings.Trim(path.Clean("/"+p), "/")
}

// extractParams extrai os nomes dos {params} do pathFmt e associa os regexes
// posicionalmente. O número de params é definido pelo próprio pathFmt.
func extractParams(pathFmt string, formats []*regexp.Regexp) []paramFormat {
	var params []paramFormat
	idx := 0
	for _, seg := range strings.Split(pathFmt, "/") {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			pf := paramFormat{name: seg[1 : len(seg)-1]}
			if idx < len(formats) {
				pf.format = formats[idx]
			}
			params = append(params, pf)
			idx++
		}
	}
	return params
}
