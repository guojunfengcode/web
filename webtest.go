package main

import (
	"net/http"
	"web"
)

func XmlRespon(c *web.Context) {
	v := &web.Servers{Version: "1"}
	v.Svs = append(v.Svs, web.Server{c.PostForm("username"), c.PostForm("password")})
	v.Svs = append(v.Svs, web.Server{"Beijing_VPN", "127.0.0.2"})
	c.XmlFmt(http.StatusOK, v)
}

func JsonRespon(c *web.Context) {
	c.JsonFmt(http.StatusOK, web.Json{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
}

func PatternRespon(c *web.Context) {
	c.StringFmt(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
}

func main() {
	server := web.New()

	server.Get("/", func(c *web.Context) {
		c.StringFmt(http.StatusOK, "hello %s", "world")
	})

	server.Post("/login", JsonRespon)
	server.Post("/xml", XmlRespon)
	server.Get("/hello/:name", PatternRespon)
	server.ListenServer(":50000")
}
