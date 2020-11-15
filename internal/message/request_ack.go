package message

import (
	"log"
	"net"
	"../net_message"
	"google.golang.org/protobuf/proto"
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

func (ra *RequestAck) Unmarshal(data []byte) RaftMessage {
	requestAck := net_messages.RequestAck{}
	err := proto.Unmarshal(data, &requestAck)
	if err == nil {
		// converting values
		ownerIp := net.IPv4(
			requestAck.Msg.Ownerip[0],
			requestAck.Msg.Ownerip[1],
			requestAck.Msg.Ownerip[2],
			requestAck.Msg.Ownerip[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(requestAck.Msg.Ownerport),
		}

		destIp := net.IPv4(
			requestAck.Msg.Ownerip[0],
			requestAck.Msg.Ownerip[1],
			requestAck.Msg.Ownerip[2],
			requestAck.Msg.Ownerip[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(requestAck.Msg.Destport),
		}

		return NewRequestAck(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: requestAck.Msg.CurrTerm,
			},
			requestAck.Voted,
		)
	}
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	return nil
}
