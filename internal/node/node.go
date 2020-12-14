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
	AllConfig *config.AllConfig

	Addr net.UDPAddr
	Neighbors []net.UDPAddr
	URL string

	Term uint32
	Entries []*message.Entry

	Voted bool

	// Raft IO
	RaftIn  chan message.RaftMessage
	RaftOut chan message.RaftMessage

	// Client IO
	ClientIn  chan message.ClientMessage
	ClientOut chan message.ClientMessage

	// DB IO
	DBIn chan message.DBMessage
	DBOut chan message.DBMessage
}

// may be change value initialization to address
func NewRaftCore() *RaftCore {
	cfg := config.NewConfig()

	// Raft IO
	var raftIn  = make(chan message.RaftMessage)
	var raftOut = make(chan message.RaftMessage)

	// Client IO
	var clientIn  = make(chan message.ClientMessage)
	var clientOut = make(chan message.ClientMessage)

	// DB IO
	var dbIn  = make(chan message.DBMessage)
	var dbOut = make(chan message.DBMessage)

	return &RaftCore{
		Config    : cfg,
		Addr      : cfg.Addr,
		Neighbors : cfg.Neighbors,
		URL       : cfg.URL,
		Term      : cfg.Term,
		Entries   : cfg.Entries,
		Voted     : false,
		RaftIn    : raftIn,
		RaftOut   : raftOut,
		ClientIn  : clientIn,
		ClientOut : clientOut,
		DBIn      : dbIn,
		DBOut     : dbOut,
	}
}

func NewAllRunRaftCore(ip string, ipPort string, urlPort string) *RaftCore {
	cfg := config.NewConfig()
	allCfg := config.NewAllConfig()

	strAddr := ip + ":" + ipPort
	addr, err := net.ResolveUDPAddr("udp4", strAddr)
	if err != nil {
		log.Fatal(err)
	}

	var servId int
	for id, serv := range allCfg.Servers {
		if addr.String() == serv.String() {
			servId = id
		}
	}

	neighbors := allCfg.Servers[:servId]
	if servId != len(allCfg.Servers) - 1 {
		neighbors = append(neighbors, allCfg.Servers[servId + 1:]...)
	}

	url := "http://localhost:" + urlPort

	// Raft IO
	var raftIn  = make(chan message.RaftMessage)
	var raftOut = make(chan message.RaftMessage)

	// Client IO
	var clientIn  = make(chan message.ClientMessage)
	var clientOut = make(chan message.ClientMessage)

	// DB IO
	var dbIn  = make(chan message.DBMessage)
	var dbOut = make(chan message.DBMessage)

	return &RaftCore{
		Config    : cfg,
		AllConfig : allCfg,
		Addr      : *addr,
		Neighbors : neighbors,
		URL       : url,
		Term      : allCfg.Terms[servId],
		Entries   : allCfg.Entries[servId],
		Voted     : false,
		RaftIn    : raftIn,
		RaftOut   : raftOut,
		ClientIn  : clientIn,
		ClientOut : clientOut,
		DBIn      : dbIn,
		DBOut     : dbOut,
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
						msgType = "Heartbeat:"
					} else {
						msgType = "AppendEntries:"
					}
				}
			case message.RequestVoteType:
				msgType = "RequestVote:"
			case message.AppendAckType:
				switch appendAck := msg.(type) {
				case *message.AppendAck:
					msgType = strconv.FormatBool(appendAck.Appended) + " AppendAck:"
				default:
				}
			case message.RequestAckType:
				switch requestAck := msg.(type) {
				case *message.RequestAck:
					msgType = strconv.FormatBool(requestAck.Voted) + " RequestAck:"
				default:
				}
			}
			log.Println("Node:", msg.DestAddr().String(), " got ", msgType, msg.Term(),
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

func (n* RaftCore) TryRecvDBMsg() message.DBMessage {
	for {
		select {
		case msg := <- n.DBIn:
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

func (n *RaftCore) SendDBMsg(msg message.DBMessage)  {
	n.DBOut <- msg
}

func (core *RaftCore) ProcessRequestVote(request *message.RequestVote) bool {
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

	log.Println("Node:", ack.Owner.String(), " send", ack.Voted, " RequestAck:", ack.CurrTerm,
		" to Node:", ack.Dest.String())
	core.SendRaftMsg(
		message.RaftMessage(ack),
	)
	return ack.Voted
}
