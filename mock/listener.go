package mock

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var _ net.Listener = (*FakeListener)(nil)
type FakeListener struct {
	Lock            sync.Mutex
	RequestsCount   int
	requestsPerConn int
	Request         []byte
	Ch              chan *fakeServerConn
	Done            chan struct{}
	Closed          bool
}

func (ln *FakeListener) Accept() (net.Conn, error) {
	ln.Lock.Lock()
	if ln.RequestsCount == 0 {
		ln.Lock.Unlock()
		for len(ln.Ch) < cap(ln.Ch) {
			time.Sleep(10 * time.Millisecond)
		}
		ln.Lock.Lock()
		if !ln.Closed {
			close(ln.Done)
			ln.Closed = true
		}
		ln.Lock.Unlock()
		return nil, io.EOF
	}
	requestsCount := ln.requestsPerConn
	if requestsCount > ln.RequestsCount {
		requestsCount = ln.RequestsCount
	}
	ln.RequestsCount -= requestsCount
	ln.Lock.Unlock()

	c := <-ln.Ch
	c.requestsCount = requestsCount
	c.closed = 0
	c.pos = 0

	return c, nil
}

func (ln *FakeListener) Close() error {
	return nil
}

func (ln *FakeListener) Addr() net.Addr {
	return &fakeAddr
}

func NewFakeListener(requestsCount, clientsCount, requestsPerConn int, request string) *FakeListener {
	ln := &FakeListener{
		RequestsCount:   requestsCount,
		requestsPerConn: requestsPerConn,
		Request:         []byte(request),
		Ch:              make(chan *fakeServerConn, clientsCount),
		Done:            make(chan struct{}),
	}
	for i := 0; i < clientsCount; i++ {
		ln.Ch <- &fakeServerConn{
			ln: ln,
		}
	}
	return ln
}

var _ net.Conn = (*fakeServerConn)(nil)

type fakeServerConn struct {
	net.TCPConn
	ln            *FakeListener
	requestsCount int
	pos           int
	closed        uint32
}

func (c *fakeServerConn) Read(b []byte) (int, error) {
	nn := 0
	reqLen := len(c.ln.Request)
	for len(b) > 0 {
		if c.requestsCount == 0 {
			if nn == 0 {
				return 0, io.EOF
			}
			return nn, nil
		}
		pos := c.pos % reqLen
		n := copy(b, c.ln.Request[pos:])
		b = b[n:]
		nn += n
		c.pos += n
		if n+pos == reqLen {
			c.requestsCount--
		}
	}
	return nn, nil
}

func (c *fakeServerConn) Write(b []byte) (int, error) {
	return len(b), nil
}

var fakeAddr = net.TCPAddr{
	IP:   []byte{1, 2, 3, 4},
	Port: 12345,
}

func (c *fakeServerConn) RemoteAddr() net.Addr {
	return &fakeAddr
}

func (c *fakeServerConn) Close() error {
	if atomic.AddUint32(&c.closed, 1) == 1 {
		c.ln.Ch <- c
	}
	return nil
}

func (c *fakeServerConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *fakeServerConn) SetWriteDeadline(t time.Time) error {
	return nil
}
