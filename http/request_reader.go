package http

import (
	"io"
	"net/http"
	"net/url"
)

type RequestReader interface{
	Cookie(name string) (*http.Cookie, error)
	ProtoAtLeast(major, minor int) bool
	UserAgent() string
	Cookies() []*http.Cookie
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() url.Values
}
