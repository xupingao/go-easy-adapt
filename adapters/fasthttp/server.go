package fasthttp

import (
	"github.com/valyala/fasthttp"
	"github.com/xupingao/go-easy-adapt/http"
	"net"
	"context"
)

var _ http.Server = server{}

type server struct {
	*fasthttp.Server
	config   *Config
}


func NewDefaultServer() http.Server {
	return NewWithConfig(DefaultConfig)
}

func New(addr string, opts ...ConfigSetter) *server {
	c := &Config{Address: addr}
	c.handler = DefaultHandler
	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

type defaultHandler struct{}

var DefaultHandler http.Handler = &defaultHandler{}

func (h *defaultHandler) ServeHTTP(ctx http.Context) {
	ctx.Response().WriteString("Handler of this server is not set")
}

func NewWithTLS(addr, certFile, keyFile string, opts ...ConfigSetter) *server {
	c := &Config{
		Address:     addr,
		TLS:         true,
		TLSCertFile: certFile,
		TLSKeyFile:  keyFile,
		handler:     DefaultHandler,
	}
	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

func NewWithConfig(c *Config) (s *server) {
	s = &server{
		Server: &fasthttp.Server{
			ReadTimeout:        c.ReadTimeout,
			WriteTimeout:       c.WriteTimeout,

			MaxConnsPerIP:      c.MaxConnsPerIP,
			MaxRequestsPerConn: c.MaxRequestsPerConn,
			MaxRequestBodySize: c.MaxRequestBodySize,
		},
		config: c,
	}
	return
}

func (s server) SetHandler(h http.Handler) {
	s.config.handler = h
	s.Server.Handler = wrapHandler(h).ServeHTTP
}

func (s server) Run() error {
	if s.config.TLS {
		//if s.config.TLSConfig != nil {
		//	s.Server.TLSConfig = s.config.TLSConfig
		//}
		return s.Server.ListenAndServeTLS(s.config.Address,s.config.TLSCertFile, s.config.TLSKeyFile)

	} else {
		return s.Server.ListenAndServe(s.config.Address)
	}
}

func (s server) Serve(ln net.Listener) error {
	if s.config.TLS {
		//if s.config.TLSConfig != nil {
		//	s.Server.Tls = s.config.TLSConfig
		//}
		s.Server.ServeTLS(ln, s.config.TLSCertFile, s.config.TLSKeyFile)
	} else {
		s.Server.Serve(ln)
	}
	return s.Server.Serve(ln)
}

func (s server) ListenAndServeTLS(addr, certFile, keyFile string, engine http.Handler) error {
	s.Server.Handler = wrapHandler(engine).ServeHTTP
	return s.Server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (s server) ListenAndServe(addr string, engine http.Handler) error {

	s.Server.Handler = wrapHandler(engine).ServeHTTP
	return s.Server.ListenAndServe(addr)
}

func (s server) DisposedServer() interface{} {
	return s.Server
}

func (s server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown()
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type handler struct {
	engine http.Handler
}

func wrapHandler(engine http.Handler) handler {
	return handler{engine: engine}
}

func (h handler) ServeHTTP(c *fasthttp.RequestCtx) {
	ctx := wrapConext(c)
	h.engine.ServeHTTP(ctx)
}