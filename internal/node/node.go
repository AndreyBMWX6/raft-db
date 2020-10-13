package node

import (
	"net"

	"../config"
	"../message"
)

// Error Types
type AuthError struct {
	Msg string
}

func (e *AuthError) Error() string {
	return e.Msg
}

// RolePlayer interface realize Raft role behavior
type RolePlayer interface {
	ReleaseNode() *RaftCore
	PlayRole() RolePlayer
}

func RunRolePlayer(player RolePlayer) {
	for {
		player = player.PlayRole()
	}
}

// Core of every Raft role
// Any RolePlayer implementation should contains *RaftCore
type RaftCore struct {
	Config *config.Config

	Addr net.Addr
	Neighbors []net.Addr

	Term int
	Entries []*message.Entry

	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage

	// Client IO
	ClientIn  <-chan message.ClientMessage
	ClientOut chan<- message.ClientMessage
}

// Message receivers wrappers
// Here there are message receiving logics
func (n *RaftCore) TryRecvRaftMsg() message.RaftMessage {
	select {
	case msg := <-n.RaftIn:
		return msg
	default:
		return nil
	}
}

func (n *RaftCore) TryRecvClientMsg() message.ClientMessage {
	select {
	case msg := <-n.ClientIn:
		return msg
	default:
		return nil
	}
}

// Message senders wrapper
// Here there are message sending logics
func (n *RaftCore) SendRaftMsg(msg message.RaftMessage) {
	n.RaftOut <- msg
}

func (n *RaftCore) SendClientMsg(msg message.ClientMessage)  {
	n.ClientOut <- msg
}
