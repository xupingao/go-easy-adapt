package http

import (
	"context"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/url"
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

type Header map[string][]string
type Values map[string][]string

type HTTPRequest interface{
	Context() context.Context
	WithContext(ctx context.Context) HTTPRequest

	Clone(ctx context.Context) *HTTPRequest

	/////////////////////////////////////////////////
	Cookies() []*Cookie
	Cookie(name string) (*Cookie, error)
	AddCookie(c *Cookie)
	/////////////////////////////////////////////////
	Header() Header
	UserAgent() string
	Referer() string
	/////////////////////////////////////////////////
	Proto() string
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() Values
	ContentLength()int64
	Host()string
	TransferEncoding() []string

	Form() Values
	PostForm() Values
	Trailer() Header
	RemoteAddr() string
	MultipartReader() (*multipart.Reader, error)
	TLS() *tls.ConnectionState

	/////////////////////////////////////////////////
	Write(w io.Writer, usingProxy bool, extraHeaders Header, waitForContinue func() bool) (err error)
}
