package message

import "net"

type RequestVote struct {
	BaseRaftMessage
}

func NewRequestVote(base *BaseRaftMessage) *RequestVote {
	return &RequestVote{
		BaseRaftMessage: *base,
	}
}

func (rv *RequestVote) DestAddr() net.Addr {
	return rv.Dest
}

func (rv *RequestVote) OwnerAddr() net.Addr {
	return rv.Owner
}

func (rv *RequestVote) Term() int {
	return rv.CurrTerm
}

func (rv RequestVote) Type() int {
	return RequestVoteType
}

