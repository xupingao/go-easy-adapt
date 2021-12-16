package standard

import (
	"context"
	"crypto/tls"
	"github.com/xupingao/go-easy-adapt/http"
	"net"
	net_HTTP "net/http"
)

var defalutAddr = ":8080"
var defalutScheme = "tcp"

var _ http.Server = (*server)(nil)

type server struct {
	*net_HTTP.Server
	config      *Config
	listener    net.Listener
	contextPool *ContextPool
}

func (s *server) SetListener(listener net.Listener) {
	s.listener = listener
}

func (s *server) SetHandler(h http.Handler) {
	s.config.Handler = h
}

func NewDefaultServer() http.Server {
	return NewWithConfig(defaultConfig())
}

func defaultConfig() *Config {
	return &Config{
		Address:            defalutAddr,
		Scheme:             defalutScheme,
		TLS:                false,
		Handler:            DefaultHandler(),
		ListenerCreator:    http.NewListener,
		TLSListenerCreator: http.NewTLSListener,
	}
}

func New(addr string, opts ...ConfigSetter) *server {
	c := defaultConfig()
	c.Address = addr
	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

func NewWithTLS(addr, certFile, keyFile string, opts ...ConfigSetter) *server {
	c := defaultConfig()

	c.Address = addr
	c.TLS = true
	c.TLSCertFile = certFile
	c.TLSKeyFile = keyFile

	for _, opt := range opts {
		opt(c)
	}
	return NewWithConfig(c)
}

func NewWithConfig(c *Config) (s *server) {
	s = &server{
		Server:      &net_HTTP.Server{},
		config:      c,
		contextPool: NewContextPool(),
	}
	return
}

func (s *server) Run() error {
	s.applyConfig()
	if err := s.initListener(); err != nil {
		return err
	}

	return s.Server.Serve(s.listener)
}

func (s server) Serve(listener net.Listener) error {
	s.SetListener(listener)
	return s.Run()
}

func (s *server) ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	s.config.Handler = handler
	s.config.Address = addr
	s.config.TLS = true
	s.config.TLSCertFile = certFile
	s.config.TLSKeyFile = keyFile
	return s.Run()
}

func (s *server) ListenAndServe(addr string, handler http.Handler) error {
	s.config.Handler = handler
	s.config.Address = addr
	s.config.TLS = false
	return s.Run()
}

func (s *server) initTlSConfig() {
	if s.config.TLSConfig != nil {
		return
	}
	s.config.TLSConfig = new(tls.Config)
	if len(s.config.TLSCertFile) > 0 && len(s.config.TLSKeyFile) > 0 {
		cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
		if err != nil {
			panic(err)
		}
		s.config.TLSConfig.Certificates = append(s.config.TLSConfig.Certificates, cert)
	}
	if !s.config.DisableHTTP2 {
		s.config.TLSConfig.NextProtos = append(s.config.TLSConfig.NextProtos, "h2")
	}
}

func (s *server) initListener() error {
	if s.listener != nil {
		return nil
	}

	var ln net.Listener
	var err error

	if s.config.TLS {
		s.initTlSConfig()
		if s.config.TLSListenerCreator != nil {
			ln, err = s.config.TLSListenerCreator(s.config.Address, s.config.Scheme, s.config.TLSConfig)
		} else {
			ln, err = http.NewTLSListener(s.config.Address, s.config.Scheme, s.config.TLSConfig)
		}
		if err == nil {
			s.listener = ln
		}
		return err
	}
	if s.config.ListenerCreator != nil {
		ln, err = s.config.ListenerCreator(s.config.Address, s.config.Scheme)
	} else {
		ln, err = http.NewListener(s.config.Address, s.config.Scheme)
	}
	if err == nil {
		s.listener = ln
	}
	return err
}

func (s *server) ApplyConfig() {
	s.applyConfig()
}
func (s *server) applyConfig() {
	if s.config.Handler == nil {
		s.config.Handler = DefaultHandler()
	}
	if s.config.TLSConfig != nil {
		s.Server.TLSConfig = s.config.TLSConfig
	}
	s.Server.Handler = wrapHandler(s.config.Handler, s)
	s.Server.ReadTimeout = s.config.ReadTimeout
	s.Server.WriteTimeout = s.config.WriteTimeout
	s.Server.Addr = s.config.Address
}

func (s server) DisposedServer() interface{} {
	return s.Server
}

func (s *server) Stop() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func DefaultHandler() http.Handler {
	return http.HandlerFunc(func(ctx http.Context) {
		ctx.Response().WriteString("Handler of this server is not set")
	})
}

type handlerWrapper struct {
	handler http.Handler
	s       *server
}

func wrapHandler(handler http.Handler, s *server) *handlerWrapper {
	return &handlerWrapper{handler: handler, s: s}
}

func (wp *handlerWrapper) ServeHTTP(w net_HTTP.ResponseWriter, r *net_HTTP.Request) {
	ctx := wp.s.contextPool.AllocateContext(w, r)
	wp.handler.ServeHTTP(ctx)
	wp.s.contextPool.ReleaseContext(ctx)
}
