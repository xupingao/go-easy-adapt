package http

import (
	"context"
	"net"
)

type Server interface {
	ListenAndServe(address string, engine Handler) error
	ListenAndServeTLS(addr, certFile, keyFile string, engine Handler) error
	DisposedServer() interface{}
	Run() error
	SetHandler(Handler)
	Serve(listener net.Listener) error
	SetListener(listener net.Listener)
	Shutdown(ctx context.Context) error
}

type Handler interface {
	ServeHTTP(ctx Context)
}

type HandlerFunc func(ctx Context)

func (f HandlerFunc) ServeHTTP(ctx Context) {
	f(ctx)
}

type Context interface {
	Request() HTTPRequest
	Response() HTTPResponse
	Redirect(int, string)
}

