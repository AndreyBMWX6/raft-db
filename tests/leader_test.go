package tests

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/AndreyBMWX6/raft-db/internal/message"
	"github.com/AndreyBMWX6/raft-db/internal/node"
)

func TestLeaderApplyAppendEntries(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var leader = node.BecomeLeader(node.NewCandidate(raftNode))
	var sync chan interface{}

	// message with equal term test
	eqMsg := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		0,
		make([]*message.Entry, 0),
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	// will return nil as candidate votes ony for candidate of bigger term
	result := leader.ApplyRaftMessage(eqMsg)

	Synchronize(sync, raftNode)
	require.Equal(t, nil, result, "AppendEntries with eq term test")

	highMsg := eqMsg
	highMsg.CurrTerm = 1 // higher than node term(0)

	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	// will return *node.Follower as next role
	// you can see role change in logs
	// also will update node term
	result = leader.ApplyRaftMessage(highMsg)

	var resultIsFollower bool

	switch result.(type) {
	case *node.Follower:
		resultIsFollower = true
	default:
		resultIsFollower = false
	}

	Synchronize(sync, raftNode)
	require.Equal(t, true, resultIsFollower, "message with higher term test")
	require.Equal(t, highMsg.Term(), result.ReleaseNode().Term, "message with higher term test")
}

// Processing of RequestVote has been already tested in node_test.go TestProcessRequestVote func
// We test Processing of RequestVote individually for candidate cause his behaviour differs
