package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func (web1handle) getRealIP(r *http.Request) string {
	ips := r.Header.Get("x-forwarded-for")
	if ips != "" {
		ipList := strings.Split(ips, ",")
		if len(ipList) > 0 && ipList[0] != "" {
			return ipList[0]
		}
	}
	return r.RemoteAddr
}

type web1handle struct{}

func (this web1handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start web1...")
	auth := r.Header.Get("Authorization")

	if auth == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="你必须输入用户名和密码"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authList := strings.Split(auth, " ")
	if len(authList) != 2 || authList[0] != "Basic" {
		w.Write([]byte("用户名或密码错误"))
		return
	}
	res, err := base64.StdEncoding.DecodeString(authList[1])
	fmt.Println(res)
	fmt.Println(string(res))
	if err != nil || string(res) != "abc:123" {
		w.Write([]byte("用户名或密码错误"))
		return
	}

	w.Write([]byte(fmt.Sprintf("<h1>web1: %s<h1>", this.getRealIP(r))))
}

type web2handle struct{}

func (web2handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>web2<h1>"))
}

func main() {
	go func() {
		http.ListenAndServe(":9091", web1handle{})
	}()

	go func() {
		http.ListenAndServe(":9092", web2handle{})
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(s)
}
