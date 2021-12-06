package standard

import (
	"crypto/tls"
	"time"

	"github.com/xupingao/go-easy-adapt/http"
)

type Config struct {
	Address      string        // TCP address to listen on.
	Scheme       string        // TCP address to listen on.
	ReadTimeout  time.Duration // Maximum duration before timing out read of the request.
	WriteTimeout time.Duration // Maximum duration before timing out write of the response.

	TLS          bool
	TLSConfig    *tls.Config
	TLSCertFile  string // TLS certificate file path.
	TLSKeyFile   string // TLS key file path.
	DisableHTTP2 bool

	Handler http.Handler

	ListenerCreator    http.ListenerCreator
	TLSListenerCreator http.TLSListenerCreator
}

// func (c *Config) Print(handler string) {
// 	var s string
// 	if c.TLSConfig != nil {
// 		s = `s`
// 	}
// 	log.Printf("%s ⇛ http%s server started on %s\n", handler, s, c.Listener.Addr())
// }

type ConfigSetter func(*Config)

func Address(v string) ConfigSetter {
	return func(c *Config) {
		c.Address = v
	}
}

func Scheme(v string) ConfigSetter {
	return func(c *Config) {
		c.Scheme = v
	}
}

func TLS(v bool) ConfigSetter {
	return func(c *Config) {
		c.TLS = v
	}
}

func DisableHTTP2(v bool) ConfigSetter {
	return func(c *Config) {
		c.DisableHTTP2 = v
	}
}

func Handler(h http.Handler) ConfigSetter {
	return func(c *Config) {
		c.Handler = h
	}
}
func TLSConfig(v *tls.Config) ConfigSetter {
	return func(c *Config) {
		c.TLSConfig = v
	}
}

// TLSCertFile TLS certificate file path.
func TLSCertFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSCertFile = v
	}
}

func TLSKeyFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSKeyFile = v
	}
}

func ReadTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.ReadTimeout = v
	}
}

func WriteTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.WriteTimeout = v
	}
}

func ListenerCreator(v http.ListenerCreator) ConfigSetter {
	return func(c *Config) {
		c.ListenerCreator = v
	}
}

func TLSListenerCreator(v http.TLSListenerCreator) ConfigSetter {
	return func(c *Config) {
		c.TLSListenerCreator = v
	}
}
