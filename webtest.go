package main

import (
	"fmt"
	"net/http"
	"web"
)

func main() {
	m := web.New()
	m.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	m.ListenServer(":50000")
}
