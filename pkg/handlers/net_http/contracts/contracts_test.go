package contracts

import (
	"net/http"
	"testing"
)

type testHandler struct{}

func (testHandler) Create(http.ResponseWriter, *http.Request) {}
func (testHandler) List(http.ResponseWriter, *http.Request)   {}
func (testHandler) Get(http.ResponseWriter, *http.Request)    {}
func (testHandler) Update(http.ResponseWriter, *http.Request) {}
func (testHandler) Save(http.ResponseWriter, *http.Request)   {}
func (testHandler) Delete(http.ResponseWriter, *http.Request) {}

var _ NetHttpHandler = testHandler{}

func TestNetHttpHandlerContract_Implementation(t *testing.T) {
	var h NetHttpHandler = testHandler{}
	_ = h
}
