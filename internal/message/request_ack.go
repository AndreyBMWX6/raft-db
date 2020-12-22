package message

import (
	"../net_message"
	"net"
)

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

func (ra *RequestAck) Unmarshal(message *net_message.Message) RaftMessage {
	switch raftMsg := message.RaftMessage.(type) {
	case *net_message.Message_RequestAck:
		// converting values
		ownerIp := net.IPv4(
			raftMsg.RequestAck.Msg.OwnerIp[0],
			raftMsg.RequestAck.Msg.OwnerIp[1],
			raftMsg.RequestAck.Msg.OwnerIp[2],
			raftMsg.RequestAck.Msg.OwnerIp[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(raftMsg.RequestAck.Msg.OwnerPort),
		}

		destIp := net.IPv4(
			raftMsg.RequestAck.Msg.DestIp[0],
			raftMsg.RequestAck.Msg.DestIp[1],
			raftMsg.RequestAck.Msg.DestIp[2],
			raftMsg.RequestAck.Msg.DestIp[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(raftMsg.RequestAck.Msg.DestPort),
		}

		return NewRequestAck(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: raftMsg.RequestAck.Msg.CurrTerm,
			},
			raftMsg.RequestAck.Voted,
		)

	default:
		return nil
	}
}
