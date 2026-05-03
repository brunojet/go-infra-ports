package net_http

import (
	"net/http"
	"testing"

	svccts "github.com/brunojet/go-infra-ports/pkg/services/contracts"
)

type testNetHTTPHandler struct{}

func (testNetHTTPHandler) Create(http.ResponseWriter, *http.Request) {}
func (testNetHTTPHandler) List(http.ResponseWriter, *http.Request)   {}
func (testNetHTTPHandler) Get(http.ResponseWriter, *http.Request)    {}
func (testNetHTTPHandler) Update(http.ResponseWriter, *http.Request) {}
func (testNetHTTPHandler) Save(http.ResponseWriter, *http.Request)   {}
func (testNetHTTPHandler) Delete(http.ResponseWriter, *http.Request) {}

var _ NetHttpHandler = testNetHTTPHandler{}

func TestNetHTTPTypeAliases_CompileCompatibility(t *testing.T) {
	var h1 NetHttpHandler = testNetHTTPHandler{}
	_ = h1

	a := ServiceCreate[string]{}
	_ = svccts.ServiceCreate[string](a)

	b := ServiceUpdate[string]{}
	_ = svccts.ServiceUpdate[string](b)

	c := ServiceSave[string]{}
	_ = svccts.ServiceSave[string](c)

	d := ServiceResponse[string]{}
	_ = svccts.ServiceResponse[string](d)

	e := ServiceResponses[string]{}
	_ = svccts.ServiceResponses[string](e)
}
