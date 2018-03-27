package mock

import (
	"net/http"
	"net/url"
)

type HTTPClient struct {
	DoFunc    func(*http.Request) (*http.Response, error)
	DoInvoked bool

	URL    *url.URL
	Method string
	Header http.Header
}

func (c *HTTPClient) Do(r *http.Request) (*http.Response, error) {
	c.DoInvoked = true
	c.URL = r.URL
	c.Method = r.Method
	c.Header = r.Header

	return c.DoFunc(r)
}
