package message

import "net"

type RequestAck struct {
	BaseRaftMessage
}

func NewRequestAck(base *BaseRaftMessage) *RequestAck {
	return &RequestAck{
		BaseRaftMessage: *base,
	}
}

func (ra *RequestAck) DestAddr() net.Addr {
	return ra.Dest
}

func (ra *RequestAck) OwnerAddr() net.Addr {
	return ra.Owner
}

func (ra *RequestAck) Term() int {
	return ra.CurrTerm
}

func (ra RequestAck) Type() int {
	return VoteType
	return VoteType
}

