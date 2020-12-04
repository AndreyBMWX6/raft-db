package main

import (
	"../internal/manager"
	"../internal/node"
)

func main() {
	var raftNode = node.NewRaftCore()

	var candidate = node.NewCandidate(raftNode)


	rm := &manager.RaftManager{
		RaftIn: raftNode.Config.RaftOut,
		RaftOut: raftNode.Config.RaftIn,
	}

	cm := &manager.ClientManager{
		ClientIn: raftNode.Config.ClientOut,
		ClientOut: raftNode.Config.ClientIn,
	}

	go rm.ProcessMessage()
	go cm.ProcessEntries()

	node.RunRolePlayer(candidate)
}
