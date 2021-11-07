package fasthttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/valyala/fasthttp"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"strings"
	"unsafe"

	"github.com/webx-top/echo/logger"
	"github.com/xupingao/go-easy-adapt/http"
)

var _ http.Context = context_{}

func wrapConext(requestCtx *fasthttp.RequestCtx) http.Context {
	return context_{
		rawRequestCtx: requestCtx,
		request:       &httpRequest{RequestCtx: requestCtx},
		response:      &httpResponse{RequestCtx: requestCtx},
	}
}

type context_ struct {
	rawRequestCtx *fasthttp.RequestCtx

	request  http.HTTPRequest
	response http.HTTPResponse
}

func (c context_) Request() http.HTTPRequest {
	return c.request
}

func (c context_) Response() http.HTTPResponse {
	return c.response
}

func (c context_) Redirect(code int, url string) {
	c.rawRequestCtx.Redirect(url, code)
}

///////////////////////////////////////////////////////////////
type httpRequest struct {
	*fasthttp.RequestCtx
	url    http.URL
	header http.Header

	query    http.Values
	postForm http.Values
	form     http.Values
}

// func (r *httpRequest) Context() context.Context {
// 	panic("")
// }
// func (r *httpRequest) WithContext(ctx context.Context) http.HTTPRequest {
// 	panic("")
// }

func (r *httpRequest)initValues() {
	if r.query != nil {
		return
	}
	r.query = http.Values{}
	r.RequestCtx.QueryArgs().VisitAll(func(k, v []byte) {
		r.query.Add(Bytes2str(k),Bytes2str(v))
	})

	r.postForm = http.Values{}
	r.RequestCtx.PostArgs().VisitAll(func(k, v []byte) {
		r.postForm.Add(Bytes2str(k),Bytes2str(v))
	})
	if r.form != nil {
		return
	}
	r.form = http.Values{}
	for key, vals := range r.query {
		r.form[key] = vals
	}
	for key, vals := range r.postForm {
		r.form[key] = vals
	}
	mf:= r.MultipartForm()
	if  mf != nil && mf.Value != nil {
		for key, vals := range mf.Value {
			r.form[key] = vals
		}
	}
}

func (r *httpRequest) Clone(ctx context.Context) *http.HTTPRequest {
	panic("")
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (r *httpRequest) Scheme() string {
	if r.IsTLS() {
		return http.SchemeHTTPS
	}
	scheme := r.RequestCtx.URI().Scheme()
	if len(scheme) > 0 {
		return Bytes2str(scheme)
	}
	return http.SchemeHTTP
}

func (r *httpRequest) Proto() string {
	return "HTTP/1.1"
}
func (r *httpRequest) Host() string {
	return Bytes2str(r.RequestCtx.Host())
}
func (r *httpRequest) SetHost(h string) {
	r.RequestCtx.Request.Header.SetHost(h)
}
func (r *httpRequest) URL() http.URL {
	if r.url == nil {
		r.url = &URL{url: r.RequestCtx.URI()}
	}
	return r.url
}
func (r *httpRequest) URI() string {
	return Bytes2str(r.RequestCtx.RequestURI())
}
func (r *httpRequest) SetURI(uri string) {
	r.RequestCtx.Request.Header.SetRequestURI(uri)
}

func (r *httpRequest) Method() string {
	return Bytes2str(r.RequestCtx.Method())
}
func (r *httpRequest) SetMethod(method string) {
	r.RequestCtx.Request.Header.SetMethod(method)
}

func (r *httpRequest) Header() http.Header {
	if r.header == nil {
		r.header = &RequestHeader{header: &r.RequestCtx.Request.Header}
	}
	return r.header
}

func (r *httpRequest) RemoteAddr() string {
	panic("")
}

// func (r *httpRequest) RealIP() string {
// 	if len(r.realIP) > 0 {
// 		return r.realIP
// 	}

// 	r.realIP = realip.XRealIP(r.header.Get(echo.HeaderXRealIP), r.header.Get(echo.HeaderXForwardedFor), r.RemoteAddress())
// 	return r.realIP
// }

func (r *httpRequest) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(r.RequestCtx.PostBody()))

}
func (r *httpRequest) SetBody(reader io.Reader) {
	r.RequestCtx.Request.SetBodyStream(reader, 0)
}

func (r *httpRequest) Form() http.Values {
	r.initValues()
	return r.form
}
func (r *httpRequest) FormValue(string) string {
	panic("")
}
func (r *httpRequest) PostForm() http.Values {
	r.initValues()
	return r.postForm
}

// MultipartForm returns the multipart form.
func (r *httpRequest) MultipartForm() *multipart.Form {
	if !strings.HasPrefix(string(r.RequestCtx.Request.Header.ContentType()), echo.MIMEMultipartForm) {
		return nil
	}
	re, _ := r.RequestCtx.MultipartForm()
	//if err != nil {
	//	r.context.Logger().Printf(err.Error())
	//}
	return re
}
func (r *httpRequest) Referer() string {
	return Bytes2str(r.RequestCtx.Referer())
}
func (r *httpRequest) UserAgent() string {
	return Bytes2str(r.RequestCtx.UserAgent())
}
func (r *httpRequest) Size() int64 {
	return int64(r.RequestCtx.Request.Header.ContentLength())
}

func (r *httpRequest) Cookies() []*http.Cookie {
	panic("")
}
func (r *httpRequest) Cookie(name string) (*http.Cookie, error) {
	panic("")
}
func (r *httpRequest) AddCookie(c *http.Cookie) {
	panic("")
}
func (r *httpRequest) Query() http.Values {
	panic("")
}

func (r *httpRequest) TransferEncoding() []string {
	panic("")
}
func (r *httpRequest) Trailer() http.Header {
	panic("")
}

func (r *httpRequest) MultipartReader() (*multipart.Reader, error) {
	panic("")
}

func (r *httpRequest) IsTLS() bool {
	return r.RequestCtx.IsTLS()
}
func (r *httpRequest) TLS() *tls.ConnectionState {
	panic("")
}

/////////////////////////////////////////////////////////////////
const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

var _ http.HTTPResponse = &httpResponse{}

type httpResponse struct {
	*fasthttp.RequestCtx

	logger logger.Logger
	writen bool
	status int
	size   int64
	header http.Header
}

func (r *httpResponse) WriteHeader(statusCode int) {
	if statusCode > 0 && r.status != statusCode {
		if r.Written() {
			r.logger.Warn("response already committed")
		}
		r.status = statusCode
	}
}

func (r *httpResponse) WriteHeaderNow() {
	if !r.Written() {
		r.writen = true
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.RequestCtx.SetStatusCode(r.status)
	}
}

func (r *httpResponse) Written() bool {
	return r.writen
}

func (r *httpResponse) Write(b []byte) (n int, err error) {
	r.WriteHeaderNow()
	n, err = r.RequestCtx.Write(b)
	r.size += int64(n)
	return
}

func (r *httpResponse) WriteString(s string) (n int, err error) {
	r.WriteHeaderNow()
	n, err = io.WriteString(r.RequestCtx, s)
	r.size += int64(n)
	return
}

func (r *httpResponse) Status() int {
	return r.status
}
func (r *httpResponse) Size() int64 {
	return r.size
}
func (r *httpResponse) Hijacker(fn func(net.Conn)) error {
	r.RequestCtx.Hijack(fasthttp.HijackHandler(fn))
	r.writen = true
	return nil
}

func (r *httpResponse) Header() http.Header {
	if r.header == nil {
		r.header = &ResponseHeader{header: &r.RequestCtx.Response.Header}
	}
	return r.header
}

func (r *httpResponse) HeaderWritten() bool {
	panic("")
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
	if v, ok := u.query[name]; ok {
		return v
	}
	return []string{}
}

func (u *URL) Query() http.Values {
	if u.query == nil {
		u.query = http.Values{}
		u.url.QueryArgs().VisitAll(func(key []byte, value []byte) {
			u.query.Set(string(key), string(value))
		})
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
		stdhdr *http.Header
	}

	ResponseHeader struct {
		header *fasthttp.ResponseHeader
		stdhdr *http.Header
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

//////////////////////////////////////////////////////////////////
