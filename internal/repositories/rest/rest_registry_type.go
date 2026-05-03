package rest

import (
	"encoding/json"

	"github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

const (
	DefaultMethodName = "default"
	DefaultStatusCode = 0
)

type (
	RestRegistry     = contracts.RestRegistry
	RestRequestSpec  = contracts.RestRequestSpec
	RestResponseSpec = contracts.RestResponseSpec
	RestEnvelopeSpec = contracts.RestEnvelopeSpec
	RestRequest      = contracts.RestRequest
	RestResponse     = contracts.RestResponse
	RestResponses    = contracts.RestResponses
	RestRepository   = contracts.RestRepository
)

type DefaultRestRequest struct {
	Body json.RawMessage
}

func (d *DefaultRestRequest) New() RestRequestSpec {
	return &DefaultRestRequest{}
}

func (d *DefaultRestRequest) SetBody(body RestRequestSpec) {
	if body == nil {
		d.Body = nil
		return
	}
	if raw, ok := body.(*DefaultRestRequest); ok {
		d.Body = append([]byte(nil), raw.Body...)
		return
	}
	data, err := json.Marshal(body)
	if err == nil {
		d.Body = data
	}
}

type DefaultRestResponse struct {
	Body json.RawMessage
}

func (d *DefaultRestResponse) New() RestResponseSpec {
	return &DefaultRestResponse{}
}

func (d *DefaultRestResponse) UnmarshalJSON(data []byte) error {
	d.Body = append([]byte(nil), data...)
	return nil
}

func (d *DefaultRestResponse) NewSlice(n int) []RestResponseSpec {
	backarray := make([]DefaultRestResponse, n)
	outarray := make([]RestResponseSpec, n)
	for i := range backarray {
		outarray[i] = &backarray[i]
	}
	return outarray
}

type responseEnvelope struct {
	Data json.RawMessage    `json:"data"`
	Meta types.ResponseMeta `json:"meta"`
}

func (e *responseEnvelope) New() RestEnvelopeSpec {
	return &responseEnvelope{}
}

func (e *responseEnvelope) EnvelopeData() json.RawMessage {
	if e.Data == nil {
		return json.RawMessage("null")
	}
	return e.Data
}

func (e *responseEnvelope) EnvelopeMeta() types.ResponseMeta {
	if e.Meta == nil {
		return types.ResponseMeta{}
	}
	return e.Meta
}
