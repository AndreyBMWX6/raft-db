package tests

import (
	"net"
	"testing"
	"log"

	"../message"
	"../node"
	"../config"
)

func TestApplyRaftMessage(t *testing.T) {
	// candidate ApplyRaftMessages testing
	// message with lower term
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:800")
	if err != nil {
		log.Fatal(err)
	}

	neighborStrings := []string{
		"127.0.0.1:8001",
		"127.0.0.1:8002",
		"127.0.0.1:8003",
		"127.0.0.1:8004",
		"127.0.0.1:8005",
	}

	var neighbors []net.Addr
	for _, neighbor := range neighborStrings {
		nborAddr, err := net.ResolveUDPAddr("udp4", neighbor)
		if err != nil {
			log.Fatal(err)
		}
		neighbors = append(neighbors, nborAddr)
	}

	var raftIn  = make(chan message.RaftMessage)
	var raftOut = make(chan message.RaftMessage)

	var clientIn  = make(chan message.ClientMessage)
	var clientOut = make(chan message.ClientMessage)

	var raftNode = &node.RaftCore{
		Config:    config.NewConfig(),
		Addr:      addr,
		Neighbors: neighbors,
		RaftIn:    raftIn,
		RaftOut:   raftOut,
		ClientIn:  clientIn,
		ClientOut: clientOut,
		Term:      1,
		Entries:   nil,
	}

	var candidate = node.NewCandidate(raftNode)


	owner, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8001")
	if err != nil {
		log.Fatal(err)
	}

	lowermessage := message.NewRequestVote(
		&message.BaseRaftMessage{
			Owner: owner,
			Dest: addr,
			CurrTerm: 0, // lower than candidate term(1)
		},
	)
	// empty EntriesLog
	lowermessage.TopTerm = 0
	lowermessage.TopIndex = 0

	result := candidate.ApplyRaftMessage(lowermessage)

	if result != nil {
		t.Error(
			"For msg.Term", lowermessage.CurrTerm,
			"expected", nil,
			"got", result)
	}
}