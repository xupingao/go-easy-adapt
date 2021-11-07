package fasthttp

import (
	"time"

	"github.com/xupingao/go-easy-adapt/http"
)

type Config struct {
	Address      string        // TCP address to listen on.
	ReadTimeout  time.Duration // Maximum duration before timing out read of the request.
	WriteTimeout time.Duration // Maximum duration before timing out write of the response.

	TLS         bool
	//TLSConfig   *tls.Config
	TLSCertFile string // TLS certificate file path.
	TLSKeyFile  string // TLS key file path.

	handler http.Handler

	MaxConnsPerIP      int
	MaxRequestsPerConn int
	MaxRequestBodySize int
}

var DefaultConfig *Config = &Config{
	Address: ":8080",
	TLS:     false,

	handler: DefaultHandler,
}

// func (c *Config) Print(engine string) {
// 	var s string
// 	if c.TLSConfig != nil {
// 		s = `s`
// 	}
// 	log.Printf("%s ⇛ http%s server started on %s\n", engine, s, c.Listener.Addr())
// }

type ConfigSetter func(*Config)

// Address TCP address to listen on.
func Address(v string) ConfigSetter {
	return func(c *Config) {
		c.Address = v
	}
}

func TLS(v bool) ConfigSetter {
	return func(c *Config) {
		c.TLS = v
	}
}

func Handler(h http.Handler) ConfigSetter {
	return func(c *Config) {
		c.handler = h
	}
}
//func TLSConfig(v *tls.Config) ConfigSetter {
//	return func(c *Config) {
//		c.TLSConfig = v
//	}
//}

// TLSCertFile TLS certificate file path.
func TLSCertFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSCertFile = v
	}
}

// TLSKeyFile TLS key file path.
func TLSKeyFile(v string) ConfigSetter {
	return func(c *Config) {
		c.TLSKeyFile = v
	}
}

// ReadTimeout Maximum duration before timing out read of the request.
func ReadTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.ReadTimeout = v
	}
}

// WriteTimeout Maximum duration before timing out write of the response.
func WriteTimeout(v time.Duration) ConfigSetter {
	return func(c *Config) {
		c.WriteTimeout = v
	}
}


func MaxConnsPerIP(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxConnsPerIP = v
	}
}

func MaxRequestsPerConn(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxRequestsPerConn = v
	}
}

func MaxRequestBodySize(v int) ConfigSetter {
	return func(c *Config) {
		c.MaxRequestBodySize = v
	}
}
