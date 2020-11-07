package message

import "net"

// Raft Message Types
const (
	// Push Types
	AppendEntriesType = iota
	RequestVoteType

	// Pull Types
	AppendAckType
	RequestAckType
)

// Client Message Types
const (
	// Pull Types
	RawClientMessageType = iota
	PreprocessedClientMessageType

	// Push Types
	ResponseClientAnswer
)

// Inter-node Message Type
// Nodes communicate each other with
// messages realizing that interface
type RaftMessage interface {
	OwnerAddr() *net.UDPAddr
	DestAddr () *net.UDPAddr

	Term() int
	Type() int
}

// Inter-client Message Type
// Clients communicate nodes with these
// messages (also nodes can redirect them
// to other nodes)
type ClientMessage interface {
	ClientAddr() net.Addr
	DestAddr()   net.Addr

	Type() int
}

// Any inter-node message bases on it structure
type BaseRaftMessage struct {
	Owner net.UDPAddr
	Dest  net.UDPAddr
	CurrTerm int
}
