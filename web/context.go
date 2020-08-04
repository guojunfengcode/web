package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Json map[string]interface{}

type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []Server `xml:"info"`
}

type Server struct {
	UserName string `xml:"userName"`
	PassWord string `xml:"password"`
}

type Context struct {
	Write  http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
	Params map[string]string
	Status int
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	value := c.Req.FormValue(key)
	return value
}

func (c *Context) Query(key string) string {
	value := c.Req.URL.Query().Get(key)
	return value
}

func (c *Context) Statu(code int) {
	c.Status = code
	c.Write.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Write.Header().Set(key, value)
}

func (c *Context) JsonFmt(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Statu(code)
	encoder := json.NewEncoder(c.Write)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Write, err.Error(), 500)
	}
}

func (c *Context) XmlFmt(code int, v *Servers) {
	c.SetHeader("Content-Type", "text/xml")
	c.Statu(code)
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	c.Write.Write([]byte(xml.Header))
	c.Write.Write(output)

}

func (c *Context) StringFmt(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Statu(code)
	c.Write.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) HtmlFmt(code int, html string) {
	uid := "33333333"
	cookie := &http.Cookie{
		Name:     "testuid",
		Value:    uid,
		Path:     "/html",
		HttpOnly: false,
		MaxAge:   100,
	}
	c.SetHeader("Content-Type", "text/html")
	//http.SetCookie(c.Write, cookie)
	c.SetCook(cookie)
	c.Statu(code)
	c.Write.Write([]byte(html))
}

func (c *Context) SetCook(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		c.Write.Header().Add("Set-Cookie", v)
		log.Printf("Set-Cookie: %v", v)
	}
}
