package standard

import (
	"context"
	"net"
	net_HTTP "net/http"

	"github.com/xupingao/go-easy-adapt/http"
)

var _ http.Server = &server{}

type server struct {
	*net_HTTP.Server
	config   *Config

}

func (s *server) SetListener(listener *net.Listener) {
	s.config.listener = listener
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
		Server: &net_HTTP.Server{
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			Addr:         c.Address,
			Handler:      wrapHandler(c.handler),
		},
		config: c,
	}
	return
}

func (s server) SetHandler(h http.Handler) {
	s.config.handler = h
	s.Server.Handler = wrapHandler(h)
}

func (s server) Run() error {
	if s.config.TLS {
		if s.config.TLSConfig != nil {
			s.Server.TLSConfig = s.config.TLSConfig
		}
		return s.Server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)

	} else {
		return s.Server.ListenAndServe()
	}
}

func (s server) Serve(ln net.Listener) error {
	if s.config.TLS {
		if s.config.TLSConfig != nil {
			s.Server.TLSConfig = s.config.TLSConfig
		}
		s.Server.ServeTLS(ln, s.config.TLSCertFile, s.config.TLSKeyFile)
	} else {
		s.Server.Serve(ln)
	}
	return s.Server.Serve(ln)
}

func (s server) ListenAndServeTLS(addr, certFile, keyFile string, engine http.Handler) error {
	s.Server.Addr = addr
	s.Server.Handler = wrapHandler(engine)
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s server) ListenAndServe(addr string, engine http.Handler) error {
	s.Server.Addr = addr
	s.Server.Handler = wrapHandler(engine)
	return s.Server.ListenAndServe()
}

func (s server) DisposedServer() interface{} {
	return s.Server
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type handler struct {
	engine http.Handler
}

func wrapHandler(engine http.Handler) handler {
	return handler{engine: engine}
}

func (h handler) ServeHTTP(w net_HTTP.ResponseWriter, r *net_HTTP.Request) {
	ctx := wrapConext(w, r)
	h.engine.ServeHTTP(ctx)
}
