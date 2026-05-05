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
	// allowed chars: letters, digits, underscore, hyphen, dot, braces and slashes.
	// Leading slash is optional: the helper will normalize and insert it.
	// Dot is allowed here but traversal segments (./ and ../) are rejected by containsTraversal.
	allowedCharsRegex = regexp.MustCompile(`^/?[A-Za-z0-9_/\-\{\}.]*$`)
	// validPathRegex: segments are plain tokens (may contain dot but cannot start with it) or {param}.
	// Accept root '/' explicitly to match net/http patterns.
	validPathRegex = regexp.MustCompile(`^/$|^(/([A-Za-z0-9_-][A-Za-z0-9_\-.]*|\{[A-Za-z0-9_]+\}))+$`)
	// templateRegex captures {name} — group 1 is always non-empty (+ quantifier).
	// validPathRegex rejects empty braces before this point.
	templateRegex = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)
)

// containsTraversal reports whether p contains any . or .. segment.
// Must be called before path.Clean, which would silently resolve traversal.
func containsTraversal(p string) bool {
	for _, seg := range strings.Split(p, "/") {
		if seg == "." || seg == ".." {
			return true
		}
	}
	return false
}

// SanitizeAndValidatePath trims whitespace, checks for emptiness, validates
// allowed characters and structure, rejects traversal, and normalizes the path.
func SanitizeAndValidatePath(template string) (string, error) {
	p := strings.TrimSpace(template)
	switch {
	case p == "":
		return "", errRepositoryBasePathEmpty
	case len(p) > maxPathLen:
		return "", errPathTooLongf(p, maxPathLen)
	case !allowedCharsRegex.MatchString(p):
		return "", errPathInvalidCharsf(p)
	case containsTraversal(p):
		return "", errPathTraversalf(p)
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

// ExtractPathParams validates the template and extracts parameter names.
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
