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
	updates map[string]chan[]*message.Entry

	// needed to define, when more than half committed and send response yo client
	replicated int
}

func BecomeLeader(player RolePlayer) *Leader {
	core := player.ReleaseNode()
	return &Leader{
		core:      core,
		heartbeat: time.NewTicker(core.Config.HeartbeatTimeout),
		ctx:       context.Background(),
		updates:   make(map[string]chan []*message.Entry, len(core.Neighbors)),
		replicated: 0,
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

	for _, neighbour := range l.core.Neighbors {
		update := make(chan []*message.Entry)
		l.updates[neighbour.String()] = update
		go NewReplicator(ctx, l, neighbour, update)()
	}

	for {
		select {
		case <-l.heartbeat.C:
			for _, update := range l.updates {
				update <- nil
			}
		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {
				switch rawClient := msg.(type) {
				case *message.RawClientMessage:
					rawClient.Entry.Term = l.core.Term
					l.core.Entries = append(l.core.Entries, rawClient.Entry)
					l.replicated++

					log.Println("Leader added new entry")
					log.Println("Leader log:     ", l.core.Entries)
					var entriesTerms []uint32
					for _,entry := range l.core.Entries {
						entriesTerms = append(entriesTerms, entry.Term)
					}
					log.Println("Log terms:      ", entriesTerms)

					var entries []*message.Entry
					entries = append(entries, rawClient.Entry)
					log.Println("Append entries: ", entries)
					entriesTerms = nil
					for _,entry := range entries {
						entriesTerms = append(entriesTerms, entry.Term)
					}
					log.Println("Entries terms:  ", entriesTerms)

					for _, update := range l.updates {
						upd := make([]*message.Entry, 1)
						upd[0] = rawClient.Entry
						update <-upd
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
				   update <-chan []*message.Entry) func() {
	return func() {
		// helps switching between sending heartbeat and new entries
		var heartbeat = false

		var entries []*message.Entry // current entries
		//var log []*message.Entry // full log
		var newIndex uint32
		var prevTerm uint32

		// one of the solutions is use
		// 1-elem slice made from 1 nil
		//for true condition of committed var
		// and 2-elem slice made from 2 nil for false
		// 0 length entries slice signals,
		//that node have committed changes
		//so entries can be cleared
		// extra bool chan can be used instead

		for {
			select {
			case <-ctx.Done():
				return
			case u := <-update:
				if u == nil {
					heartbeat = true
				} else {
					if len(u) == 1 && u[0] == nil {
						// clear entries request
						entries = nil
					} else {
						if len(u) == 2 && u[0] == nil && u[1] == nil {
							// retry append request
							log.Println("retry request")
							heartbeat = false

							if newIndex != 0 {
								newIndex--
								entries = l.core.Entries[newIndex:]
							}
							if newIndex > 0 {
								prevTerm = l.core.Entries[newIndex - 1].Term
							} else {
								prevTerm = 0
							}
						} else {
							heartbeat = false

							if len(l.core.Entries) == 0 {
								newIndex = 0
							} else {
								newIndex = uint32(len(l.core.Entries) - 1)
							}

							if newIndex < 1 {
								prevTerm = 0
							} else {
								prevTerm = l.core.Entries[newIndex - 1].Term
							}

							entries = l.core.Entries[newIndex:]
						}
					}
				}

				msg := message.NewAppendEntries(
					&message.BaseRaftMessage{
						Owner:    l.core.Addr,
						Dest:     follower,
						CurrTerm: l.core.Term,
					},
					prevTerm,
					newIndex,
					entries,
				)


				// make loop sending appendEntries until got response
				if heartbeat == true {
				//	msg.Entries = nil
				}


				var msgType string
				if len(msg.Entries) == 0 {
					msgType = "Heartbeat:"
				} else {
					msgType = "AppendEntries:"
				}

				log.Println("Node:", msg.Owner.String(), " send ", msgType, msg.CurrTerm,
					" to Node:", msg.Dest.String())


					if msgType == "AppendEntries:" && len(u) == 2 && u[0] == nil && u[1] == nil {
						// retry append request
						log.Println("Append entries: ", msg.Entries)
						var entriesTerms []uint32
						for _,entry := range msg.Entries {
							entriesTerms = append(entriesTerms, entry.Term)
						}
						log.Println("Entries terms:  ", entriesTerms)
					}


				l.core.SendRaftMsg(message.RaftMessage(msg))
			default:
			}
		}
	}
}
