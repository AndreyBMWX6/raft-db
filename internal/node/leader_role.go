package node

import (
	"context"
	"log"
	"net"
	"time"

	"../manager"
	"../message"
)

type FollowerView struct {
	LogIndex int
	
}

type Leader struct {
	core *RaftCore
	ctx  context.Context
	heartbeat *time.Ticker

	// needed to define, when more than half committed and send response yo client
	//replicated int
}

func BecomeLeader(player RolePlayer) *Leader {
	core := player.ReleaseNode()
	return &Leader{
		core:      core,
		heartbeat: time.NewTicker(core.Config.HeartbeatTimeout),
		ctx:       context.Background(),
	}
}

func (l *Leader) ReleaseNode() *RaftCore {
	core := l.core
	l.core = nil
	return core
}

func (l *Leader) PlayRole() RolePlayer {
	cm := &manager.ClientManager{
		ClientIn: l.core.Config.ClientOut,
		ClientOut: l.core.Config.ClientIn,
	}

	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	go cm.ProcessEntries(ctx)

	updates := make([]chan *message.Entry, 0)

	for _, neighbour := range l.core.Neighbors {
		update := make(chan *message.Entry)
		go NewReplicator(ctx, l, neighbour, update)()
		updates = append(updates, update)
	}

	for {
		select {
		case <-l.heartbeat.C:
			for _, update := range updates {
				update <- nil
			}
		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {
				switch rawClient := msg.(type) {
				case *message.RawClientMessage:
					rawClient.Entry.Term = l.core.Term
					l.core.Entries = append(l.core.Entries, rawClient.Entry)

					log.Println("Leader added new entry")
					log.Println("Leader log:    ", l.core.Entries)
					var entriesTerms []uint32
					for _,entry := range l.core.Entries {
						entriesTerms = append(entriesTerms, entry.Term)
					}
					log.Println("Log terms:     ", entriesTerms)

					var entries []*message.Entry
					entries = append(entries, rawClient.Entry)
					log.Println("Append entries:", entries)
					entriesTerms = nil
					for _,entry := range entries {
						entriesTerms = append(entriesTerms, entry.Term)
					}
					log.Println("Entries terms: ", entriesTerms)

					for _, update := range updates {
						update <-rawClient.Entry
					}
				default:
					log.Print("`RawClientMessage` expected, got another type")
				}
			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := l.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
		}
	}
}

func NewReplicator(ctx context.Context,
				   l *Leader,
				   follower net.UDPAddr,
				   update <-chan *message.Entry) func() {
	return func() {
		// helps switching between sending heartbeat and new entries
		var heartbeat = false

		var entries []*message.Entry
		var newEntries []*message.Entry

		for {
			select {
			case <-ctx.Done():
				return
			case u := <-update:
				if u == nil {
					heartbeat = true
				} else {
					heartbeat = false
					entries = append(entries, u)
					newEntries = append(newEntries, u)
				}

				// make loop sending appendEntries until got response
				if heartbeat == false {

				} else {
					newEntries = nil
				}

				var newIndex uint32 = 0
				if len(l.core.Entries) != 0 {
					newIndex = uint32(len(l.core.Entries) - 1)
				}
				var prevTerm uint32 = 0
				if newIndex > 0 {
					prevTerm = l.core.Entries[newIndex - 1].Term
				}

				msg := message.NewAppendEntries(
					&message.BaseRaftMessage{
						Owner:    l.core.Addr,
						Dest:     follower,
						CurrTerm: l.core.Term,
					},
					prevTerm,
					newIndex,
					newEntries,
				)
				var msgType string
				if len(msg.Entries) == 0 {
					msgType = "Heartbeat:"
				} else {
					msgType = "AppendEntries:"
				}

				log.Println("Node:", msg.Owner.String(), " send ", msgType, msg.CurrTerm,
					" to Node:", msg.Dest.String())

				// may be will log append entries for each node later to show entries in case,
				//when some nodes won't commit updates for the first time
				/*
				if msgType == "AppendEntries:" {
					log.Println("append entries:", msg.Entries)
					var entriesTerms []uint32
					for _,entry := range msg.Entries {
						entriesTerms = append(entriesTerms, entry.Term)
					}
					log.Println("entries terms: ", entriesTerms)
				}
				*/

				l.core.SendRaftMsg(message.RaftMessage(msg))
				newEntries = nil
			default:
			}
		}
	}
}
