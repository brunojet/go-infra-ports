package rest

import (
	"context"
	"testing"

	repcts "github.com/brunojet/go-infra-ports/pkg/repositories/rest/contracts"
	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
	"github.com/brunojet/go-infra-ports/pkg/types"
)

type typeAliasRepo struct{}

func (typeAliasRepo) Create(context.Context, repcts.RestRequest, *repcts.RestResponse) error {
	return nil
}
func (typeAliasRepo) List(context.Context, types.RequestContext, *repcts.RestResponses) error {
	return nil
}
func (typeAliasRepo) Get(context.Context, types.RequestContext, *repcts.RestResponse) error {
	return nil
}
func (typeAliasRepo) Update(context.Context, repcts.RestRequest, *repcts.RestResponse) error {
	return nil
}
func (typeAliasRepo) Save(context.Context, repcts.RestRequest, *repcts.RestResponse) error {
	return nil
}
func (typeAliasRepo) Delete(context.Context, types.RequestContext, *repcts.RestResponse) error {
	return nil
}

var _ RestRepository = typeAliasRepo{}

func TestRestServiceTypeAliases_CompileCompatibility(t *testing.T) {
	var repo RestRepository = typeAliasRepo{}
	_ = repo

	createReq := ServiceCreate[string]{Body: "x"}
	_ = svccts.ServiceCreate[string](createReq)

	updateReq := ServiceUpdate[string]{Body: "y"}
	_ = svccts.ServiceUpdate[string](updateReq)

	saveReq := ServiceSave[string]{Body: "z"}
	_ = svccts.ServiceSave[string](saveReq)

	resp := ServiceResponse[string]{}
	_ = svccts.ServiceResponse[string](resp)

	resps := ServiceResponses[string]{}
	_ = svccts.ServiceResponses[string](resps)
}
