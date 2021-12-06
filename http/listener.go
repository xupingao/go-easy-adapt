package http

import (
	"crypto/tls"
	"net"
	"time"
)

type ListenerCreator func(address string, scheme string) (net.Listener, error)

type TLSListenerCreator func(address string, scheme string,config*tls.Config) (net.Listener, error)

var _ ListenerCreator = NewListener
var _ TLSListenerCreator = NewTLSListener


func NewTLSListener(address string, scheme string,config *tls.Config) (net.Listener, error) {
	ln, err := NewListener(address, scheme)
	if err != nil {
		return nil, err
	}
	return tls.NewListener(ln, config), nil
}

func NewListener(address string, scheme string) (net.Listener, error) {
	l, err := net.Listen(scheme, address)
	if err != nil {
		return nil, err
	}
	if listener, ok := l.(*net.TCPListener);ok {
		return &tcpKeepAliveWrapper{listener}, nil
	}
	return l, nil
}

type tcpKeepAliveWrapper struct {
	*net.TCPListener
}

func (ln tcpKeepAliveWrapper) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return tc, err
	}
	err = tc.SetKeepAlive(true)
	if err != nil {
		return tc, err
	}
	err = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, err
}
