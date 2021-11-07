package standard

import (
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	net_HTTP "net/http"
	"net/url"

	"github.com/webx-top/echo/logger"
	"github.com/xupingao/go-easy-adapt/http"
)

var _ http.Context = context_{}

func wrapConext(writer net_HTTP.ResponseWriter, req *net_HTTP.Request) http.Context {
	return context_{
		rawRequest:        req,
		rawResponseWriter: writer,
		request:           &httpRequest{Request: req},
		response:          &httpResponse{ResponseWriter: writer},
	}
}

type context_ struct {
	rawRequest        *net_HTTP.Request
	rawResponseWriter net_HTTP.ResponseWriter

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
	net_HTTP.Redirect(c.rawResponseWriter, c.rawRequest, url, code)
}

///////////////////////////////////////////////////////////////
type httpRequest struct {
	*net_HTTP.Request
	// context Context
	url *URL
}

// func (r *httpRequest) Context() context.Context {
// 	panic("")
// }
// func (r *httpRequest) WithContext(ctx context.Context) http.HTTPRequest {
// 	panic("")
// }

func (r *httpRequest) Clone(ctx context.Context) *http.HTTPRequest {
	panic("")
}

func (r *httpRequest) Scheme() string {
	// if r.IsTLS() {
	// 	return echo.SchemeHTTPS
	// }
	return r.Request.URL.Scheme
	// if len(r.request.URL.Scheme) > 0 {
	// 	return r.request.URL.Scheme
	// }
	// return echo.SchemeHTTP
}
func (r *httpRequest) Proto() string {
	return r.Request.Proto
}
func (r *httpRequest) Host() string {
	return r.Request.Host
}
func (r *httpRequest) SetHost(h string) {
	r.Request.Host = h
}
func (r *httpRequest) URL() http.URL {
	if r.url == nil {
		r.url = &URL{url: r.Request.URL}
	}
	return r.url
}
func (r *httpRequest) URI() string {
	return r.Request.RequestURI
}
func (r *httpRequest) SetURI(uri string) {
	r.Request.RequestURI = uri
}

func (r *httpRequest) Method() string {
	return r.Request.Method
}
func (r *httpRequest) SetMethod(method string) {
	r.Request.Method = method
}

func (r *httpRequest) Header() http.Header {
	panic("")
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
	return r.Request.Body
}
func (r *httpRequest) SetBody(reader io.Reader) {
	if readCloser, ok := reader.(io.ReadCloser); ok {
		r.Request.Body = readCloser
	} else {
		r.Request.Body = ioutil.NopCloser(reader)
	}
}

func (r *httpRequest) Form() http.Values {
	panic("")
}
func (r *httpRequest) FormValue(string) string {
	panic("")
}
func (r *httpRequest) PostForm() http.Values {
	panic("")
}

// MultipartForm returns the multipart form.
func (r *httpRequest) MultipartForm() *multipart.Form {
	panic("")
}
func (r *httpRequest) Referer() string {
	panic("")
}
func (r *httpRequest) UserAgent() string {
	panic("")
}
func (r *httpRequest) Size() int64 {
	panic("")
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
	panic("")
}
func (r *httpRequest) TLS() *tls.ConnectionState {
	panic("")
}

/////////////////////////////////////////////////////////////////
const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type httpResponse struct {
	net_HTTP.ResponseWriter
	logger logger.Logger
	writen bool
	status int
	size   int64
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
		r.ResponseWriter.WriteHeader(r.status)
	}
}

func (r *httpResponse) Written() bool {
	return r.writen
}

func (r *httpResponse) Write(b []byte) (n int, err error) {
	// if !r.committed {
	// 	if r.status == 0 {
	// 		r.status = http.StatusOK
	// 	}
	// 	r.WriteHeader(r.status)
	// }
	// if r.keepBody {
	// 	r.body = append(r.body, b...)
	// }
	r.WriteHeaderNow()
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	return n, err
}

func (r *httpResponse) WriteString(s string) (n int, err error) {
	r.WriteHeaderNow()
	n, err = io.WriteString(r.ResponseWriter, s)
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
	conn, bufrw, err := r.ResponseWriter.(net_HTTP.Hijacker).Hijack()
	if err != nil {
		return err
	}
	_ = bufrw
	fn(conn)
	conn.Close()
	r.writen = true
	return nil
}

func (r *httpResponse) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *httpResponse) HeaderWritten() bool {
	panic("")
}

func (r *httpResponse) Flush() {
	r.WriteHeaderNow()
	r.ResponseWriter.(net_HTTP.Flusher).Flush()
}

func (r *httpResponse) Pusher() net_HTTP.Pusher {
	if pusher, ok := r.ResponseWriter.(net_HTTP.Pusher); ok {
		return pusher
	}
	return nil
}

type URL struct {
	url   *url.URL
	query http.Values
}

func (u *URL) SetPath(path string) {
	u.url.Path = path
}

func (u *URL) RawPath() string {
	return u.url.EscapedPath()
}

func (u *URL) Path() string {
	return u.url.Path
}

func (u *URL) QueryValue(name string) string {
	if u.query == nil {
		u.query = http.Values(u.url.Query())
	}
	return u.query.Get(name)
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
		u.query = http.Values(u.url.Query())
	}
	return u.query
}

func (u *URL) reset(url *url.URL) {
	u.url = url
	u.query = nil
}

func (u *URL) RawQuery() string {
	return u.url.RawQuery
}

func (u *URL) SetRawQuery(rawQuery string) {
	u.url.RawQuery = rawQuery
}

func (u *URL) String() string {
	return u.url.String()
}

func (u *URL) Object() interface{} {
	return u.url
}
