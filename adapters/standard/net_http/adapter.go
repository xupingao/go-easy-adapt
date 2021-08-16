package net_http
import  "github.com/xupingao/go-easy-adapt/http"
import HTTP "net/http"
var _ http.Server = Server{}
type Server struct {

}
func (s Server) ListenAndServe(address string, engine http.Engine) error {
	return HTTP.ListenAndServe(address, NewHandler(engine))
}
/////////////////////////////////////////////////////////////////////////////////////////////////

type Handler  struct {
	engine http.Engine
}
func NewHandler(engine http.Engine) Handler{
	return Handler{engine:engine}
}

func (h Handler) ServeHTTP(w HTTP.ResponseWriter,r *HTTP.Request) {
	ctx := wrapConext(w, r)
	h.engine.ServeHTTP(ctx)
}

/////////////////////////////////////////////////////////////////////////////////////////////////

var _ http.Context = Context{}

func wrapConext(HTTP.ResponseWriter, *HTTP.Request) Context {
	return Context{}
}

type Context struct {

}

func (c Context) RequestReader() http.RequestReader {
	panic("implement me")
}

func (c Context) ResponseWriter() http.ResponseWriter {
	panic("implement me")
}

func (c Context) Redirect(int, string) {
	panic("implement me")
}



