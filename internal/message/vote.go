package message

import "net"

type Vote struct {
	BaseRaftMessage
}

func NewVote(base *BaseRaftMessage) *Vote {
	return &Vote{
		BaseRaftMessage: *base,
	}
}

func (v *Vote) DestAddr() net.Addr {
	return v.Dest
}

func (v *Vote) OwnerAddr() net.Addr {
	return v.Owner
}

func (v *Vote) Term() int {
	return v.CurrTerm
}

func (v Vote) Type() int {
	return VoteType
}
