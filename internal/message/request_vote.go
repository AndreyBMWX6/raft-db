package message

import "net"

type RequestVote struct {
	BaseRaftMessage

	TopIndex uint32
	TopTerm  uint32
}

func NewRequestVote(base *BaseRaftMessage,
					topindex uint32,
					topterm uint32) *RequestVote {
	return &RequestVote{
		BaseRaftMessage: *base,
		TopIndex: topindex,
		TopTerm: topterm,
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
