package rest

import (
	"encoding/json"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type restRegistry struct {
	cfg *registryOptions
}

// NewRestRegistry builds a REST registry using the provided options.
func NewRestRegistry(options ...RegistryOption) RestRegistry {
	return &restRegistry{cfg: newRegistryConfig(options...)}
}

func (r *restRegistry) Merge(other RestRegistry) RestRegistry {
	otherRegistry, ok := other.(*restRegistry)
	if !ok || otherRegistry == nil {
		panic("rest registry: Merge requires the same concrete *restRegistry type")
	}
	merged := &restRegistry{cfg: r.cloneConfig()}
	mergeRegistryOptions(merged.cfg, otherRegistry.cfg)
	return merged
}

func (r *restRegistry) ResolveRequest(body RestDataSpec, requestBody *[]byte) error {
	if requestBody == nil {
		return errRestResolveRequestBodyNil
	}
	if body == nil {
		return errRestResolveRequestSpecBodyNil
	}
	return marshalInto(requestBody, body, errRestResolveRequestMarshal)
}

func (r *restRegistry) ResolveEnvelopeRequest(restMethod RestMethod, dataBody *[]byte) error {
	if dataBody == nil {
		return errRestResolveEnvelopeRequestBodyNil
	}
	prototype := r.resolveRequestEnvelopeSpec(restMethod)
	if prototype == nil {
		return nil
	}
	envelope := prototype.New()
	if envelope == nil {
		return errRestResolveEnvelopeRequestSpecNewNil
	}
	envelope.SetEnvelopeData(json.RawMessage(append([]byte(nil), (*dataBody)...)))
	return marshalInto(dataBody, envelope, errRestResolveEnvelopeRequestMarshal)
}

func (r *restRegistry) ResolveResponse(status int, responseBody []byte, body *RestDataSpec) error {
	if body == nil {
		return errRestResolveResponseBodyNil
	}
	instance, err := r.newResponseSpec(status)
	if err != nil {
		return err
	}
	if err := unmarshalInto(responseBody, instance, errRestResolveResponseUnmarshal); err != nil {
		return err
	}
	*body = instance
	return nil
}

func (r *restRegistry) ResolveResponses(status int, responseBody []byte, bodies *[]RestDataSpec) error {
	if bodies == nil {
		return errRestResolveResponsesBodiesNil
	}
	var rawPayloads []json.RawMessage
	if err := unmarshalInto(responseBody, &rawPayloads, errRestResolveResponsesUnmarshalRawList); err != nil {
		return err
	}
	prototype, err := r.resolveResponseSpec(status)
	if err != nil {
		return err
	}
	resolved, err := resolveResponseSlice(prototype, len(rawPayloads))
	if err != nil {
		return err
	}
	for i := range rawPayloads {
		if resolved[i] == nil {
			return errRestResolveResponsesSpecNewNil
		}
		if err := unmarshalInto(rawPayloads[i], resolved[i], func(err error) error {
			return errRestResolveResponsesUnmarshalItem(i, err)
		}); err != nil {
			return err
		}
	}
	*bodies = resolved
	return nil
}

func (r *restRegistry) ResolveEnvelopeResponse(status int, dataBody *[]byte, meta *types.ResponseMeta) error {
	if dataBody == nil {
		return errRestResolveEnvelopeResponseBodyNil
	}
	if meta == nil {
		return errRestResolveEnvelopeResponseMetaNil
	}
	spec := r.resolveResponseEnvelopeSpec(status)
	if spec == nil {
		return nil
	}
	instance := spec.New()
	if instance == nil {
		return errRestResolveEnvelopeResponseSpecNewNil
	}
	if err := unmarshalInto(*dataBody, instance, errRestResolveEnvelopeResponseUnmarshal); err != nil {
		return err
	}
	*dataBody = append([]byte(nil), instance.EnvelopeData()...)
	*meta = instance.EnvelopeMeta()
	return nil
}

func (r *restRegistry) NewRequestSpec(method RestMethod) (RestDataSpec, error) {
	instance, err := r.newRequestSpec(method)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (r *restRegistry) ReleaseRequestSpec(spec RestDataSpec) {
	_ = spec
}
