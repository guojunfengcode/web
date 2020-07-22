package web

import (
	"fmt"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Manage struct {
	conn   string
	handle Handler
}

func New() *Manage {
	m := new(Manage)
	return m
}

func (m *Manage) AddReq(method string, path string, handle Handler) {
	m.conn = method + "-" + path
	m.handle = handle
}

func (m *Manage) Get(path string, handle Handler) {
	m.AddReq("GET", path, handle)
}

func (m *Manage) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn := req.Method + "-" + req.URL.Path
	if m.conn == conn {
		m.handle(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func (m *Manage) ListenServer(addr string) error {
	err := http.ListenAndServe(addr, m)
	return err
}
