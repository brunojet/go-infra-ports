package http_helper

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
)

// ApplyURLQueryParams merges the provided query values `q` into `u`'s existing
// query string in-place and encodes the result into `u.RawQuery`.
// Behavior notes:
//   - For each key, incoming values are appended to existing values, then the
//     resulting slice is sorted and de-duplicated using `slices.Sort` and
//     `slices.Compact`.
//   - Because values are sorted, insertion order is NOT preserved. Use this
//     helper when a canonical, deduplicated representation is desired (for
//     deterministic requests or cache keys). If you need to preserve order,
//     merge values manually without sorting.
//   - Complexity: per-key sort dominates (O(m log m) for m values of a key).
//   - The function modifies `u` in-place.
func ApplyURLQueryParams(u *url.URL, q url.Values) error {
	if u == nil {
		return errNilURL
	}
	existing := u.Query()
	if err := ApplyQueryParams(existing, q); err != nil {
		return fmt.Errorf("failed to apply query params: %w", err)
	}
	u.RawQuery = existing.Encode()
	return nil
}

func ApplyQueryParams(dst, src url.Values) error {
	if dst == nil {
		return fmt.Errorf("destination query values cannot be nil")
	}
	if len(src) == 0 {
		return nil
	}
	for key, values := range src {
		dst[key] = append(dst[key], values...)
		slices.Sort(dst[key])
		dst[key] = slices.Compact(dst[key])
	}
	return nil
}

// ApplyHeaderParams merges headers from `src` into `dst`.
// Implementation details and policy:
//
//   - Uses `maps.Copy(dst, src)`, which performs a shallow copy of the header
//     map entries. The header values are slices of strings and their backing
//     arrays are NOT deep-cloned by `maps.Copy`.
//
//   - Shallow copying is intentional for performance: it avoids allocating new
//     slices on every hop. If callers treat headers as immutable after merge,
//     this is safe and efficient.
//
//   - If a caller intends to mutate header value slices after merging, they
//     MUST deep-clone the values to avoid aliasing. Prefer the stdlib
//     `http.Header.Clone()` at edge boundaries (for example, when a handler
//     receives a transaction or just before sending an HTTP response) to
//     obtain a safe, independent copy. Example using stdlib:
//
//     cloned := src.Clone()
//
//   - For query parameters (`url.Values`) there is no stdlib `Clone()` helper;
//     when cloning at an edge boundary use an explicit per-key copy, e.g.:
//
//     q2 := make(url.Values, len(q))
//     for k, vs := range q { q2[k] = append([]string(nil), vs...) }
//
//   - `src` takes precedence on key conflicts (it will overwrite keys in `dst`).
func ApplyHeaderParams(dst, src http.Header) error {
	if dst == nil {
		return errNilDstHeader
	}
	if len(src) == 0 {
		return nil
	}
	maps.Copy(dst, src)
	return nil
}
