package tests

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/AndreyBMWX6/raft-db/internal/message"
	"github.com/AndreyBMWX6/raft-db/internal/node"
)

func GetRaftMsg(ch chan message.RaftMessage, result *message.AppendAck, sync chan interface{}) {
	for {
		select {
		case msg := <-ch:
			switch ack := msg.(type) {
			case *message.AppendAck:
				*result = *ack
			}
			sync<-nil
			return
		}
	}
}

func Synchronize(sync chan interface{}, raftNode *node.RaftCore) {
	for {
		select {
		case <-sync:
			break
		default:
			raftNode.Config.Delay(raftNode.Config.DelayTime)
		}
		break
	}
	return
}

func TestProcessRequestVote(t *testing.T) {
	var raftNode = node.NewRaftCore("127.0.0.1", "8001", "8081")
	var sync chan interface{}

	msg := message.NewRequestVote(
		&message.BaseRaftMessage{
			Dest:     raftNode.Addr,
			CurrTerm: 0,
		},
		0,
		0,
	)


	// VOTED CHECKS TESTS

	// correct work test
	raftNode.Voted = false
	go GetRaftMsg(raftNode.RaftOut, nil, sync)

	result := raftNode.ProcessRequestVote(msg)
	Synchronize(sync, raftNode)
	require.Equal(t, true, result, "Correct work test")

	// node already voted test
	raftNode.Voted = true
	go GetRaftMsg(raftNode.RaftOut, nil, sync)

	result = raftNode.ProcessRequestVote(msg)

	Synchronize(sync, raftNode)
	require.Equal(t, false, result, "Node already voted test")
	raftNode.Voted = false

	// LOG METADATA CHECKS TESTS

	// same log meta test
	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	result = raftNode.ProcessRequestVote(msg)

	Synchronize(sync, raftNode)
	require.Equal(t, true, result, "Same log meta test")
	// if got true raftNode has voted so need to initiate voted with false back
	raftNode.Voted = false


	// adding entry to raftNode
	entry := &message.Entry {
		Term:  1,
		Query: nil,
	}
	raftNode.Entries = append(raftNode.Entries, entry)

	// less actual topTerm test
	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	result = raftNode.ProcessRequestVote(msg)

	Synchronize(sync, raftNode)
	require.Equal(t, false, result, "less actual topTerm test")

	// actual TopTerm, less actual TopIndex
	msg.TopTerm = raftNode.Entries[0].Term

	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	result = raftNode.ProcessRequestVote(msg)

	Synchronize(sync, raftNode)
	require.Equal(t, false, result, "actual TopTerm, less actual TopIndex test")

	// more actual log meta test
	msg.TopTerm  = raftNode.Entries[0].Term + 1
	msg.TopIndex = uint32(len(raftNode.Entries) + 1)

	go GetRaftMsg(raftNode.RaftOut, nil, sync)
	result = raftNode.ProcessRequestVote(msg)

	Synchronize(sync, raftNode)
	require.Equal(t, true, result, "more actual log meta test")
	raftNode.Voted = false
}
