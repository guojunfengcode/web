package web

import (
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Manage struct {
	route *Route
}

func New() *Manage {
	m := &Manage{
		route: NewRoute(),
	}
	return m
}

func (m *Manage) AddReq(method string, path string, handle HandleFunc) {
	m.route.AddRoute(method, path, handle)
}

func (m *Manage) Get(path string, handle HandleFunc) {
	m.AddReq("GET", path, handle)
}

func (m *Manage) Post(path string, handle HandleFunc) {
	m.AddReq("POST", path, handle)
}

func (m *Manage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	m.route.Handle(context)
}

func (m *Manage) ListenServer(addr string) error {
	err := http.ListenAndServe(addr, m)
	return err
}
