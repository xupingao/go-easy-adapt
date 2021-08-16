package http

import "net"

type Server interface {
	ListenAndServe(address string ,engine Engine) error
	ListenAndServeTLS(addr, certFile, keyFile string, engine Engine) error
	Serve(ln net.Listener, engine Engine) error
	ServeTLS(ln net.Listener, certFile, keyFile string, engine Engine) error
	DisposedServer()interface{}
}

type Engine interface {
	ServeHTTP(ctx Context)
}

type Context interface {
	Request() HTTPRequest
	Response() HTTPResponse
	Redirect(int, string)
}


