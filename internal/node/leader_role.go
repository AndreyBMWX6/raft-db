package node

import (
	"context"
	"log"
	"net"
	"time"

	"../message"
)

type Leader struct {
	core *RaftCore
	ctx  context.Context
	heartbeat *time.Ticker
	updates map[string]chan[]*message.Entry

	// needed to define, when more than half committed and send response to client
	replicated []int
	// needed to send DB response to client
	responses []*message.DBResponse
}

func BecomeLeader(player RolePlayer) *Leader {
	core := player.ReleaseNode()
	return &Leader{
		core:       core,
		heartbeat:  time.NewTicker(core.Config.HeartbeatTimeout),
		ctx:        context.Background(),
		updates:    make(map[string]chan []*message.Entry, len(core.Neighbors)),
		replicated: make([]int, 0),
		responses:  make([]*message.DBResponse, 0),
	}
}

func (l *Leader) ReleaseNode() *RaftCore {
	core := l.core
	l.core = nil
	return core
}

func (l *Leader) PlayRole() RolePlayer {
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

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
				l.ApplyClientMessage(msg)
			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := l.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
			if msg := l.core.TryRecvDBMsg(); msg != nil {
				l.ApplyDBMessage(msg)
			}

			Delay()
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
					l.core.URL,
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
			}
		}
	}
}
