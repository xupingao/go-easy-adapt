package http

import (
	"context"
	"crypto/tls"
	"io"
	"mime/multipart"
	"time"
)

type SameSite int
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

type Values interface {
	Get(key string) string
	Gets(key string) []string
	Add(key, value string)
	Set(key, value string)
	Del(key string)
}

// URL defines an interface for HTTP request url.
type URL interface {
	SetPath(string)
	RawPath() string
	Path() string
	QueryValue(string) string
	QueryValues(string) []string
	Query() Values
	RawQuery() string
	SetRawQuery(string)
	String() string
	Object() interface{}
}

type HTTPRequest interface {
	// Context() context.Context
	// WithContext(ctx context.Context) HTTPRequest

	Clone(ctx context.Context) *HTTPRequest

	Scheme() string
	Proto() string
	Host() string
	SetHost(string)

	URL() URL
	URI() string
	SetURI(string)

	Method() string
	SetMethod(string)

	Header() Header

	RemoteAddr() string
	// RealIP() string

	Body() io.ReadCloser
	SetBody(io.Reader)

	Form() Values
	FormValue(string) string
	PostForm() Values
	// MultipartForm returns the multipart form.
	MultipartForm() *multipart.Form
	Referer() string
	UserAgent() string
	Size() int64

	Cookies() []*Cookie
	Cookie(name string) (*Cookie, error)
	AddCookie(c *Cookie)

	Query() Values

	TransferEncoding() []string
	Trailer() Header

	MultipartReader() (*multipart.Reader, error)

	IsTLS() bool
	TLS() *tls.ConnectionState

	/////////////////////////////////////////////////
	// Write(w io.Writer, usingProxy bool, extraHeaders Header, waitForContinue func() bool) (err error)
	// FormFile returns the multipart form file for the provided name.
	// FormFile(string) (multipart.File, *multipart.FileHeader, error)

	// BasicAuth() (string, string, bool)
}
