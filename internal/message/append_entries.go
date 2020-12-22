package message

import (
	"net"

	"github.com/AndreyBMWX6/raft-db/internal/net_message"
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
	LeaderURL string
}

func NewAppendEntries(base *BaseRaftMessage,
                      prevTerm uint32,
                      newIdx uint32,
                      entries []*Entry,
                      leaderURL string) *AppendEntries {
	return &AppendEntries{
		BaseRaftMessage: *base,
		PrevTerm:    prevTerm,
		NewIndex:    newIdx,
		Entries:     entries,
		LeaderURL:   leaderURL,
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

func (e *AppendEntries) Unmarshal(message *net_message.Message) RaftMessage {
	switch raftMsg := message.RaftMessage.(type) {
	case *net_message.Message_AppendEntries:
		// converting values
		ownerIp := net.IPv4(
			raftMsg.AppendEntries.Msg.OwnerIp[0],
			raftMsg.AppendEntries.Msg.OwnerIp[1],
			raftMsg.AppendEntries.Msg.OwnerIp[2],
			raftMsg.AppendEntries.Msg.OwnerIp[3])
		ownerUdp := net.UDPAddr{
			IP:   ownerIp,
			Port: int(raftMsg.AppendEntries.Msg.OwnerPort),
		}

		destIp := net.IPv4(
			raftMsg.AppendEntries.Msg.DestIp[0],
			raftMsg.AppendEntries.Msg.DestIp[1],
			raftMsg.AppendEntries.Msg.DestIp[2],
			raftMsg.AppendEntries.Msg.DestIp[3])
		destUdp := net.UDPAddr{
			IP:   destIp,
			Port: int(raftMsg.AppendEntries.Msg.DestPort),
		}

		entries := make([]*Entry, 0)
		if raftMsg.AppendEntries.Entries == nil {
			entries = nil
		} else {
			for _, protoEntry := range raftMsg.AppendEntries.Entries {
				entry := &Entry{
					protoEntry.Term,
					protoEntry.Query,
				}
				entries = append(entries, entry)
			}
		}

		leaderURL := raftMsg.AppendEntries.URL

		return NewAppendEntries(
			&BaseRaftMessage{
				Owner:    ownerUdp,
				Dest:     destUdp,
				CurrTerm: raftMsg.AppendEntries.Msg.CurrTerm,
			},
			raftMsg.AppendEntries.PrevTerm,
			raftMsg.AppendEntries.NewIndex,
			entries,
			leaderURL,
		)

	default:
		return nil
	}
}
