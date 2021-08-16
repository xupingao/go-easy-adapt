package http

import (
	"bufio"
	"net/http"
)

type HTTPResponse interface {
	Write([]byte) (int, error)
	SetHeader(key, value string)
	GetHeader() map[string][]string
	WriteHeader(statusCode int)
	Hijack() (interface{},*bufio.ReadWriter, error)
	Flush()
	Status() int
	Size() int
	WriteString(string) (int, error)
	Written() bool
	WriteHeaderNow()
	Pusher() http.Pusher
}

