package message

import "net"

type AppendAck struct {
	BaseRaftMessage

	Appended bool
}

func NewEntriesAck(base *BaseRaftMessage, appended bool) *AppendAck {
	return &AppendAck{
		BaseRaftMessage: *base,
		Appended:    appended,
	}
}

func (aa *AppendAck) DestAddr() *net.UDPAddr {
	return &aa.Dest
}

func (aa *AppendAck) OwnerAddr() *net.UDPAddr {
	return &aa.Owner
}

func (aa *AppendAck) Term() uint32 {
	return aa.CurrTerm
}

func (aa AppendAck) Type() int {
	return AppendAckType
}

