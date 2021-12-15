package adapters

import (
	"bytes"
	"fmt"
	fast "github.com/valyala/fasthttp"
	fast_adapter "github.com/xupingao/go-easy-adapt/adapters/fasthttp"
	std_adapter "github.com/xupingao/go-easy-adapt/adapters/standard"
	"github.com/xupingao/go-easy-adapt/http"
	"github.com/xupingao/go-easy-adapt/mock"
	"io/ioutil"
	"net"
	std "net/http"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var defaultClientsCount = runtime.NumCPU()

func BenchmarkRequestCtxRedirect(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var ctx fast.RequestCtx
		for pb.Next() {
			ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
			ctx.Redirect("/foo/bar?baz=111", http.StatusFound)
		}
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkFastServerGet1ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerGet(b, defaultClientsCount, 1)
}

func BenchmarkFastServerGet2ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerGet(b, defaultClientsCount, 2)
}

func BenchmarkFastServerGet10ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerGet(b, defaultClientsCount, 10)
}

func BenchmarkFastServerGet10KReqPerConn(b *testing.B) {
	benchmarkFastHttpServerGet(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkNetHTTPServerGet1ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 1)
}

func BenchmarkNetHTTPServerGet2ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 2)
}

func BenchmarkNetHTTPServerGet10ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 10)
}

func BenchmarkNetHTTPServerGet10KReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerGet(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkStdAdapterServerGet1ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerGet(b, defaultClientsCount, 1)
}

func BenchmarkStdAdapterServerGet2ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerGet(b, defaultClientsCount, 2)
}

func BenchmarkStdAdapterServerGet10ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerGet(b, defaultClientsCount, 10)
}

func BenchmarkStdAdapterServerGet10KReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerGet(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkFastAdapterServerGet1ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerGet(b, defaultClientsCount, 1)
}

func BenchmarkFastAdapterServerGet2ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerGet(b, defaultClientsCount, 2)
}

func BenchmarkFastAdapterServerGet10ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerGet(b, defaultClientsCount, 10)
}

func BenchmarkFastAdapterServerGet10KReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerGet(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkFastServerPost1ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerPost(b, defaultClientsCount, 1)
}

func BenchmarkFastServerPost2ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerPost(b, defaultClientsCount, 2)
}

func BenchmarkFastServerPost10ReqPerConn(b *testing.B) {
	benchmarkFastHttpServerPost(b, defaultClientsCount, 10)
}

func BenchmarkFastServerPost10KReqPerConn(b *testing.B) {
	benchmarkFastHttpServerPost(b, defaultClientsCount, 10000)
}
//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkNetHTTPServerPost1ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerPost(b, defaultClientsCount, 1)
}

func BenchmarkNetHTTPServerPost2ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerPost(b, defaultClientsCount, 2)
}

func BenchmarkNetHTTPServerPost10ReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerPost(b, defaultClientsCount, 10)
}

func BenchmarkNetHTTPServerPost10KReqPerConn(b *testing.B) {
	benchmarkNetHTTPServerPost(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkStdAdapterServerPost1ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerPost(b, defaultClientsCount, 1)
}

func BenchmarkStdAdapterServerPost2ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerPost(b, defaultClientsCount, 2)
}

func BenchmarkStdAdapterServerPost10ReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerPost(b, defaultClientsCount, 10)
}

func BenchmarkStdAdapterServerPost10KReqPerConn(b *testing.B) {
	benchmarkStdAdapterServerPost(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkFastAdapterServerPost1ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerPost(b, defaultClientsCount, 1)
}

func BenchmarkFastAdapterServerPost2ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerPost(b, defaultClientsCount, 2)
}

func BenchmarkFastAdapterServerPost10ReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerPost(b, defaultClientsCount, 10)
}

func BenchmarkFastAdapterServerPost10KReqPerConn(b *testing.B) {
	benchmarkFastAdapterServerPost(b, defaultClientsCount, 10000)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkFastServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkFastHttpServerGet(b, 10000, 1)
}

func BenchmarkFastServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkFastHttpServerGet(b, 10000, 2)
}

func BenchmarkFastServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkFastHttpServerGet(b, 10000, 10)
}

func BenchmarkFastServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkFastHttpServerGet(b, 10000, 100)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkNetHTTPServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 1)
}

func BenchmarkNetHTTPServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 2)
}

func BenchmarkNetHTTPServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 10)
}

func BenchmarkNetHTTPServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkNetHTTPServerGet(b, 10000, 100)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkFastAdapterServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkFastAdapterServerGet(b, 10000, 1)
}

func BenchmarkFastAdapterServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkFastAdapterServerGet(b, 10000, 2)
}

func BenchmarkFastAdapterServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkFastAdapterServerGet(b, 10000, 10)
}

func BenchmarkFastAdapterServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkFastAdapterServerGet(b, 10000, 100)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func BenchmarkStdAdapterServerGet1ReqPerConn10KClients(b *testing.B) {
	benchmarkStdAdapterServerGet(b, 10000, 1)
}

func BenchmarkStdAdapterServerGet2ReqPerConn10KClients(b *testing.B) {
	benchmarkStdAdapterServerGet(b, 10000, 2)
}

func BenchmarkStdAdapterServerGet10ReqPerConn10KClients(b *testing.B) {
	benchmarkStdAdapterServerGet(b, 10000, 10)
}

func BenchmarkStdAdapterServerGet100ReqPerConn10KClients(b *testing.B) {
	benchmarkStdAdapterServerGet(b, 10000, 100)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkServerHijack(b *testing.B) {
	clientsCount := 1000
	requestsPerConn := 10000
	ch := make(chan struct{}, b.N)
	responseBody := []byte("123")
	s := &fast.Server{
		Handler: func(ctx *fast.RequestCtx) {
			ctx.Hijack(func(c net.Conn) {
				// emulate server loop :)
				err := fast.ServeConn(c, func(ctx *fast.RequestCtx) {
					ctx.Success("foobar", responseBody)
					registerServedRequest(b, ch)
				})
				if err != nil {
					b.Fatalf("error when serving connection")
				}
			})
			ctx.Success("foobar", responseBody)
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	req := "GET /foo HTTP/1.1\r\nHost: google.com\r\n\r\n"
	benchmarkServer(b, s, clientsCount, requestsPerConn, req)
	verifyRequestsServed(b, ch)
}

func BenchmarkServerMaxConnsPerIP(b *testing.B) {
	clientsCount := 1000
	requestsPerConn := 10
	ch := make(chan struct{}, b.N)
	responseBody := []byte("123")
	s := &fast.Server{
		Handler: func(ctx *fast.RequestCtx) {
			ctx.Success("foobar", responseBody)
			registerServedRequest(b, ch)
		},
		MaxConnsPerIP: clientsCount * 2,
		Concurrency:   16 * clientsCount,
	}
	req := "GET /foo HTTP/1.1\r\nHost: google.com\r\n\r\n"
	benchmarkServer(b, s, clientsCount, requestsPerConn, req)
	verifyRequestsServed(b, ch)
}

func BenchmarkServerTimeoutError(b *testing.B) {
	clientsCount := 10
	requestsPerConn := 1
	ch := make(chan struct{}, b.N)
	n := uint32(0)
	responseBody := []byte("123")
	s := &fast.Server{
		Handler: func(ctx *fast.RequestCtx) {
			if atomic.AddUint32(&n, 1)&7 == 0 {
				ctx.TimeoutError("xxx")
				go func() {
					ctx.Success("foobar", responseBody)
				}()
			} else {
				ctx.Success("foobar", responseBody)
			}
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	req := "GET /foo HTTP/1.1\r\nHost: google.com\r\n\r\n"
	benchmarkServer(b, s, clientsCount, requestsPerConn, req)
	verifyRequestsServed(b, ch)
}
//////////////////////////////////////////////////////////////////////////////////////////////////
var (
	fakeResponse = []byte("Hello, world!")
	getRequest   = "GET /foobar?baz HTTP/1.1\r\nHost: google.com\r\nUser-Agent: aaa/bbb/ccc/ddd/eee Firefox Chrome MSIE Opera\r\n" +
		"Referer: http://xxx.com/aaa?bbb=ccc\r\nCookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n"
	postRequest = fmt.Sprintf("POST /foobar?baz HTTP/1.1\r\nHost: google.com\r\nContent-Type: foo/bar\r\nContent-Length: %d\r\n"+
		"User-Agent: Opera Chrome MSIE Firefox and other/1.2.34\r\nReferer: http://google.com/aaaa/bbb/ccc\r\n"+
		"Cookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n%s",
		len(fakeResponse), fakeResponse)
)

func benchmarkFastHttpServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := &fast.Server{
		Handler: func(ctx *fast.RequestCtx) {
			if !ctx.IsGet() {
				b.Fatalf("Unexpected request method: %s", ctx.Method())
			}
			ctx.Success("text/plain", fakeResponse)
			if requestsPerConn == 1 {
				ctx.SetConnectionClose()
			}
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}

func benchmarkNetHTTPServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := &std.Server{
		Handler: std.HandlerFunc(func(w std.ResponseWriter, req *std.Request) {
			if req.Method != http.MethodGet {
				b.Fatalf("Unexpected request method: %s", req.Method)
			}
			h := w.Header()
			h.Set("Content-Type", "text/plain")
			if requestsPerConn == 1 {
				h.Set(http.HeaderConnection, "close")
			}
			w.Write(fakeResponse) //nolint:errcheck
			registerServedRequest(b, ch)
		}),
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}


func benchmarkFastAdapterServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := fast_adapter.NewDefaultServer()
	s.SetHandler(http.HandlerFunc(func(ctx http.Context) {
		if ctx.Request().Method() != http.MethodGet {
			b.Fatalf("Unexpected request method: %s", ctx.Request().Method())
		}
		h := ctx.Request().Header()
		h.Set("Content-Type", "text/plain")
		if requestsPerConn == 1 {
			h.Set(http.HeaderConnection, "close")
		}
		ctx.Response().Write(fakeResponse) //nolint:errcheck
		registerServedRequest(b, ch)
	}))
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}


func benchmarkStdAdapterServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := std_adapter.NewDefaultServer()
	s.SetHandler(http.HandlerFunc(func(ctx http.Context) {
		if ctx.Request().Method() != http.MethodGet {
			b.Fatalf("Unexpected request method: %s", ctx.Request().Method())
		}
		h := ctx.Request().Header()
		h.Set("Content-Type", "text/plain")
		if requestsPerConn == 1 {
			h.Set(http.HeaderConnection, "close")
		}
		ctx.Response().Write(fakeResponse) //nolint:errcheck
		registerServedRequest(b, ch)
	}))
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func benchmarkFastHttpServerPost(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := &fast.Server{
		Handler: func(ctx *fast.RequestCtx) {
			if !ctx.IsPost() {
				b.Fatalf("Unexpected request method: %s", ctx.Method())
			}
			body := ctx.Request.Body()
			if !bytes.Equal(body, fakeResponse) {
				b.Fatalf("Unexpected body %q. Expected %q", body, fakeResponse)
			}
			ctx.Success("text/plain", body)
			if requestsPerConn == 1 {
				ctx.SetConnectionClose()
			}
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, postRequest)
	verifyRequestsServed(b, ch)
}

func benchmarkNetHTTPServerPost(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := &std.Server{
		Handler: std.HandlerFunc(func(w std.ResponseWriter, req *std.Request) {
			if req.Method != http.MethodPost {
				b.Fatalf("Unexpected request method: %s", req.Method)
			}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				b.Fatalf("Unexpected error: %s", err)
			}
			req.Body.Close()
			if !bytes.Equal(body, fakeResponse) {
				b.Fatalf("Unexpected body %q. Expected %q", body, fakeResponse)
			}
			h := w.Header()
			h.Set("Content-Type", "text/plain")
			if requestsPerConn == 1 {
				h.Set(http.HeaderConnection, "close")
			}
			w.Write(body) //nolint:errcheck
			registerServedRequest(b, ch)
		}),
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, postRequest)
	verifyRequestsServed(b, ch)
}

func benchmarkFastAdapterServerPost(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := fast_adapter.NewDefaultServer()
	s.SetHandler(http.HandlerFunc(func(ctx http.Context){
		if ctx.Request().Method() != http.MethodPost {
			b.Fatalf("Unexpected request method: %s", ctx.Request().Method())
		}
		body, err := ioutil.ReadAll(ctx.Request().Body())
		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
		ctx.Request().Body().Close()
		if !bytes.Equal(body, fakeResponse) {
			b.Fatalf("Unexpected body %q. Expected %q", body, fakeResponse)
		}
		h := ctx.Response().Header()
		h.Set("Content-Type", "text/plain")
		if requestsPerConn == 1 {
			h.Set(http.HeaderConnection, "close")
		}
		ctx.Response().Write(body) //nolint:errcheck
		registerServedRequest(b, ch)
	}))

	benchmarkServer(b, s, clientsCount, requestsPerConn, postRequest)
	verifyRequestsServed(b, ch)
}

func benchmarkStdAdapterServerPost(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := std_adapter.NewDefaultServer()
	s.SetHandler(http.HandlerFunc(func(ctx http.Context){
		if ctx.Request().Method() != http.MethodPost {
			b.Fatalf("Unexpected request method: %s", ctx.Request().Method())
		}
		body, err := ioutil.ReadAll(ctx.Request().Body())
		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
		ctx.Request().Body().Close()
		if !bytes.Equal(body, fakeResponse) {
			b.Fatalf("Unexpected body %q. Expected %q", body, fakeResponse)
		}
		h := ctx.Response().Header()
		h.Set("Content-Type", "text/plain")
		if requestsPerConn == 1 {
			h.Set(http.HeaderConnection, "close")
		}
		ctx.Response().Write(body) //nolint:errcheck
		registerServedRequest(b, ch)
	}))

	benchmarkServer(b, s, clientsCount, requestsPerConn, postRequest)
	verifyRequestsServed(b, ch)
}


//////////////////////////////////////////////////////////////////////////////////////////////////

func registerServedRequest(b *testing.B, ch chan<- struct{}) {
	select {
	case ch <- struct{}{}:
	default:
		b.Fatalf("More than %d requests served", cap(ch))
	}
}

func verifyRequestsServed(b *testing.B, ch <-chan struct{}) {
	requestsServed := 0
	for len(ch) > 0 {
		<-ch
		requestsServed++
	}
	requestsSent := b.N
	for requestsServed < requestsSent {
		select {
		case <-ch:
			requestsServed++
		case <-time.After(100 * time.Millisecond):
			b.Fatalf("Unexpected number of requests served %d. Expected %d", requestsServed, requestsSent)
		}
	}
}

type testServer interface {
	Serve(ln net.Listener) error
}

func benchmarkServer(b *testing.B, s testServer, clientsCount, requestsPerConn int, request string) {
	ln := mock.NewFakeListener(b.N, clientsCount, requestsPerConn, request)
	ch := make(chan struct{})
	go func() {
		s.Serve(ln)
		//err := s.Serve(ln)
		//if err != nil {
		//	b.Fatalf("Server.Serve() reports: %s", err.Error())
		//}
		ch <- struct{}{}
	}()

	<-ln.Done

	select {
	case <-ch:
	case <-time.After(10 * time.Second):
		b.Fatalf("Server.Serve() didn't stop")
	}
}
