package message

import "net"

type RequestAck struct {
	BaseRaftMessage

	Voted bool
}

func NewRequestAck(base *BaseRaftMessage, voted bool) *RequestAck {
	return &RequestAck{
		BaseRaftMessage: *base,
		Voted: voted,
	}
}

func (ra *RequestAck) DestAddr() *net.UDPAddr {
	return &ra.Dest
}

func (ra *RequestAck) OwnerAddr() *net.UDPAddr {
	return &ra.Owner
}

func (ra *RequestAck) Term() uint32 {
	return ra.CurrTerm
}

func (ra RequestAck) Type() int {
	return RequestAckType
}

