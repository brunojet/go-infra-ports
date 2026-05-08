package rest

import (
	"encoding/json"
	"maps"
)

func mergeRegistryOptions(dst, src *registryOptions) {
	maps.Copy(dst.requests, src.requests)
	maps.Copy(dst.requestsEnvelopes, src.requestsEnvelopes)
	maps.Copy(dst.responses, src.responses)
	maps.Copy(dst.responseEnvelopes, src.responseEnvelopes)
	maps.Copy(dst.informations, src.informations)
	maps.Copy(dst.redirections, src.redirections)
	maps.Copy(dst.problems, src.problems)
}

func resolveResponseInstance(prototype RestDataSpec) (RestDataSpec, error) {
	resolved := prototype.New()
	if resolved == nil {
		return nil, errRestResolveResponseSpecNewNil
	}
	return resolved, nil
}

func resolveResponseSlice(prototype RestDataSpec, n int) ([]RestDataSpec, error) {
	resolved := prototype.NewSlice(n)
	if len(resolved) != n {
		return nil, errRestResolveResponsesNewSliceLen(len(resolved), n)
	}
	return resolved, nil
}

func resolveResponseSpec(status int, source map[int]RestDataSpec) RestDataSpec {
	if spec, ok := source[status]; ok && spec != nil {
		return spec
	}
	return source[defaultStatusCode]
}

func marshalInto(target *[]byte, value any, wrapErr func(error) error) error {
	data, err := json.Marshal(value)
	if err != nil {
		return wrapErr(err)
	}
	*target = data
	return nil
}

func unmarshalInto(data []byte, target any, wrapErr func(error) error) error {
	if err := json.Unmarshal(data, target); err != nil {
		return wrapErr(err)
	}
	return nil
}
