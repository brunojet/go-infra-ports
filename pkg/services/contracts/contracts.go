package contracts

import (
	"encoding/json"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type HTTPRequesOptions = types.HTTPRequesOptions
type HTTPResponseOptions = types.HTTPResponseOptions

type ServiceCreate[C any] struct {
	HTTPRequesOptions
	Body C
}

type ServiceUpdate[U any] struct {
	HTTPRequesOptions
	Body U
}

type ServiceSave[C any] struct {
	HTTPRequesOptions
	Body C
}

type ServiceMeta struct {
	Message string
	Details json.RawMessage
}

type BaseServiceResponse struct {
	HTTPResponseOptions
	Meta ServiceMeta
}

type ServiceResponse[R any] struct {
	BaseServiceResponse
	Body R
}

type ServiceResonses[R any] struct {
	BaseServiceResponse
	Bodies []R
}

type Service[C, R, U any] interface {
	Create(request ServiceCreate[C], response *ServiceResponse[R]) error
	List(opts HTTPRequesOptions, response *ServiceResonses[R]) error
	Get(opts HTTPRequesOptions, response *ServiceResponse[R]) error
	Update(request ServiceUpdate[U], response *ServiceResponse[R]) error
	Save(request ServiceSave[C], response *ServiceResponse[R]) error
	Delete(opts HTTPRequesOptions, response *ServiceResponse[R]) error
}
