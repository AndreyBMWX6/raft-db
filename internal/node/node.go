package node

import (
	"../config"
	"../message"
	"net"
	"time"
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

	Addr net.UDPAddr
	Neighbors []net.UDPAddr

	Term uint32
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

func (core *RaftCore) ProcessRequestVote(request *message.RequestVote) {
	ack := message.NewRequestAck(
		&message.BaseRaftMessage{
			Owner:    core.Addr,
			Dest:     request.Owner,
			CurrTerm: core.Term,
		}, false,
	)

	// if no Entries, Topterm = 0
	var topterm uint32 = 0
	if core.Entries != nil {
		topterm = core.Entries[len(core.Entries) - 1].Term
	}
	if request.TopTerm < topterm {
		// made for clarity
		ack.Voted = false
	} else {
		var topindex = uint32(len(core.Entries))
		if (request.TopTerm == topterm) && (request.TopIndex < topindex) {
			// made for clarity
			ack.Voted = false
		} else {
			ack.Voted = true
		}
	}

	core.SendRaftMsg(
		message.RaftMessage(ack),
	)

	if (ack.Voted) {
		time.Sleep(core.Config.VotingTimeout)

	}
}
