package tests

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/AndreyBMWX6/raft-db/internal/message"
	"github.com/AndreyBMWX6/raft-db/internal/node"
)

func TestFollowerApplyLowerTermMessage(t *testing.T) {
	// each node in any role ignores any raft messages
	// as ApplyRaftMessage gets message.RaftMessage as parameter
	// and calls Term() method, that return msg.Term for any type of raft message

	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	raftNode.Term = 1
	var follower = node.BecomeFollower(node.NewCandidate(raftNode))

	// AppendEntries with lower term test
	lowAppEnt := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0, // lower than candidate term(1)
		},
		0,
		0,
		make([]*message.Entry, 0),
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil cause message term is lower than candidate term
	result := follower.ApplyRaftMessage(lowAppEnt)

	require.Equal(t, nil, result, "AppendEntries with lower term test")

	// RequestVote with lower term test
	lowReqVot := message.NewRequestVote(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0, // lower than candidate term(1)
		},
		0,
		0,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil cause message term is lower than candidate term
	result = follower.ApplyRaftMessage(lowReqVot)

	require.Equal(t, nil, result, "RequestVote with lower term test")

	// RequestAck with lower term test
	lowAppAck := message.NewAppendAck(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0, // lower than candidate term(1)
		},
		false,
		false,
		0,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil cause message term is lower than candidate term
	result = follower.ApplyRaftMessage(lowAppAck)

	require.Equal(t, nil, result, "AppendAck with lower term test")

	// Candidate don't process RequestAck
}

func TestFollowerApplyAppendEntries(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var follower = node.BecomeFollower(node.NewCandidate(raftNode))
	var res message.AppendAck

	// in AppendEntries case follower refreshes timer
	// and add Entries to log in case of successful request

	// heartbeat message test
	heartbeatMsg := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		0,
		nil,
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, &res)

	follower.ApplyAppendEntries(heartbeatMsg)

	time.Sleep(1)
	require.Equal(t, true, res.Heartbeat, "heartbeat message test")
	require.Equal(t, true, res.Appended, "heartbeat message test")

	// entries to make AppendEntries instead of heartbeat
	entries := make([]*message.Entry, 0)
	entries = append(entries, &message.Entry{
		Term:  0,
		Query: nil,
	})

	// bad NewIndex message test
	badNewIdxMsg := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		12, // more than follower's log length(0)
		entries,
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, &res)

	follower.ApplyAppendEntries(badNewIdxMsg)

	time.Sleep(1)
	require.Equal(t, false, res.Heartbeat, "bad NewIndex message test")
	require.Equal(t, false, res.Appended, "bad NewIndex message test")


	// good message test
	goodMsg := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		0, // ids start from 0
		entries,
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, &res)

	follower.ApplyAppendEntries(goodMsg)

	time.Sleep(1)
	require.Equal(t, false, res.Heartbeat, "good message test")
	require.Equal(t, true, res.Appended, "good message test")
	require.Equal(t, goodMsg.NewIndex, res.TopIndex, "good message test")


	// bad PrevTerm message test
	badPrevTermMsg := message.NewAppendEntries(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		12, // more than follower's log length(0)
		entries,
		"",
	)

	go GetRaftMsg(raftNode.RaftOut, &res)

	follower.ApplyAppendEntries(badPrevTermMsg)

	time.Sleep(1)
	require.Equal(t, false, res.Heartbeat, "bad PrevTerm message test")
	require.Equal(t, false, res.Appended, "bad PrevTerm message test")

}

// Processing of RequestVote has been already tested in node_test.go TestProcessRequestVote func
// We test Processing of RequestVote individually for candidate cause his behaviour differs
