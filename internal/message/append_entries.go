package message

import "net"

type Entry struct {
	Term int
	Query []byte
}

type AppendEntries struct {
	BaseRaftMessage

	PrevTerm  int
	NewIndex int
	Entries  []*Entry
}

func NewAppendEntries(base *BaseRaftMessage,
                      prevTerm,
                      newIdx int,
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

func (e *AppendEntries) Term() int {
	return e.CurrTerm
}

func (e AppendEntries) Type() int {
	return AppendEntriesType
}

