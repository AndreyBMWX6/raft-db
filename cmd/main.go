package main

import (
	"flag"
	"log"

	"github.com/AndreyBMWX6/raft-db/internal/manager"
	"github.com/AndreyBMWX6/raft-db/internal/node"
)

func main() {
	flag.Parse()
	ip := flag.Arg(0)
	if ip == "" {
		log.Fatalf("Error: no ip")
	}
	ipPort := flag.Arg(1)
	if ipPort == "" {
		log.Fatalf("Error: no port")
	}
	urlPort := flag.Arg(2)
	if urlPort == "" {
		log.Fatalf("Error: no URL port")
	}
	username := flag.Arg(3)
	if username == "" {
		log.Fatalf("Error: no username")
	}

	var raftNode = node.NewRaftCore(ip, ipPort, urlPort)

	var candidate = node.NewCandidate(raftNode)

	rm := &manager.RaftManager{
		RaftIn:  raftNode.RaftOut,
		RaftOut: raftNode.RaftIn,
		Addr:    raftNode.Addr,
	}

	cm := &manager.ClientManager{
		ClientIn:  raftNode.ClientOut,
		ClientOut: raftNode.ClientIn,
		UrlPort:   urlPort,
	}

	dbm := &manager.MongoDBManager{
		DBIn:       raftNode.DBOut,
		DBOut:      raftNode.DBIn,
		Username:   username,
	}

	go rm.ProcessMessage()
	go cm.ProcessEntries()
	go dbm.ProcessMessage()

	node.RunRolePlayer(candidate)
}
