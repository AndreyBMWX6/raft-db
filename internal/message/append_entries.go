package message

import "net"

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

