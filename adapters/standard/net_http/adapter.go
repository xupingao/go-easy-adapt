package net_http
import (
	"bufio"
	"context"
	"crypto/tls"
	"github.com/xupingao/go-easy-adapt/http"
	"io"
	"mime/multipart"
	"net"
	net_HTTP "net/http"
	"net/url"
)

var _ http.Server = serverAdapter{}

func NewServer() http.Server {
	return serverAdapter{
		Server: net_HTTP.Server{},
	}
}

type serverAdapter struct {
	net_HTTP.Server
}
func (s serverAdapter) ListenAndServeTLS(addr, certFile, keyFile string, engine http.Engine) error {
	s.Addr = addr
	s.Handler = wrapHandler(engine)
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s serverAdapter) Serve(ln net.Listener, engine http.Engine) error {
	s.Handler = wrapHandler(engine)
	return s.Server.Serve(ln)
}

func (s serverAdapter) ServeTLS(ln net.Listener, certFile, keyFile string, engine http.Engine) error {
	s.Handler = wrapHandler(engine)
	return s.Server.ServeTLS(ln, certFile, keyFile)
}

func (s serverAdapter) DisposedServer() interface{} {
	return s.Server
}

func (s serverAdapter) ListenAndServe(address string, engine http.Engine) error {
	s.Addr = address
	s.Handler = wrapHandler(engine)
	return s.Server.ListenAndServe()
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type handler struct {
	engine http.Engine
}
func wrapHandler(engine http.Engine) handler {
	return handler{engine: engine}
}

func (h handler) ServeHTTP(w net_HTTP.ResponseWriter,r *net_HTTP.Request) {
	ctx := wrapConext(w, r)
	h.engine.ServeHTTP(ctx)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

var _ http.Context = httpContext{}

func wrapConext(net_HTTP.ResponseWriter, *net_HTTP.Request) httpContext {
	return httpContext{}
}

type httpContext struct {
	request *net_HTTP.Request
	responseWriter net_HTTP.ResponseWriter
}

func (c httpContext) Request() http.HTTPRequest {
	return &httpRequest{ctx:c}
}

func (c httpContext) Response() http.HTTPResponse {
	return &httpResponseWriter{ctx:c}
}

func (c httpContext) Redirect(code int, url string) {
	net_HTTP.Redirect(c.responseWriter, c.request, url, code )
}

///////////////////////////////////////////////////////////////
type httpRequest struct {
	ctx httpContext
}

func (h httpRequest) Context() context.Context {
	panic("implement me")
}

func (h httpRequest) WithContext(ctx context.Context) http.HTTPRequest {
	panic("implement me")
}

func (h httpRequest) Clone(ctx context.Context) *http.HTTPRequest {
	panic("implement me")
}

func (h httpRequest) Cookies() []*http.Cookie {
	panic("implement me")
}

func (h httpRequest) Cookie(name string) (*http.Cookie, error) {
	panic("implement me")
}

func (h httpRequest) AddCookie(c *http.Cookie) {
	panic("implement me")
}

func (h httpRequest) Header() http.Header {
	panic("implement me")
}

func (h httpRequest) UserAgent() string {
	panic("implement me")
}

func (h httpRequest) Referer() string {
	panic("implement me")
}

func (h httpRequest) Proto() string {
	panic("implement me")
}

func (h httpRequest) Body() io.ReadCloser {
	panic("implement me")
}

func (h httpRequest) Method() string {
	panic("implement me")
}

func (h httpRequest) URL() *url.URL {
	panic("implement me")
}

func (h httpRequest) Query() http.Values {
	panic("implement me")
}

func (h httpRequest) ContentLength() int64 {
	panic("implement me")
}

func (h httpRequest) Host() string {
	panic("implement me")
}

func (h httpRequest) TransferEncoding() []string {
	panic("implement me")
}

func (h httpRequest) Form() http.Values {
	panic("implement me")
}

func (h httpRequest) PostForm() http.Values {
	panic("implement me")
}

func (h httpRequest) Trailer() http.Header {
	panic("implement me")
}

func (h httpRequest) RemoteAddr() string {
	panic("implement me")
}

func (h httpRequest) MultipartReader() (*multipart.Reader, error) {
	panic("implement me")
}

func (h httpRequest) TLS() *tls.ConnectionState {
	panic("implement me")
}

func (h httpRequest) Write(w io.Writer, usingProxy bool, extraHeaders http.Header, waitForContinue func() bool) (err error) {
	panic("implement me")
}

/////////////////////////////////////////////////////////////////
type httpResponseWriter struct {
	ctx httpContext
}

func (h httpResponseWriter) Write([]byte) (int, error) {
	panic("implement me")
}

func (h httpResponseWriter) SetHeader(key, value string) {
	panic("implement me")
}

func (h httpResponseWriter) GetHeader() map[string][]string {
	panic("implement me")
}

func (h httpResponseWriter) WriteHeader(statusCode int) {
	panic("implement me")
}

func (h httpResponseWriter) Hijack() (interface{}, *bufio.ReadWriter, error) {
	panic("implement me")
}

func (h httpResponseWriter) Flush() {
	panic("implement me")
}

func (h httpResponseWriter) Status() int {
	panic("implement me")
}

func (h httpResponseWriter) Size() int {
	panic("implement me")
}

func (h httpResponseWriter) WriteString(string) (int, error) {
	panic("implement me")
}

func (h httpResponseWriter) Written() bool {
	panic("implement me")
}

func (h httpResponseWriter) WriteHeaderNow() {
	panic("implement me")
}

func (h httpResponseWriter) Pusher() net_HTTP.Pusher {
	panic("implement me")
}
