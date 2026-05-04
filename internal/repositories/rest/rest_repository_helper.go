package rest

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strings"
)

var (
	allowedCharsRegex = regexp.MustCompile(`^/[A-Za-z0-9_:/\-\{\}]*$`)
	validPathRegex    = regexp.MustCompile(`^(/([A-Za-z0-9_-]+|:[A-Za-z0-9_]+|\{[A-Za-z0-9_]+\}))+$`)
	templateRegex     = regexp.MustCompile(`:([A-Za-z0-9_]+)|\{([A-Za-z0-9_]+)\}`)
)

func sanitizeAndValidatePath(basePath string) (string, error) {
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

func extractPathParams(template string) (format string, names []string, err error) {
	sanitized, err := sanitizeAndValidatePath(template)
	if err != nil {
		return "", nil, errRepositoryInvalidPathTemplate(err)
	}
	matches := templateRegex.FindAllStringSubmatch(sanitized, -1)
	if len(matches) == 0 {
		return sanitized, nil, nil // path válido sem parâmetros
	}
	names = make([]string, 0, len(matches))
	for _, m := range matches {
		switch {
		case m[1] != "":
			names = append(names, m[1])
		case m[2] != "":
			names = append(names, m[2])
		}
	}
	return templateRegex.ReplaceAllString(sanitized, "%s"), names, nil
}

// applyQueryParams appends q onto u's existing query string, skipping exact key+value duplicates.
func applyQueryParams(u *url.URL, q url.Values) {
	if len(q) == 0 {
		return
	}
	existing := u.Query()
	for key, values := range q {
		current := make(map[string]struct{}, len(existing[key]))
		for _, v := range existing[key] {
			current[v] = struct{}{}
		}
		for _, value := range values {
			if _, dup := current[value]; !dup {
				existing.Add(key, value)
				current[value] = struct{}{}
			}
		}
	}
	u.RawQuery = existing.Encode()
}

// applyHeaderParams copies src headers into dst, overriding existing keys.
func applyHeaderParams(dst, src http.Header) {
	if len(src) == 0 {
		return
	}
	for key, values := range src {
		dst[key] = slices.Clone(values)
	}
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}
