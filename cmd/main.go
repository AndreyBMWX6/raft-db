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

	go rm.ProcessMessage()

	node.RunRolePlayer(candidate)
}
