package mycrawler

import (
	"net/http"
	"net/url"
)

type Handler interface {
	Handle(string, url.URL, *http.Response, error)
}

type HandlerFunc func(string, url.URL, *http.Response, error)

// Handle is the Handler interface implementation for the HandlerFunc type.
func (h HandlerFunc) Handle(method string, url url.URL, res *http.Response, err error) {
	h(method, url, res, err)
}
