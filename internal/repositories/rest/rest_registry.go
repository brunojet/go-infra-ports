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
	merged := &restRegistry{cfg: r.cloneConfig()}
	if other == nil {
		return merged
	}
	otherRegistry, ok := other.(*restRegistry)
	if !ok || otherRegistry == nil {
		return merged
	}
	mergeRegistryOptions(merged.cfg, otherRegistry.cfg)
	return merged
}

func (r *restRegistry) ResolveRequest(body RestRequestSpec, requestBody *[]byte) error {
	if requestBody == nil {
		return errRestResolveRequestBodyNil
	}
	if body == nil {
		return errRestResolveRequestSpecBodyNil
	}
	return marshalInto(requestBody, body, errRestResolveRequestMarshal)
}

func (r *restRegistry) ResolveEnvelopeRequest(name string, dataBody *[]byte) error {
	if dataBody == nil {
		return errRestResolveEnvelopeRequestBodyNil
	}
	prototype := r.resolveRequestEnvelopeSpec(name)
	if prototype == nil {
		return nil
	}
	envelope := prototype.New()
	if envelope == nil {
		return errRestResolveEnvelopeRequestSpecNewNil
	}
	envelope.SetBody(&DefaultRestRequest{Body: append([]byte(nil), (*dataBody)...)})
	return marshalInto(dataBody, envelope, errRestResolveEnvelopeRequestMarshal)
}

func (r *restRegistry) ResolveResponse(status int, responseBody []byte, body *RestResponseSpec) error {
	if body == nil {
		return errRestResolveResponseBodyNil
	}
	prototype, err := r.resolveResponseSpec(status)
	if err != nil {
		return err
	}
	instance, err := resolveResponseInstance(prototype)
	if err != nil {
		return err
	}
	if err := unmarshalInto(responseBody, instance, errRestResolveResponseUnmarshal); err != nil {
		return err
	}
	*body = instance
	return nil
}

func (r *restRegistry) ResolveResponses(status int, responseBody []byte, bodies *[]RestResponseSpec) error {
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

func (r *restRegistry) NewRequestSpec(name string) (RestRequestSpec, error) {
	prototype := r.resolveRequestSpec(name)
	instance := prototype.New()
	if instance == nil {
		return nil, errRestNewRequestSpecNil
	}
	return instance, nil
}

func (r *restRegistry) ReleaseRequestSpec(spec RestRequestSpec) {
	_ = spec
}
