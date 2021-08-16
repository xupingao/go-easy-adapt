package http

type Server interface {
	ListenAndServe(address string ,engine Engine) error
}

type Engine interface {
	ServeHTTP(ctx Context)
}

type Context interface {
	RequestReader() RequestReader
	ResponseWriter() ResponseWriter
	Redirect(int, string)
}


