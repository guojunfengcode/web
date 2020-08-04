package web

import (
	"net/http"
	"strings"
)

type HandleFunc func(c *Context)

type Route struct {
	roots    map[string]*Node
	handlers map[string]HandleFunc
}

func NewRoute() *Route {
	return &Route{
		roots:    make(map[string]*Node),
		handlers: make(map[string]HandleFunc),
	}
}

func ParsePattern(pattern string) []string {
	s := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range s {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Route) AddRoute(method string, path string, handler HandleFunc) {
	parts := ParsePattern(path)
	key := method + "-" + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &Node{}
	}
	r.roots[method].Insert(path, parts, 0)

	r.handlers[key] = handler
}

func (r *Route) GetRoute(method string, path string) (*Node, map[string]string) {
	searchParts := ParsePattern(path)
	//log.Println(searchParts)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.Search(searchParts, 0)

	if n != nil {
		parts := ParsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Route) Handle(c *Context) {
	n, param := r.GetRoute(c.Method, c.Path)
	if n != nil {
		c.Params = param
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.StringFmt(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	/*
		key := c.Method + "-" + c.Path
		if handler, ok := r.handlers[key]; ok {
			handler(c)
		} else {
			fmt.Fprintf(c.Write, "404 NOT FOUND: %s\n", c.Path)
		}
	*/
}
