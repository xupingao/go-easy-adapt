package http

import "context"

type Server interface {
	ListenAndServe(address string, engine Handler) error
	ListenAndServeTLS(addr, certFile, keyFile string, engine Handler) error
	// Serve(ln net.Listener, engine Engine) error
	// ServeTLS(ln net.Listener, certFile, keyFile string, engine Engine) error
	DisposedServer() interface{}

	SetHandler(Handler)

	// Stop() error
	Shutdown(ctx context.Context) error
}

type Handler interface {
	ServeHTTP(ctx Context)
}

type Context interface {
	Request() HTTPRequest
	Response() HTTPResponse
	Redirect(int, string)
}

type Header interface {
	// Add adds the key, value pair to the header. It appends to any existing values
	// associated with key.
	Add(string, string)

	// Del deletes the values associated with key.
	Del(string)

	// Get gets the first value associated with the given key. If there are
	// no values associated with the key, Get returns "".
	Get(string) string

	// Set sets the header entries associated with key to the single element value.
	// It replaces any existing values associated with key.
	Set(string, string)

	// value() map[string][]string
}
