package contracts

import (
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type HTTPRequesOptions = types.HTTPRequesOptions
type HTTPResponseOptions = types.HTTPResponseOptions

type RestRequest struct {
	HTTPRequesOptions
	Body any
}

type BaseRestResponse struct {
	HTTPResponseOptions
	Information any
	Redirection any
	Error       any
}

type RestResponse struct {
	BaseRestResponse
	Body any
}

type RestResonses struct {
	BaseRestResponse
	Bodies []any
}

type RestRepository interface {
	Create(request RestRequest, response *RestResponse) error
	List(opts HTTPRequesOptions, response *RestResonses) error
	Get(opts HTTPRequesOptions, response *RestResponse) error
	Update(request RestRequest, response *RestResponse) error
	Save(request RestRequest, response *RestResponse) error
	Delete(opts HTTPRequesOptions, response *RestResponse) error
}
