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

	Term() uint32
	Type() int
	Unmarshal([]byte) RaftMessage
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
	CurrTerm uint32
}

type BaseClientMessage struct {
	Owner net.Addr
	Dest  net.Addr
}

type RawClientMessage struct {
	BaseClientMessage

	Entry *Entry
}

func NewRawClientMessage(base *BaseClientMessage,
					  entry *Entry) *RawClientMessage {
	return &RawClientMessage{
		BaseClientMessage: *base,
		Entry:           entry,
	}
}


func (rc *RawClientMessage) ClientAddr() net.Addr {
	return rc.Owner
}

func (rc *RawClientMessage) DestAddr() net.Addr {
	return rc.Dest
}

func (rc *RawClientMessage) Type() int {
	return RawClientMessageType
}

type ResponseClientMessage struct {
	BaseClientMessage
}

func NewResponseClientMessage(base *BaseClientMessage) *RawClientMessage {
	return &RawClientMessage{
		BaseClientMessage: *base,
	}
}

func (rc *ResponseClientMessage) ClientAddr() net.Addr {
	return rc.Owner
}

func (rc *ResponseClientMessage) DestAddr() net.Addr {
	return rc.Dest
}

func (rc *ResponseClientMessage) Type() int {
	return ResponseClientAnswer
}
