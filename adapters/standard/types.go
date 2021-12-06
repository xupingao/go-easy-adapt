package standard

import (
	"github.com/xupingao/go-easy-adapt/http"
	"net/url"
)

type values struct {
	rawValues *url.Values
}

func (v *values) Add(key string, value string) {
	v.rawValues.Add(key, value)
}

func (v *values) Del(key string) {
	v.rawValues.Del(key)
}

func (v *values) Get(key string) string {
	return v.rawValues.Get(key)
}

func (v *values) Gets(key string) []string {
	form := *v.rawValues
	if v, ok := form[key]; ok {
		return v
	}
	return []string{}
}

func (v *values) Set(key string, value string) {
	v.rawValues.Set(key, value)
}

type URL struct {
	url   *url.URL
	query http.Values
}

func (u *URL) SetPath(path string) {
	u.url.Path = path
}

func (u *URL) RawPath() string {
	return u.url.EscapedPath()
}

func (u *URL) Path() string {
	return u.url.Path
}

func (u *URL) QueryValue(name string) string {
	u.initQuery()
	return u.query.Get(name)
}

func (u *URL) QueryValues(name string) []string {
	u.initQuery()
	return u.query.Gets(name)
}

func (u *URL) Query() http.Values {
	u.initQuery()
	return u.query
}

func (u *URL) initQuery() {
	if u.query != nil {
		return
	}
	query := u.url.Query()
	u.query = &values{rawValues:&query}
}

func (u *URL) reset(url *url.URL) {
	u.url = url
	u.query = nil
}

func (u *URL) RawQuery() string {
	return u.url.RawQuery
}

func (u *URL) SetRawQuery(rawQuery string) {
	u.url.RawQuery = rawQuery
}

func (u *URL) String() string {
	return u.url.String()
}

func (u *URL) Object() interface{} {
	return u.url
}


