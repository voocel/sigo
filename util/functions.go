package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func CloneHead(src http.Header, dest http.Header) {
	for k, v := range src {
		dest.Set(k, v[0])
	}
}

func RequestUrl(w http.ResponseWriter, r *http.Request, url string) {
	fmt.Println(r.RemoteAddr)
	newreq, _ := http.NewRequest(r.Method, url, r.Body)
	CloneHead(r.Header, newreq.Header)                 // 将用户请求头转发给后端服务器
	newreq.Header.Add("x-forwarded-for", r.RemoteAddr) // 将用户IP转发给后端服务器
	//newres, _ := http.DefaultClient.Do(newreq)  // 向后端服务器发起请求

	dt := http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	newres, _ := dt.RoundTrip(newreq)

	CloneHead(newres.Header, w.Header()) // 将后端服务器响应的头返回给用户
	w.WriteHeader(newres.StatusCode)     // 将后端服务器响应的状态码返回给用户

	defer newres.Body.Close()
	content, _ := ioutil.ReadAll(newres.Body)
	w.Write(content) // 响应给用户
}
