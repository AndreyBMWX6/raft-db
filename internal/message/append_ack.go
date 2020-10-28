package message

import "net"

type AppendAck struct {
	BaseRaftMessage

	Appended bool
}

func NewEntriesAck(base *BaseRaftMessage) *AppendAck {
	return &AppendAck{
		BaseRaftMessage: *base,
		Appended:    false,
	}
}

func (aa *AppendAck) DestAddr() *net.UDPAddr {
	return &aa.Dest
}

func (aa *AppendAck) OwnerAddr() *net.UDPAddr {
	return &aa.Owner
}

func (aa *AppendAck) Term() int {
	return aa.CurrTerm
}

func (aa AppendAck) Type() int {
	return AppendAckType
}

