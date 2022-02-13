package main

import (
	"fmt"
	"gate/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Endpoint struct {
	Host string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type ProxyHandle struct{}

func (ProxyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
	}()

	//if r.URL.Path == "/a" {
	//	util.RequestUrl(w, r, "http://127.0.0.1:9091")
	//	return
	//}
	//if r.URL.Path == "/b" {
	//	util.RequestUrl(w, r, "http://127.0.0.1:9092")
	//	return
	//}

	//for k,v := range util.ProxyConfig{
	//	if matched, _ := regexp.MatchString(k, r.URL.Path); matched == true {
	//		//util.RequestUrl(w, r, v)  // 开始反向代理
	//		target, _ := url.Parse(v)
	//		proxy := httputil.NewSingleHostReverseProxy(target)
	//		proxy.ServeHTTP(w, r)
	//
	//		return
	//	}
	//}
	//
	//w.Write([]byte("default index"))

	//myurl, _ := url.Parse(lb.SelectByRand().Host)
	//myurl, _ := url.Parse(lb.SelectByIPHash(r.RemoteAddr).Host)
	//myurl, _ := url.Parse(lb.SelectByWeightRand2().Host)
	myurl, _ := url.Parse(util.LB.RoundRobin().Host)
	proxy := httputil.NewSingleHostReverseProxy(myurl)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.ListenAndServe(":8080", ProxyHandle{})
}
