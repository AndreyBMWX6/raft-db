package main

import (
	"../internal/manager"
	"../internal/node"
	"../internal/router"
)

func main() {
	var raftNode = node.NewRaftCore()

	var candidate = node.NewCandidate(raftNode)

	rm := &manager.RaftManager{
		RaftIn:  raftNode.RaftOut,
		RaftOut: raftNode.RaftIn,
	}

	cm := &manager.ClientManager{
		ClientIn:  raftNode.ClientOut,
		ClientOut: raftNode.ClientIn,
	}

	dbm := &manager.MongoDBManager{
		DBIn:       raftNode.DBOut,
		DBOut:      raftNode.DBIn,
	}

	router := router.NewRouter()

	go rm.ProcessMessage()
	go cm.ProcessEntries()
	go dbm.ProcessMessage()
	go router.RunRouter()

	node.RunRolePlayer(candidate)
}
