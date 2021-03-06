package standard

import (
	"crypto/tls"
	"github.com/webx-top/echo/engine"
	"github.com/xupingao/go-easy-adapt/http"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	net_HTTP "net/http"
	"sync"
)

var _ http.Context = (*std_context)(nil)

func NewContext(writer net_HTTP.ResponseWriter, req *net_HTTP.Request) http.Context {
	return wrapStdConext(writer, req)
}

type ContextPool struct {
	context  sync.Pool
	request  sync.Pool
	response sync.Pool
	url      sync.Pool
	//requestHeader  sync.Pool
	//responseHeader sync.Pool
}

func NewContextPool() *ContextPool {
	return &ContextPool{
		context: sync.Pool{
			New: func() interface{} {
				return &std_context{}
			},
		},
		request: sync.Pool{
			New: func() interface{} {
				return &httpRequest{}
			},
		},
		response: sync.Pool{
			New: func() interface{} {
				return &httpResponse{}
			},
		},
		//requestHeader: sync.Pool{
		//	New: func() interface{} {
		//		return &Header{}
		//	},
		//},
		//responseHeader: sync.Pool{
		//	New: func() interface{} {
		//		return &Header{}
		//	},
		//},
		url: sync.Pool{
			New: func() interface{} {
				return &URL{}
			},
		},
	}
}

func (p *ContextPool) AllocateContext(writer net_HTTP.ResponseWriter, request *net_HTTP.Request) http.Context {
	// Request
	ctx := p.context.Get().(*std_context)
	req := p.request.Get().(*httpRequest)
	ctx.rawResponseWriter = writer
	ctx.rawRequest = request
	req.Request = request

	url := p.url.Get().(*URL)
	url.url = request.URL
	req.url = url
	//reqHdr := s.pool.requestHeader.Get().(*Header)
	//reqHdr.reset(r.Header)
	//reqURL := s.pool.url.Get().(*URL)
	//reqURL.reset(r.URL)
	//req.reset(r, reqHdr, reqURL)
	//req.config = s.config

	// Response
	res := p.response.Get().(*httpResponse)
	res.ResponseWriter = writer
	res.writer = writer

	ctx.response = res
	ctx.request = req
	return ctx
	////resHdr := s.pool.responseHeader.Get().(*Header)
	////resHdr.reset(w.Header())
	////res.reset(w, r, resHdr)
	////res.config = s.config
	//
	//s.handler.ServeHTTP(req, res)
	//
	//s.pool.request.Put(req)
	//s.pool.requestHeader.Put(reqHdr)
	//s.pool.url.Put(reqURL)
	//s.pool.response.Put(res)
	//s.pool.responseHeader.Put(resHdr)
}

func (p *ContextPool) ReleaseContext(ctx http.Context) {
	context, ok := ctx.(*std_context)
	if !ok {
		return
	}
	p.url.Put(context.request.url)
	context.request.reset()
	context.response.reset()
	p.request.Put(context.request)
	p.response.Put(context.response)
	context.reset()
	p.context.Put(context)
}

func wrapStdConext(writer net_HTTP.ResponseWriter, req *net_HTTP.Request) http.Context {
	return std_context{
		rawRequest:        req,
		rawResponseWriter: writer,
		request:           &httpRequest{Request: req},
		response:          &httpResponse{ResponseWriter: writer, writer: writer},
	}
}

type std_context struct {
	rawRequest        *net_HTTP.Request
	rawResponseWriter net_HTTP.ResponseWriter
	request           *httpRequest
	response          *httpResponse
}

func (c *std_context) reset() {
	c.rawRequest = nil
	c.rawResponseWriter = nil
	c.request = nil
	c.response = nil
}

func (c std_context) Request() http.HTTPRequest {
	return c.request
}

func (c std_context) Response() http.HTTPResponse {
	return c.response
}

func (c std_context) Redirect(code int, url string) {
	net_HTTP.Redirect(c.rawResponseWriter, c.rawRequest, url, code)
}

//----------------
// Param
//----------------

var defaultMaxRequestBodySize int64 = 32 << 20 // 32 MB

var _ http.HTTPRequest = (*httpRequest)(nil)

type httpRequest struct {
	*net_HTTP.Request

	url      http.URL
	query    http.Values
	postForm http.Values
	form     http.Values
	header   http.Header
	trailer  http.Header
}

func (r *httpRequest) reset() {
	r.url = nil
	r.query = nil
	r.postForm = nil
	r.form = nil
	r.header = nil
	r.trailer = nil
}

func (r *httpRequest) Scheme() string {
	if r.IsTLS() {
		return http.SchemeHTTPS
	}

	if len(r.Request.URL.Scheme) > 0 {
		return r.Request.URL.Scheme
	}
	return http.SchemeHTTP
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
	//if r.url == nil {
	//	r.url = &URL{url: r.Request.URL}
	//}
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
	if r.header == nil {
		r.header = &header{r.Request.Header}
	}
	return r.header
}

func (r *httpRequest) RemoteAddr() string {
	return r.Request.RemoteAddr
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
	if r.form == nil {
		r.MultipartForm()
		r.form = &values{&r.Request.Form}
	}

	return r.form
}

func (r *httpRequest) FormValue(name string) string {
	r.Form()
	return r.form.Get(name)
}

func (r *httpRequest) PostForm() http.Values {
	if r.postForm == nil {
		r.postForm = &values{rawValues: &r.Request.PostForm}
	}
	return r.postForm
}

func (r *httpRequest) MultipartForm() (*multipart.Form, error) {
	if r.Request.MultipartForm == nil {
		err := r.Request.ParseMultipartForm(defaultMaxRequestBodySize)
		if err != nil {
			return nil, err
		}
	}
	return r.Request.MultipartForm, nil
}

func (r *httpRequest) Referer() string {
	return r.Request.Referer()
}

func (r *httpRequest) UserAgent() string {
	return r.Request.UserAgent()
}

func (r *httpRequest) Size() int64 {
	return r.Request.ContentLength
}

//func (r *httpRequest) Cookies() []*http.Cookie {
//	panic("")
//}

func (r *httpRequest) Cookie(name string) string {
	if cookie, err := r.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ``
}

//func (r *httpRequest) AddCookie(c *http.Cookie) {
//	panic("")
//}

func (r *httpRequest) Query() http.Values {
	if r.query == nil {
		if r.url == nil {
			r.url = &URL{url: r.Request.URL}
		}
		r.query = r.url.Query()
	}
	return r.query
}

func (r *httpRequest) TransferEncoding() []string {
	return r.Request.TransferEncoding
}
func (r *httpRequest) Trailer() http.Header {
	if r.trailer == nil {
		r.trailer = &header{r.Request.Trailer}
	}
	return r.trailer
}

func (r *httpRequest) MultipartReader() (*multipart.Reader, error) {
	return r.Request.MultipartReader()
}

func (r *httpRequest) IsTLS() bool {
	return r.Request.TLS != nil
}

func (r *httpRequest) TLS() *tls.ConnectionState {
	return r.Request.TLS
}

//----------------
// Param
//----------------
const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type httpResponse struct {
	net_HTTP.ResponseWriter
	//logger logger.Logger
	writen bool
	status int
	size   int64
	writer io.Writer

	header http.Header
}

func (r *httpResponse) reset() {
	r.writen = false
	r.status = 0
	r.size = 0
	r.writer = nil
	r.header = nil
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

func (r *httpResponse) SetCookie(cookie *net_HTTP.Cookie) {
	r.Header().Add(engine.HeaderSetCookie, cookie.String())
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
	if r.header == nil {
		r.header = &header{r.ResponseWriter.Header()}
	}
	return r.header
}

func (r *httpResponse) HeaderWritten() bool {
	return r.writen
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

func (r *httpResponse) SetWriter(w io.Writer) {
	r.writer = w
}

func (r *httpResponse) Writer() io.Writer {
	return r.writer
}
