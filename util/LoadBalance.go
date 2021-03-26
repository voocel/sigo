package util

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

type HttpServer struct {
	Host   string
	Weight int
}

func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{Host: host, Weight: weight}
}

type LoadBalance struct {
	Servers []*HttpServer
}

func NewLoadBalance() *LoadBalance {
	return &LoadBalance{
		Servers: make([]*HttpServer, 0),
	}
}

func (lb *LoadBalance) AddServer(server *HttpServer) {
	lb.Servers = append(lb.Servers, server)
}

func (lb *LoadBalance) SelectByRand() *HttpServer { // 随机算法
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(lb.Servers))
	return lb.Servers[index]
}

func (lb *LoadBalance) SelectByIPHash(ip string) *HttpServer { // IP Hash
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(lb.Servers)
	return lb.Servers[index]
}

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
