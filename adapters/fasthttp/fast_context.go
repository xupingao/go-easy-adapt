package fasthttp

import (
	"bytes"
	"crypto/tls"
	"github.com/valyala/fasthttp"
	"github.com/webx-top/echo/logger"
	"github.com/xupingao/go-easy-adapt/http"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	net_HTTP "net/http"
	"net/url"
	"strings"
	"unsafe"
)

var _ http.Context = (*fast_context)(nil)

func wrapFastConext(requestCtx *fasthttp.RequestCtx) http.Context {
	return fast_context{
		rawRequestCtx: requestCtx,
		request:       &httpRequest{RequestCtx: requestCtx},
		response:      &httpResponse{RequestCtx: requestCtx, writer:requestCtx},
	}
}

type fast_context struct {
	rawRequestCtx *fasthttp.RequestCtx

	request  http.HTTPRequest
	response http.HTTPResponse
}

func (c fast_context) Request() http.HTTPRequest {
	return c.request
}

func (c fast_context) Response() http.HTTPResponse {
	return c.response
}

func (c fast_context) Redirect(code int, url string) {
	c.rawRequestCtx.Redirect(url, code)
}

///////////////////////////////////////////////////////////////
var _ http.HTTPRequest = (*httpRequest)(nil)

type httpRequest struct {
	*fasthttp.RequestCtx

	url    http.URL
	header http.Header

	query    http.Values
	postForm http.Values
	form     http.Values
}

func (r *httpRequest) initValues() {
	//if r.query != nil {
	//	return
	//}
	//r.query = http.Values{}
	//r.RequestCtx.QueryArgs().VisitAll(func(k, v []byte) {
	//	r.query.Add(Bytes2str(k), Bytes2str(v))
	//})
	//
	//r.postForm = http.Values{}
	//r.RequestCtx.PostArgs().VisitAll(func(k, v []byte) {
	//	r.postForm.Add(Bytes2str(k), Bytes2str(v))
	//})
	//if r.form != nil {
	//	return
	//}
	//r.form = http.Values{}
	//for key, vals := range r.query {
	//	r.form[key] = vals
	//}
	//for key, vals := range r.postForm {
	//	r.form[key] = vals
	//}
	//mf := r.MultipartForm()
	//if mf != nil && mf.Value != nil {
	//	for key, vals := range mf.Value {
	//		r.form[key] = vals
	//	}
	//}
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
	return r.RequestCtx.RemoteAddr().String()
}

func (r *httpRequest) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(r.RequestCtx.PostBody()))

}
func (r *httpRequest) SetBody(reader io.Reader) {
	r.RequestCtx.Request.SetBodyStream(reader, 0)
}

func (r *httpRequest) Form() http.Values {
	if r.form != nil {
		return r.form
	}
	form := url.Values{}

	r.RequestCtx.QueryArgs().VisitAll(func(k, v []byte) {
		form.Add(Bytes2str(k), Bytes2str(v))
	})

	r.RequestCtx.PostArgs().VisitAll(func(k, v []byte) {
		form.Add(Bytes2str(k), Bytes2str(v))
	})
	mf, err := r.MultipartForm()
	if err == nil && mf != nil && mf.Value != nil {
		for key, vals := range mf.Value {
			form[key] = vals
		}
	}
	r.form = &formValues{rawValues: &form}
	return r.form
}
func (r *httpRequest) FormValue(key string) string {
	return r.Form().Get(key)
}
func (r *httpRequest) PostForm() http.Values {
	if r.postForm == nil {
		r.postForm = &argValues{args: r.PostArgs()}
	}
	r.initValues()
	return r.postForm
}

func (r *httpRequest) MultipartForm() (*multipart.Form, error) {
	if !strings.HasPrefix(string(r.RequestCtx.Request.Header.ContentType()), http.MIMEMultipartForm) {
		return nil, nil
	}
	return r.RequestCtx.MultipartForm()
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

func (r *httpRequest) Cookie(key string) string {
	return Bytes2str(r.RequestCtx.Request.Header.Cookie(key))
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
	writer io.Writer
	header http.Header
}

func (r *httpResponse) WriteHeader(statusCode int) {
	if statusCode > 0 && r.status != statusCode {
		//if r.Written() {
		//	r.logger.Warn("response already committed")
		//}
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

func (r *httpResponse) SetCookie(cookie *net_HTTP.Cookie) {
	r.header.Add(http.HeaderSetCookie, cookie.String())
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
	return r.writen
}

func (r *httpResponse) SetWriter(w io.Writer) {
	r.writer = w
}

func (r *httpResponse) Writer() io.Writer {
	return r.writer
}
