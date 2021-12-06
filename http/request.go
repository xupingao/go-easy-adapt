package http

import (
	"crypto/tls"
	"io"
	"mime/multipart"
)

type HTTPRequest interface {
	// Context() context.Context
	// WithContext(ctx context.Context) HTTPRequest

	//Clone(ctx context.Context) *HTTPRequest

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
	MultipartForm() (*multipart.Form, error)
	Referer() string
	UserAgent() string
	Size() int64

	//Cookies() []*Cookie
	//Cookie(name string) (*Cookie, error)
	//AddCookie(c *Cookie)
	Cookie(name string) string


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
