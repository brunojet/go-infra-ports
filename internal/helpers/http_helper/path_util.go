package http_helper

import (
	"path"
	"regexp"
	"strings"
)

const (
	maxPathLen    = 256
	maxPathParams = 8
)

var (
	// allowed chars: letters, digits, underscore, hyphen, braces and slashes
	// Leading slash is optional: the helper will normalize and insert it.
	allowedCharsRegex = regexp.MustCompile(`^/?[A-Za-z0-9_/\-\{\}]*$`)
	// validPathRegex: segments are either plain tokens or {param} (no :param support)
	// Accept root '/' explicitly to match net/http patterns.
	validPathRegex = regexp.MustCompile(`^/$|^(/([A-Za-z0-9_-]+|\{[A-Za-z0-9_]+\}))+$`)
	// templateRegex captures {name} — group 1 is always non-empty (+ quantifier).
	templateRegex = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)
)

// SanitizeAndValidatePath trims whitespace, checks for emptiness, validates allowed characters and structure, and normalizes the path.
func SanitizeAndValidatePath(template string) (string, error) {
	p := strings.TrimSpace(template)
	switch {
	case p == "":
		return "", errRepositoryBasePathEmpty
	case len(p) > maxPathLen:
		return "", errPathTooLongf(p, maxPathLen)
	case !allowedCharsRegex.MatchString(p):
		return "", errPathInvalidCharsf(p)
	}
	p = "/" + strings.TrimPrefix(strings.TrimSuffix(path.Clean(p), "/"), "/")
	if !validPathRegex.MatchString(p) {
		return "", errPathInvalidStructuref(template)
	}
	return p, nil
}

func PathParamsFmt(template string) (string, error) {
	sanitized, err := SanitizeAndValidatePath(template)
	if err != nil {
		return "", errInvalidPathTemplate(err)
	}
	return templateRegex.ReplaceAllString(sanitized, "%s"), nil
}

// ExtractPathParams validates the template and extracts parameter names, returning a format string with %s placeholders.
func ExtractPathParams(template string) (names []string, err error) {
	sanitized, err := SanitizeAndValidatePath(template)
	if err != nil {
		return nil, errInvalidPathTemplate(err)
	}
	matches := templateRegex.FindAllStringSubmatch(sanitized, -1)
	lenMatches := len(matches)
	if lenMatches > maxPathParams {
		return nil, errPathParametersExceedMaxf(lenMatches, maxPathParams)
	}
	if lenMatches == 0 {
		return nil, nil
	}
	names = make([]string, 0, lenMatches)
	for _, m := range matches {
		names = append(names, m[1]) // m[1] is always non-empty, guaranteed by templateRegex
	}
	return names, nil
}
