package web

import (
	"fmt"
	"log"
)

type HandleFunc func(c *Context)

type Route struct {
	handlers map[string]HandleFunc
}

func NewRoute() *Route {
	return &Route{handlers: make(map[string]HandleFunc)}
}

func (r *Route) AddRoute(method string, path string, handler HandleFunc) {
	log.Printf("Route %s - %s", method, path)
	key := method + "-" + path
	r.handlers[key] = handler
}

func (r *Route) Handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		fmt.Fprintf(c.Write, "404 NOT FOUND: %s\n", c.Path)
	}
}
