package rest

import (
	"encoding/json"
	"maps"
)

func mergeRegistryOptions(dst, src *registryOptions) {
	maps.Copy(dst.Requests, src.Requests)
	maps.Copy(dst.RequestsEnvelopes, src.RequestsEnvelopes)
	maps.Copy(dst.Responses, src.Responses)
	maps.Copy(dst.ResponseEnvelopes, src.ResponseEnvelopes)
	maps.Copy(dst.Informations, src.Informations)
	maps.Copy(dst.Redirections, src.Redirections)
	maps.Copy(dst.Problems, src.Problems)
}

func resolveResponseInstance(prototype RestResponseSpec) (RestResponseSpec, error) {
	resolved := prototype.New()
	if resolved == nil {
		return nil, errRestResolveResponseSpecNewNil
	}
	return resolved, nil
}

func resolveResponseSlice(prototype RestResponseSpec, n int) ([]RestResponseSpec, error) {
	resolved := prototype.NewSlice(n)
	if len(resolved) != n {
		return nil, errRestResolveResponsesNewSliceLen(len(resolved), n)
	}
	return resolved, nil
}

func resolveResponseSpec(status int, source map[int]RestResponseSpec) RestResponseSpec {
	if spec, ok := source[status]; ok && spec != nil {
		return spec
	}
	return source[DefaultStatusCode]
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
