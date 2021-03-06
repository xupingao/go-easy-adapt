package http

import (
	"io"
	"net"
	"net/http"
)

type HTTPResponse interface {
	Write([]byte) (int, error)

	Status() int
	Size() int64
	Hijacker(fn func(net.Conn)) error

	Header() Header

	WriteHeader(statusCode int)
	HeaderWritten() bool

	//Flush()

	WriteString(string) (int, error)

	WriteHeaderNow()

	SetCookie(*http.Cookie)
	// KeepBody(bool)
	SetWriter(io.Writer)
	Writer() io.Writer
	// Body() []byte
	// Redirect(string, int)
	// NotFound()
	// ServeFile(string)
	// Stream(func(io.Writer) bool)
	// Error(string, ...int)
}

type HTTP2Response interface {
	Pusher() http.Pusher
}
