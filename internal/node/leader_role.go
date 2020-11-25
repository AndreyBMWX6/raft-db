package node

import (
	"context"
	"log"
	"time"

	"../message"
)

type FollowerView struct {
	LogIndex int
	
}

type Leader struct {
	core *RaftCore
	ctx  context.Context
	heartbeat *time.Ticker
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
	for {
		select {
		case <-l.heartbeat.C:
			// sending heartbeat
			// no entries in leader case no need to assign newIdx as it's assigned by len(Entries)
			var prevTerm uint32 = 0
			if len(l.core.Entries) > 1 {
				prevTerm = l.core.Entries[len(l.core.Entries) - 2].Term
			}
			for _, neighbour := range l.core.Neighbors {
				msg := message.NewAppendEntries(
					&message.BaseRaftMessage{
					Owner:    l.core.Addr,
					Dest:     neighbour,
					CurrTerm: l.core.Term,
					},
					prevTerm,
					uint32(len(l.core.Entries)),
					make([]*message.Entry, 0),
				)
				var msgType string
				if len(msg.Entries) == 0 {
					msgType = "Heartbeat:"
				} else {
					msgType = "AppendEntries:"
				}

					log.Println("Node:", msg.Owner.String(), " send ", msgType, msg.CurrTerm,
						" to Node:", msg.Dest.String())
				go l.core.SendRaftMsg(message.RaftMessage(msg))
			}

		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {
				// as heartbeat, but with data
			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := l.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
		}
	}
}
