package http_util

import (
	"path"
	"regexp"
	"strings"
)

var (
	// allowed chars: letters, digits, underscore, hyphen, braces and slashes
	// Leading slash is optional: the helper will normalize and insert it.
	allowedCharsRegex = regexp.MustCompile(`^/?[A-Za-z0-9_/\-\{\}]*$`)
	// validPathRegex: segments are either plain tokens or {param} (no :param support)
	// Accept root '/' explicitly to match net/http patterns.
	validPathRegex = regexp.MustCompile(`^/$|^(/([A-Za-z0-9_-]+|\{[A-Za-z0-9_]+\}))+$`)
	// templateRegex only matches {name}
	templateRegex = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)
)

// SanitizeAndValidatePath trims whitespace, checks for emptiness, validates allowed characters and structure, and normalizes the path.
func SanitizeAndValidatePath(basePath string) (string, error) {
	p := strings.TrimSpace(basePath)
	switch {
	case p == "":
		return "", errRepositoryBaseURLEmpty
	case !allowedCharsRegex.MatchString(p):
		return "", errRepositoryPathInvalidCharsf(p)
	}
	p = "/" + strings.TrimPrefix(strings.TrimSuffix(path.Clean(p), "/"), "/")
	if !validPathRegex.MatchString(p) {
		return "", errRepositoryPathInvalidStructuref(basePath)
	}
	return p, nil
}

// ExtractPathParams validates the template and extracts parameter names, returning a format string with %s placeholders.
func ExtractPathParams(template string) (format string, names []string, err error) {
	sanitized, err := SanitizeAndValidatePath(template)
	if err != nil {
		return "", nil, errRepositoryInvalidPathTemplate(err)
	}
	matches := templateRegex.FindAllStringSubmatch(sanitized, -1)
	if len(matches) == 0 {
		return sanitized, nil, nil // path válido sem parâmetros
	}
	names = make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			names = append(names, m[1])
		}
	}
	return templateRegex.ReplaceAllString(sanitized, "%s"), names, nil
}
