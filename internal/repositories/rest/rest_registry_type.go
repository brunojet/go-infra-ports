package rest

import (
	"encoding/json"
	"fmt"

	"github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type (
	RestRegistry     = contracts.RestRegistry
	RestEnvelopeSpec = contracts.RestEnvelopeSpec
	RestDataSpec     = contracts.RestDataSpec
	RestMethod       = contracts.RestMethod
)

const (
	MethodCreate                    = contracts.MethodCreate
	MethodList                      = contracts.MethodList
	MethodGet                       = contracts.MethodGet
	MethodUpdate                    = contracts.MethodUpdate
	MethodSave                      = contracts.MethodSave
	MethodDelete                    = contracts.MethodDelete
	allWriteMethods      RestMethod = MethodCreate | MethodUpdate | MethodSave
	allCollectionMethods RestMethod = MethodCreate | MethodList
	allInstanceMethods   RestMethod = MethodGet | MethodUpdate | MethodSave | MethodDelete
)

var (
	writeMethodsList = []RestMethod{MethodCreate, MethodUpdate, MethodSave}
)

type DefaultRestRequest = json.RawMessage

type DefaultRestResponse = json.RawMessage

// dataSpecOf implements the common behavior shared by request and response specs.
type dataSpecOf[T any] struct {
	Prototype T
}

func NewDataSpecOf[T any]() RestDataSpec {
	return &dataSpecOf[T]{}
}

func (d *dataSpecOf[T]) New() RestDataSpec {
	return &dataSpecOf[T]{}
}

func (d *dataSpecOf[T]) NewSlice(n int) []RestDataSpec {
	backarray := make([]dataSpecOf[T], n)
	outarray := make([]RestDataSpec, n)
	for i := range backarray {
		outarray[i] = &backarray[i]
	}
	return outarray
}

func (d *dataSpecOf[T]) Body() any {
	return d.Prototype
}

func (d *dataSpecOf[T]) SetBody(body any) error {
	if typed, ok := body.(T); ok {
		d.Prototype = typed
		return nil
	}
	return fmt.Errorf("invalid body type: expected %T, got %T", d.Prototype, body)
}

func (d *dataSpecOf[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Prototype)
}

func (d *dataSpecOf[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d.Prototype)
}

type envelopeSpecOf struct {
	Data      json.RawMessage
	Meta      types.ResponseMeta
	DataField string
	MetaField string
}

func NewEnvelopeSpec(dataField, metaField string) RestEnvelopeSpec {
	if dataField == "" {
		dataField = "data"
	}
	return &envelopeSpecOf{
		DataField: dataField,
		MetaField: metaField,
	}
}

func (r *envelopeSpecOf) New() RestEnvelopeSpec {
	return &envelopeSpecOf{
		DataField: r.DataField,
		MetaField: r.MetaField,
	}
}

func (r *envelopeSpecOf) MarshalJSON() ([]byte, error) {
	envelope := map[string]any{
		r.DataField: r.Data,
	}
	if r.MetaField != "" {
		envelope[r.MetaField] = r.Meta
	}
	return json.Marshal(envelope)
}

func (r *envelopeSpecOf) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if d, ok := raw[r.DataField]; ok {
		r.Data = d
	} else {
		return fmt.Errorf("missing data field %q in response envelope", r.DataField)
	}
	if r.MetaField == "" {
		return nil
	}
	if m, ok := raw[r.MetaField]; ok {
		var meta map[string]any
		if err := json.Unmarshal(m, &meta); err != nil {
			return err
		}
		r.Meta = meta
	} else {
		return fmt.Errorf("missing meta field %q in response envelope", r.MetaField)
	}
	return nil
}

func (r *envelopeSpecOf) EnvelopeData() json.RawMessage {
	if r.Data == nil {
		return json.RawMessage("null")
	}
	return r.Data
}

func (r *envelopeSpecOf) EnvelopeMeta() types.ResponseMeta {
	if r.Meta == nil {
		return types.ResponseMeta{}
	}
	return r.Meta
}

func (r *envelopeSpecOf) SetEnvelopeData(data json.RawMessage) {
	r.Data = data
}
