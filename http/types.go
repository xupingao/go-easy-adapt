package http

type Header interface {
	Add(string, string)
	Del(string)
	Get(string) string
	Set(string, string)
	All() map[string][]string
}

type Values interface {
	Get(key string) string
	Gets(key string) []string
	Add(key, value string)
	Set(key, value string)
	Del(key string)
	All() map[string][]string
}

type URL interface {
	SetPath(string)
	RawPath() string
	Path() string
	QueryValue(string) string
	QueryValues(string) []string
	Query() Values
	RawQuery() string
	SetRawQuery(string)
	String() string
	Object() interface{}
}