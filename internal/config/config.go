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

	// Raft IO
	RaftIn  chan message.RaftMessage
	RaftOut chan message.RaftMessage

	// Client IO
	ClientIn  <-chan message.ClientMessage
	ClientOut chan<- message.ClientMessage
}

func NewConfig() *Config {
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:800")
	if err != nil {
		log.Fatal(err)
	}

	neighbourStrings := []string{
		"127.0.0.1:8001",
		//"127.0.0.1:8002",
		//"127.0.0.1:8003",
		//"127.0.0.1:8004",
		//"127.0.0.1:8005",
	}

	var neighbors []net.UDPAddr
	for _, neighbour := range neighbourStrings {
		nborAddr, err := net.ResolveUDPAddr("udp4", neighbour)
		if err != nil {
			log.Fatal(err)
		}
		neighbors = append(neighbors, *nborAddr)
	}

	var raftIn  = make(chan message.RaftMessage)
	var raftOut = make(chan message.RaftMessage)

	var clientIn  = make(chan message.ClientMessage)
	var clientOut = make(chan message.ClientMessage)

	return &Config{
		FollowerTimeout:  1000*time.Millisecond,
		VotingTimeout:    time.Duration(rand.Intn(1000) + 100)*time.Millisecond,
		HeartbeatTimeout: 1000*time.Millisecond,
		Addr: *addr,
		Neighbors: neighbors,
		Term:      1,
		Entries:   nil,
		RaftIn:    raftIn,
		RaftOut:   raftOut,
		ClientIn:  clientIn,
		ClientOut: clientOut,
	}
}
