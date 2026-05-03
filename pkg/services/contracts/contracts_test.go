package contracts

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-ports/pkg/types"
)

type testService struct{}

func (testService) Create(context.Context, ServiceCreate[string], *ServiceResponse[string]) error {
	return nil
}
func (testService) List(context.Context, types.RequestContext, *ServiceResponses[string]) error {
	return nil
}
func (testService) Get(context.Context, types.RequestContext, *ServiceResponse[string]) error {
	return nil
}
func (testService) Update(context.Context, ServiceUpdate[string], *ServiceResponse[string]) error {
	return nil
}
func (testService) Save(context.Context, ServiceSave[string], *ServiceResponse[string]) error {
	return nil
}
func (testService) Delete(context.Context, types.RequestContext, *ServiceResponse[string]) error {
	return nil
}

var _ Service[string, string, string] = testService{}

func TestServiceContracts_BasicStructuresAndInterface(t *testing.T) {
	create := ServiceCreate[string]{Body: "create"}
	update := ServiceUpdate[string]{Body: "update"}
	save := ServiceSave[string]{Body: "save"}
	one := ServiceResponse[string]{Data: "one"}
	many := ServiceResponses[string]{Data: []string{"a", "b"}}

	if create.Body != "create" || update.Body != "update" || save.Body != "save" {
		t.Fatal("unexpected payload values in service input structs")
	}
	if one.Data != "one" || len(many.Data) != 2 {
		t.Fatal("unexpected payload values in service output structs")
	}

	var svc Service[string, string, string] = testService{}
	_ = svc
}
