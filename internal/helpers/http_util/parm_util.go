package http_util

import (
	"maps"
	"net/http"
	"net/url"
	"slices"
)

// ApplyQueryParams merges the provided query values `q` into `u`'s existing
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
func ApplyQueryParams(u *url.URL, q url.Values) {
	if len(q) == 0 {
		return
	}
	existing := u.Query()
	for key, values := range q {
		existing[key] = append(existing[key], values...)
		slices.Sort(existing[key])
		existing[key] = slices.Compact(existing[key])
	}
	u.RawQuery = existing.Encode()
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
func ApplyHeaderParams(dst, src http.Header) {
	if len(src) == 0 {
		return
	}
	maps.Copy(dst, src)
}
