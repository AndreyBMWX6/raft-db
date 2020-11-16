package message

import (
	"log"
	"net"
	"../net_message"
	"google.golang.org/protobuf/proto"
)

type RequestVote struct {
	BaseRaftMessage

	TopIndex uint32
	TopTerm  uint32
}

func NewRequestVote(base *BaseRaftMessage,
					topIndex uint32,
					topTerm uint32) *RequestVote {
	return &RequestVote{
		BaseRaftMessage: *base,
		TopIndex:        topIndex,
		TopTerm:         topTerm,
	}
}

func (rv *RequestVote) DestAddr() *net.UDPAddr {
	return &rv.Dest
}

func (rv *RequestVote) OwnerAddr() *net.UDPAddr {
	return &rv.Owner
}

func (rv *RequestVote) Term() uint32 {
	return rv.CurrTerm
}

func (rv RequestVote) Type() int {
	return RequestVoteType
}

func (rv *RequestVote) Unmarshal(data []byte) RaftMessage {
	requestVote := net_messages.RequestVote{}
	err := proto.Unmarshal(data, &requestVote)
	if err == nil {
		// converting values
		ownerIp := net.IPv4(
			requestVote.Msg.Ownerip[0],
			requestVote.Msg.Ownerip[1],
			requestVote.Msg.Ownerip[2],
			requestVote.Msg.Ownerip[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(requestVote.Msg.Ownerport),
		}

		destIp := net.IPv4(
			requestVote.Msg.Dest[0],
			requestVote.Msg.Dest[1],
			requestVote.Msg.Dest[2],
			requestVote.Msg.Dest[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(requestVote.Msg.Destport),
		}

		return NewRequestVote(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: requestVote.Msg.CurrTerm,
			},
			requestVote.TopIndex,
			requestVote.TopTerm,
		)

	}
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	return nil
}
