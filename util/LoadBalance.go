package util

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

var LB *LoadBalance

type HttpServer struct {
	Host   string
	Weight int
}

func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host: host, Weight: weight}
}

type LoadBalance struct {
	Servers  []*HttpServer
	CurIndex int // 指向当前服务器索引
}

func init(){
	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://127.0.0.1:9091", 1))
	LB.AddServer(NewHttpServer("http://127.0.0.1:9092", 5))
	LB.AddServer(NewHttpServer("http://127.0.0.1:9093", 5))
}

func NewLoadBalance() *LoadBalance {
	fmt.Println("初始化lb")
	return &LoadBalance{
		Servers: make([]*HttpServer, 0),
		CurIndex: 0,
	}
}

func (lb *LoadBalance) AddServer(server *HttpServer) {
	lb.Servers = append(lb.Servers, server)
}

// 随机算法
func (lb *LoadBalance) SelectByRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(lb.Servers))
	return lb.Servers[index]
}

// IP Hash
func (lb *LoadBalance) SelectByIPHash(ip string) *HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(lb.Servers)
	return lb.Servers[index]
}

// 加权随机
func (lb *LoadBalance) SelectByWeightRand() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	var weightSlice []int
	for i, server := range lb.Servers {
		if server.Weight > 0 {
			for j := 0; j < server.Weight; j++ {
				weightSlice = append(weightSlice, i)
			}
		}
	}
	fmt.Println(weightSlice)
	index := rand.Intn(len(weightSlice))
	return lb.Servers[weightSlice[index]]
}

// 1 2 3
// 1 3 6
// 2 5 1
// 2 7 8
// 加权随机2
func (lb *LoadBalance) SelectByWeightRand2() *HttpServer {
	rand.Seed(time.Now().UnixNano())
	var sum int
	var sumList = make([]int, len(lb.Servers))
	for i, server := range lb.Servers {
		sum += server.Weight
		sumList[i] = sum
	}
	index := rand.Intn(sum)
	for k, v := range sumList {
		if index < v {
			return lb.Servers[k]
		}
	}
	return lb.Servers[0]
}

// 轮询
func (lb *LoadBalance) RoundRobin() *HttpServer {
	server := lb.Servers[lb.CurIndex]
	//lb.CurIndex++
	//if lb.CurIndex >= len(lb.Servers) {
	//	lb.CurIndex = 0
	//}
	lb.CurIndex = (lb.CurIndex + 1) % len(lb.Servers)

	return server
}
