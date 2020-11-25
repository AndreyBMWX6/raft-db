package node

import (
	"../config"
	"../message"
	"log"
	"net"

	"strconv"
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

	Voted bool

	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage

	// Client IO
	ClientIn  <-chan message.ClientMessage
	ClientOut chan<- message.ClientMessage
}

// may be change value initialization to address
func NewRaftCore() *RaftCore {
	cfg := config.NewConfig()
	return &RaftCore{
		Config:    cfg,
		Addr:      cfg.Addr,
		Neighbors: cfg.Neighbors,
		Term:      cfg.Term,
		Entries:   cfg.Entries,
		Voted:     cfg.Voted,
		RaftIn:    cfg.RaftIn,
		RaftOut:   cfg.RaftOut,
		ClientIn:  cfg.ClientIn,
		ClientOut: cfg.ClientOut,
	}
}

// Message receivers wrappers
// Here there are message receiving logics
func (n *RaftCore) TryRecvRaftMsg() message.RaftMessage {
	for {
		select {
		case msg := <-n.RaftIn:
			var msgType string
			switch msg.Type() {
			case message.AppendEntriesType:
				switch appendEntries := msg.(type) {
				case *message.AppendEntries:
					if len(appendEntries.Entries) == 0 {
						msgType = "Heartbeat"
					} else {
						msgType = "AppendEntries"
					}
				}
			case message.RequestVoteType:
				msgType = "RequestVote"
			case message.AppendAckType:
				msgType = "AppendAck"
			case message.RequestAckType:
				var Voted bool
				switch requestAck := msg.(type) {
				case *message.RequestAck:
					Voted = requestAck.Voted
				}
				msgType = strconv.FormatBool(Voted) + " RequestAck"
			}
			log.Println("Node:", msg.DestAddr().String(), " got ", msgType,
				" from Node:", msg.OwnerAddr().String())
			return msg
		default:
			return nil
		}
	}
}

func (n *RaftCore) TryRecvClientMsg() message.ClientMessage {
	for {
		select {
		case msg := <-n.ClientIn:
			return msg
		default:
			return nil
		}
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

	// node votes once in 1 term, so if node have already voted, it won't vote again
	if core.Voted {
		// made for clarity
		ack.Voted = false
	} else {
		// if no Entries, topTerm = 0
		var topTerm uint32 = 0
		if core.Entries != nil {
			topTerm = core.Entries[len(core.Entries) - 1].Term
		}
		if request.TopTerm < topTerm {
			ack.Voted = false
		} else {
			var topIndex = uint32(len(core.Entries))
			if (request.TopTerm == topTerm) && (request.TopIndex < topIndex) {
				ack.Voted = false
			} else {
				ack.Voted = true
				core.Voted = true
			}
		}
	}

	log.Println("Node:", ack.Owner.String(), " send", ack.Voted, " RequestAck to Node:", ack.Dest.String())
	core.SendRaftMsg(
		message.RaftMessage(ack),
	)
}
