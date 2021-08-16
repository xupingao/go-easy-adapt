package http

import "net/http"

type ResponseWriter interface {
	//http.ResponseWriter
	Write([]byte) (int, error)
	SetHeader(key, value string)
	GetHeader() map[string][]string
	SetStatusCode(statusCode int)
	SetCookie(name, value string, maxAge int, path, domain string, sameSite http.SameSite, secure, httpOnly bool)
	Status() int
	Size() int
	WriteString(string) (int, error)
	Written() bool
	WriteHeaderNow()

	Pusher() http.Pusher
}

