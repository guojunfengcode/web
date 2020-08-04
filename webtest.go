package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"src/db"
	"strings"
	"web"

	"gopkg.in/gcfg.v1"
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

var i int
var j int

func PatternRespon(c *web.Context) {
	i++
	c.StringFmt(http.StatusOK, "hello %s, you're at %s, %v\n", c.Param("name"), c.Path, i)
}

func HtmlRespon(c *web.Context) {
	//c.HtmlFmt(http.StatusOK, "<h1>Hello</h1>")
	c.HtmlFmt(http.StatusOK, `<!DOCTYPE html>
	<html>
		<head>
			<meta charset=\"utf-8\">
			<title>测试web(*.com)</title>
		</head>
	<body>
		<h1>我的第一个标题</h1>
		<p>我的第一个段落。</p>
	</body>
	</html>`)
}

func BasicRespon(c *web.Context) {
	req := c.Req
	w := c.Write
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="User Login"`)
		//w.WriteHeader(http.StatusOK)
		//c.HtmlFmt(http.StatusOK, "<h1>Hello</h1>")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(auth)
	auths := strings.SplitN(auth, " ", 2)
	if len(auths) != 2 {
		fmt.Println("error")
		return
	}
	authMethod := auths[0]
	authB64 := auths[1]
	switch authMethod {
	case "Basic":
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			fmt.Println(err)
			io.WriteString(w, "Unauthorized!\n")
			return
		}
		fmt.Println(string(authstr))
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			fmt.Println("error")
			return
		}
		username := userPwd[0]
		password := userPwd[1]
		fmt.Println("Username:", username)
		fmt.Println("Password:", password)
		notfint := db.Select(username, password)
		if notfint == 0 {
			db.Insert(username, password, authB64)
			//info := fmt.Sprintf(html, `E:\image\shohoku2.jpg`, username)

			//c.HtmlFmt(http.StatusOK, info)
		} else if notfint == 1 {
			db.UpdateLastTime(username)
			//info := fmt.Sprintf("<h1>WELCOME %v LOGIN</h1>", username)
			//info := fmt.Sprintf(html, `E:\image\ROSE.jpg`, username)
			//c.HtmlFmt(http.StatusOK, info)
			t, err := template.ParseFiles("login.html")
			if err != nil {
				fmt.Fprintf(c.Write, "parse template error: %s", err.Error())
				return
			}
			t.Execute(c.Write, nil)
		}
		fmt.Println()
	default:
		fmt.Println("error")
		return
	}
	//io.WriteString(w, "hello, world!\n")
	//c.HtmlFmt(http.StatusOK, "<h1>LOGIN SUCCESS</h1>")
}

func loginSuccess(c *web.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	notfint := db.Select(username, password)
	if notfint == -1 {
		c.HtmlFmt(http.StatusOK, "<h1>LOGIN FAILED, this user does not exist, need REGISTER </h1>")
	} else if notfint == -2 {
		c.HtmlFmt(http.StatusOK, "<h1>LOGIN FAILED, password error</h1>")
	} else if notfint == 1 {
		db.UpdateLastTime(username)
		c.HtmlFmt(http.StatusOK, "<h1>LOGIN SUCCESS</h1>")
	}
}

func register(c *web.Context) {
	t, err := template.ParseFiles("register.html")
	if err != nil {
		fmt.Fprintf(c.Write, "parse template error: %s", err.Error())
		return
	}
	t.Execute(c.Write, nil)
}

func photo(c *web.Context) {
	http.ServeFile(c.Write, c.Req, "test.html")
}

func UserRegister(c *web.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	isexist := db.IsExistUser(username)
	if isexist == 1 {
		c.HtmlFmt(http.StatusBadRequest, "<h1>REGISTER ERROR, USERNAME IS EXIST</h1>")
		return
	}
	basic := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(basic))
	db.Insert(username, password, encoded)
	c.HtmlFmt(http.StatusOK, "<h1>REGISTER SUCCESS</h1>")
}

func main() {
	server := web.New()
	config := struct {
		Mysql struct {
			Username string
			Password string
			Hostip   string
			Port     string
			Database string
		}
	}{}
	err := gcfg.ReadFileInto(&config, "conf.ini")
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}
	db.MysqlPing(config.Mysql.Username, config.Mysql.Password, config.Mysql.Hostip, config.Mysql.Port, config.Mysql.Database)
	defer db.Mysql.Close()

	server.Get("/", func(c *web.Context) {
		cookie := c.Req.Cookies()
		for value, _ := range cookie {
			log.Printf("cookie %v:%v", cookie[value], value)
		}
		t, err := template.ParseFiles("login.html")
		if err != nil {
			fmt.Fprintf(c.Write, "parse template error: %s", err.Error())
			return
		}
		t.Execute(c.Write, nil)
	})
	server.Post("/login", loginSuccess)
	server.Post("/userregister", UserRegister)
	server.Post("/register", register)
	server.Post("/photo", photo)
	//server.Post("/login", JsonRespon)
	server.Post("/xml", XmlRespon)
	server.Get("/hello/:name", PatternRespon)
	server.Get("/html", HtmlRespon)
	server.Get("/basic", BasicRespon)
	/*
		http.HandleFunc("/ROSE.jpg", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	*/
	server.Post("/file", FileRespon)
	server.Get("/test", TestRespon)
	server.ListenServer(":50000")
}

func FileRespon(c *web.Context) {
	http.ServeFile(c.Write, c.Req, "SHOHOKU.jpg")
}

func TestRespon(c *web.Context) {
	c.Req.ParseForm()
	if c.Req.Method == "GET" {
		//http.ServeFile(c.Write, c.Req, `E:\gowork\src\phtot\ROSE.jpg`)
		t, err := template.ParseFiles("login.html")
		if err != nil {
			fmt.Fprintf(c.Write, "parse template error: %s", err.Error())
			return
		}
		t.Execute(c.Write, nil)
	} else {
		username := c.Req.Form["username"]
		password := c.Req.Form["password"]
		fmt.Fprintf(c.Write, "username = %s, password = %s", username, password)
	}

}
