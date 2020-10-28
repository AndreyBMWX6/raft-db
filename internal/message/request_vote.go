package message

import "net"

type RequestVote struct {
	BaseRaftMessage

	TopIndex int
	TopTerm  int
}

func NewRequestVote(base *BaseRaftMessage) *RequestVote {
	return &RequestVote{
		BaseRaftMessage: *base,
	}
}

func (rv *RequestVote) DestAddr() *net.UDPAddr {
	return &rv.Dest
}

func (rv *RequestVote) OwnerAddr() *net.UDPAddr {
	return &rv.Owner
}

func (rv *RequestVote) Term() int {
	return rv.CurrTerm
}

func (rv RequestVote) Type() int {
	return RequestVoteType
}
