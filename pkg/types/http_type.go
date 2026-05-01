package types

import (
	"net/http"
	"net/url"
)

type HTTPRequesOptions struct {
	Header http.Header
	Query  url.Values
}

type HTTPResponseOptions struct {
	StatusCode int
	Header     http.Header
}
