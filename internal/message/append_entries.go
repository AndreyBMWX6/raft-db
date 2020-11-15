package message

import (
	"log"
	"net"
	"../net_message"
	"google.golang.org/protobuf/proto"
)

type Entry struct {
	Term uint32
	Query []byte
}

type AppendEntries struct {
	BaseRaftMessage

	PrevTerm  uint32
	NewIndex uint32
	Entries  []*Entry
}

func NewAppendEntries(base *BaseRaftMessage,
                      prevTerm uint32,
                      newIdx uint32,
                      entries []*Entry) *AppendEntries {
	return &AppendEntries{
		BaseRaftMessage: *base,
		PrevTerm:    prevTerm,
		NewIndex:    newIdx,
		Entries:     entries,
	}
}

func (e *AppendEntries) DestAddr() *net.UDPAddr {
	return &e.Dest
}

func (e *AppendEntries) OwnerAddr() *net.UDPAddr {
	return &e.Owner
}

func (e *AppendEntries) Term() uint32 {
	return e.CurrTerm
}

func (e AppendEntries) Type() int {
	return AppendEntriesType
}

func (e *AppendEntries) Unmarshal(data []byte) RaftMessage {
	AppendEntries := net_messages.AppendEntries{}
	err := proto.Unmarshal(data, &AppendEntries)
	if err == nil {
		// converting values
		ownerIp := net.IPv4(
			AppendEntries.Msg.Ownerip[0],
			AppendEntries.Msg.Ownerip[1],
			AppendEntries.Msg.Ownerip[2],
			AppendEntries.Msg.Ownerip[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(AppendEntries.Msg.Ownerport),
		}

		destIp := net.IPv4(
			AppendEntries.Msg.Dest[0],
			AppendEntries.Msg.Dest[1],
			AppendEntries.Msg.Dest[2],
			AppendEntries.Msg.Dest[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(AppendEntries.Msg.Destport),
		}

		entries := make([]*Entry, 0)
		for _, protoEntry := range AppendEntries.Entries {
			entry := &Entry{
				protoEntry.Term,
				protoEntry.Query,
			}
			entries = append(entries, entry)
		}

		return NewAppendEntries(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: AppendEntries.Msg.CurrTerm,
			},
			AppendEntries.PrevTerm,
			AppendEntries.NewIndex,
			entries,
		)
	}
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	return nil
}