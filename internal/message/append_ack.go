package message

import (
	"net"

	"github.com/AndreyBMWX6/raft-db/internal/net_message"
)

type AppendAck struct {
	BaseRaftMessage

	Appended bool
	// field defining heartbeat or append
	Heartbeat bool
	// field defining index of added entry
	TopIndex uint32
}

func NewAppendAck(base *BaseRaftMessage, appended bool, heartbeat bool, topIndex uint32) *AppendAck {
	return &AppendAck{
		BaseRaftMessage: *base,
		Appended:    	 appended,
		Heartbeat: 		 heartbeat,
		TopIndex:        topIndex,
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

func (aa *AppendAck) Unmarshal(message *net_message.Message) RaftMessage {
	switch raftMsg := message.RaftMessage.(type) {
	case *net_message.Message_AppendAck:
		// converting values
		ownerIp := net.IPv4(
			raftMsg.AppendAck.Msg.OwnerIp[0],
			raftMsg.AppendAck.Msg.OwnerIp[1],
			raftMsg.AppendAck.Msg.OwnerIp[2],
			raftMsg.AppendAck.Msg.OwnerIp[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(raftMsg.AppendAck.Msg.OwnerPort),
		}

		destIp := net.IPv4(
			raftMsg.AppendAck.Msg.DestIp[0],
			raftMsg.AppendAck.Msg.DestIp[1],
			raftMsg.AppendAck.Msg.DestIp[2],
			raftMsg.AppendAck.Msg.DestIp[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(raftMsg.AppendAck.Msg.DestPort),
		}

		return NewAppendAck(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: raftMsg.AppendAck.Msg.CurrTerm,
			},
			raftMsg.AppendAck.Appended,
			raftMsg.AppendAck.Heartbeat,
			raftMsg.AppendAck.TopIndex,
		)

	default:
		return nil
	}
}
