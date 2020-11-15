package message

import (
	"log"
	"net"
	"../net_message"
	"google.golang.org/protobuf/proto"
)

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

func (aa *AppendAck) Unmarshal(data []byte) RaftMessage {
	appendAck := net_messages.AppendAck{}
	err := proto.Unmarshal(data, &appendAck)
	if err == nil {
		// converting values
		ownerIp := net.IPv4(
			appendAck.Msg.Ownerip[0],
			appendAck.Msg.Ownerip[1],
			appendAck.Msg.Ownerip[2],
			appendAck.Msg.Ownerip[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(appendAck.Msg.Ownerport),
		}

		destIp := net.IPv4(
			appendAck.Msg.Dest[0],
			appendAck.Msg.Dest[1],
			appendAck.Msg.Dest[2],
			appendAck.Msg.Dest[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(appendAck.Msg.Destport),
		}

		return NewEntriesAck(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: appendAck.Msg.CurrTerm,
			},
			appendAck.Appended,
		)
	}
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	return nil
}
