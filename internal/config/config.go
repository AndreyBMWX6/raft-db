package config

import (
	"log"
	"net"
	"time"
	"math/rand"
	"../message"
)

type Config struct {
	FollowerTimeout  time.Duration
	VotingTimeout    time.Duration
	HeartbeatTimeout time.Duration

	Addr net.UDPAddr
	Neighbors []net.UDPAddr

	Term uint32
	Entries []*message.Entry
}

type RouterConfig struct {
	URLs  []string
}

func NewConfig() *Config {
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:800")
	if err != nil {
		log.Fatal(err)
	}

	neighbourStrings := []string{
		"127.0.0.1:8001",
		"127.0.0.1:8002",
		"127.0.0.1:8003",
		"127.0.0.1:8004",
		"127.0.0.1:8005",
	}

	var neighbors []net.UDPAddr
	for _, neighbour := range neighbourStrings {
		nborAddr, err := net.ResolveUDPAddr("udp4", neighbour)
		if err != nil {
			log.Fatal(err)
		}
		neighbors = append(neighbors, *nborAddr)
	}



	return &Config{
		FollowerTimeout:  4000*time.Millisecond,
		VotingTimeout:    time.Duration(rand.Intn(1000) + 1000)*time.Millisecond,
		HeartbeatTimeout: 1000*time.Millisecond,
		Addr: *addr,
		Neighbors: neighbors,
		Term:      0,
		Entries:   nil,
	}
}

func NewRouterConfig() *RouterConfig {
	urls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
		"http://localhost:8084",
		"http://localhost:8085",
		"http://localhost:8086",
	}

	return &RouterConfig{
		URLs:  urls,
	}
}
