package tests

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/AndreyBMWX6/raft-db/internal/message"
	"github.com/AndreyBMWX6/raft-db/internal/node"
)

func TestCandidateApplyLowerTermMessage(t *testing.T) {
	// each node in any role ignores any raft messages
	// as ApplyRaftMessage gets message.RaftMessage as parameter
	// and calls Term() method, that return msg.Term for any type of raft message

	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	raftNode.Term = 1
	var candidate = node.NewCandidate(raftNode)

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
	result := candidate.ApplyRaftMessage(lowAppEnt)

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
	result = candidate.ApplyRaftMessage(lowReqVot)

	require.Equal(t, nil, result, "RequestVote with lower term test")

	// RequestAck with lower term test
	lowReqAck := message.NewRequestAck(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0, // lower than candidate term(1)
		},
		false,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil cause message term is lower than candidate term
	result = candidate.ApplyRaftMessage(lowReqAck)

	require.Equal(t, nil, result, "RequestAck with lower term test")

	// Candidate don't process AppendAck
}


func TestCandidateApplyAppendEntries(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var candidate = node.NewCandidate(raftNode)

	// in AppendEntries case candidate becomes follower regardless of message term
	// it can be equal or higher


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

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil as candidate votes ony for candidate of bigger term
	result := candidate.ApplyRaftMessage(eqMsg)

	var resultIsFollower bool

	switch result.(type) {
	case *node.Follower:
		resultIsFollower = true
	default:
		resultIsFollower = false
	}

	require.Equal(t, true, resultIsFollower, "message with equal term test")

	require.Equal(t, eqMsg.Term(), result.ReleaseNode().Term, "message with equal term test")


	//message with higher term test
	// need to call NewCandidate helper cause candidate var is follower now
	candidate = node.NewCandidate(raftNode)

	highMsg := eqMsg
	highMsg.CurrTerm = 1 // higher than node term(0)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return *node.Follower as next role
	// you can see role change in logs
	// also will update node term
	result = candidate.ApplyRaftMessage(highMsg)

	switch result.(type) {
	case *node.Follower:
		resultIsFollower = true
	default:
		resultIsFollower = false
	}

	require.Equal(t, true, resultIsFollower, "message with higher term test")
	require.Equal(t, highMsg.Term(), result.ReleaseNode().Term, "message with higher term test")
}

func TestCandidateApplyRequestVote(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var candidate = node.NewCandidate(raftNode)

	// message with equal term test
	eqMsg := message.NewRequestVote(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		0,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil as candidate votes ony for candidate of bigger term
	result := candidate.ApplyRaftMessage(eqMsg)

	require.Equal(t, nil, result, "message with equal term test")

	//message with higher term test
	highMsg := eqMsg
	highMsg.CurrTerm = 1 // higher than node term(0)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return *node.Follower as next role
	// you can see role change in logs
	// also will update node term
	result = candidate.ApplyRaftMessage(highMsg)

	var resultIsFollower bool

	switch result.(type) {
	case *node.Follower:
		resultIsFollower = true
	default:
		resultIsFollower = false
	}

	require.Equal(t, true, resultIsFollower, "message with higher term test")
	require.Equal(t, highMsg.Term(), result.ReleaseNode().Term, "message with higher term test")
}

func TestCandidateApplyRequestAck(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var candidate = node.NewCandidate(raftNode)

	// false RequestAck with equal term test
	falEqMsg := message.NewRequestAck(
		&message.BaseRaftMessage{
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		false,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil as node votes against candidate
	result := candidate.ApplyRaftMessage(falEqMsg)

	require.Equal(t, nil, result, "false RequestAck with equal term test")

	// true RequestAck with equal term test
	oldVotersCount := len(candidate.GetVoters())

	owner := raftNode.Neighbors[0]
	truEqMsg := message.NewRequestAck(
		&message.BaseRaftMessage{
			Owner: owner,
			Dest: raftNode.Addr,
			CurrTerm: 0,
		},
		true,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil as node votes against candidate
	result = candidate.ApplyRaftMessage(truEqMsg)

	voters := candidate.GetVoters()

	require.Equal(t, nil, result, "true RequestAck with equal term test")
	require.Equal(t, oldVotersCount + 1, len(voters), "true RequestAck with equal term test")
	require.Equal(t, struct{}{}, voters[owner.String()], "true RequestAck with equal term test")

	// true RequestAck with equal term test
	oldVotersCount = len(candidate.GetVoters())

	// using same message from prev test

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return nil as node votes against candidate
	result = candidate.ApplyRaftMessage(truEqMsg)

	voters = candidate.GetVoters()

	require.Equal(t, nil, result, "true RequestAck with equal term test")
	require.Equal(t, oldVotersCount, len(voters), "true RequestAck with equal term test")
	require.Equal(t, struct{}{}, voters[owner.String()], "true RequestAck with equal term test")


	// true RequestAck with higher term test
	truHiMsg := message.NewRequestAck(
		&message.BaseRaftMessage{
			Owner: owner,
			Dest: raftNode.Addr,
			CurrTerm: 1,
		},
		false,
	)

	go GetRaftMsg(raftNode.RaftOut, nil)
	// will return follower as nextRole
	result = candidate.ApplyRaftMessage(truHiMsg)

	var resultIsFollower bool

	switch result.(type) {
	case *node.Follower:
		resultIsFollower = true
	default:
		resultIsFollower = false
	}

	require.Equal(t, true, resultIsFollower, "true RequestAck with higher term test")
}
