package message

import (
	"net"

	"github.com/AndreyBMWX6/raft-db/internal/net_message"
)

type RequestVote struct {
	BaseRaftMessage

	TopIndex uint32
	TopTerm  uint32
}

func NewRequestVote(base *BaseRaftMessage,
					topIndex uint32,
					topTerm uint32) *RequestVote {
	return &RequestVote{
		BaseRaftMessage: *base,
		TopIndex:        topIndex,
		TopTerm:         topTerm,
	}
}

func (rv *RequestVote) DestAddr() *net.UDPAddr {
	return &rv.Dest
}

func (rv *RequestVote) OwnerAddr() *net.UDPAddr {
	return &rv.Owner
}

func (rv *RequestVote) Term() uint32 {
	return rv.CurrTerm
}

func (rv RequestVote) Type() int {
	return RequestVoteType
}

func (rv *RequestVote) Unmarshal(message *net_message.Message) RaftMessage {
	switch raftMsg := message.RaftMessage.(type) {
	case *net_message.Message_RequestVote:
		// converting values
		ownerIp := net.IPv4(
			raftMsg.RequestVote.Msg.OwnerIp[0],
			raftMsg.RequestVote.Msg.OwnerIp[1],
			raftMsg.RequestVote.Msg.OwnerIp[2],
			raftMsg.RequestVote.Msg.OwnerIp[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(raftMsg.RequestVote.Msg.OwnerPort),
		}

		destIp := net.IPv4(
			raftMsg.RequestVote.Msg.DestIp[0],
			raftMsg.RequestVote.Msg.DestIp[1],
			raftMsg.RequestVote.Msg.DestIp[2],
			raftMsg.RequestVote.Msg.DestIp[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(raftMsg.RequestVote.Msg.DestPort),
		}

		return NewRequestVote(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: raftMsg.RequestVote.Msg.CurrTerm,
			},
			raftMsg.RequestVote.TopIndex,
			raftMsg.RequestVote.TopTerm,
		)

	default:
		return nil
	}
}
