package fasthttp

import (
	"github.com/valyala/fasthttp"
	"github.com/webx-top/echo/engine"
	"github.com/xupingao/go-easy-adapt/http"
	"net/url"
	std "net/http"
)

type argValues struct {
	args   *fasthttp.Args
	values *url.Values
}

func (u *argValues) Add(key string, value string) {
	u.args.Add(key, value)
}

func (u *argValues) Set(key string, value string) {
	u.args.Set(key, value)
}
func (u *argValues) Del(key string) {
	u.args.Del(key)
}

func (u *argValues) Get(key string) string {
	return engine.Bytes2str(u.args.Peek(key))
}

func (u *argValues) Gets(key string) []string {
	v := engine.Bytes2str(u.args.Peek(key))
	if len(v) != 0 {
		return []string{v}
	}
	return []string{}
}

func (v *argValues) All() map[string][]string {
	if v.values != nil {
		return *v.values
	}
	r := url.Values{}
	v.args.VisitAll(func(k, v []byte) {
		key := engine.Bytes2str(k)
		r.Add(key, engine.Bytes2str(v))
	})
	v.values = &r
	return *v.values
}

type formValues struct {
	rawValues *url.Values
}

func (v *formValues) Add(key string, value string) {
	v.rawValues.Add(key, value)
}

func (v *formValues) Del(key string) {
	v.rawValues.Del(key)
}

func (v *formValues) Get(key string) string {
	return v.rawValues.Get(key)
}

func (v *formValues) Gets(key string) []string {
	form := *v.rawValues
	if v, ok := form[key]; ok {
		return v
	}
	return []string{}
}

func (v *formValues) Set(key string, value string) {
	v.rawValues.Set(key, value)
}

func (v *formValues) All() map[string][]string {
	return *(v.rawValues)
}

type URL struct {
	url   *fasthttp.URI
	query http.Values
}

func (u *URL) SetPath(path string) {
	u.url.SetPath(path)
}

func (u *URL) RawPath() string {
	return engine.Bytes2str(u.url.PathOriginal())
}

func (u *URL) Path() string {
	return engine.Bytes2str(u.url.Path())
}

func (u *URL) QueryValue(name string) string {
	return engine.Bytes2str(u.url.QueryArgs().Peek(name))
}

func (u *URL) QueryValues(name string) []string {
	u.Query()
	return u.query.Gets(name)
}

func (u *URL) Query() http.Values {
	if u.query == nil {
		u.query = &argValues{u.url.QueryArgs(), nil}
	}
	return u.query
}

func (u *URL) RawQuery() string {
	return engine.Bytes2str(u.url.QueryString())
}

func (u *URL) SetRawQuery(rawQuery string) {
	u.url.SetQueryString(rawQuery)
}

func (u *URL) String() string {
	return u.url.String()
}

func (u *URL) Object() interface{} {
	return u.url
}

func (u *URL) reset(url *fasthttp.URI) {
	u.url = url
}

type (
	RequestHeader struct {
		header *fasthttp.RequestHeader
		stdhdr *std.Header
	}

	ResponseHeader struct {
		header *fasthttp.ResponseHeader
		stdhdr *std.Header
	}
)

func (h *RequestHeader) Add(key, val string) {
	h.header.Set(key, val)
}

func (h *RequestHeader) Del(key string) {
	h.header.Del(key)
}

func (h *RequestHeader) Get(key string) string {
	return engine.Bytes2str(h.header.Peek(key))
}

func (h *RequestHeader) Set(key, val string) {
	h.header.Set(key, val)
}
func (h *RequestHeader) All()map[string][]string {
	if h.stdhdr != nil {
		return *h.stdhdr
	}
	hdr := std.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	return hdr
}


func (h *RequestHeader) Object() interface{} {
	return h.header
}

func (h *ResponseHeader) Add(key, val string) {
	h.header.Set(key, val)
}

func (h *RequestHeader) reset(hdr *fasthttp.RequestHeader) {
	h.header = hdr
}
func (h *ResponseHeader) Del(key string) {
	h.header.Del(key)
}

func (h *ResponseHeader) Get(key string) string {
	return engine.Bytes2str(h.header.Peek(key))
}

func (h *ResponseHeader) Set(key, val string) {
	h.header.Set(key, val)
}

func (h *ResponseHeader) Object() interface{} {
	return h.header
}

func (h *ResponseHeader) reset(hdr *fasthttp.ResponseHeader) {
	h.header = hdr
}

func (h *ResponseHeader) All()map[string][]string {
	if h.stdhdr != nil {
		return *h.stdhdr
	}
	hdr := std.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	return hdr
}