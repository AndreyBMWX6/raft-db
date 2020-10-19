package leader

import (
	"../../node"
	"../../message"
	"context"
	"time"
)

type FollowerView struct {
	LogIndex int
	
}

type Leader struct {
	core *node.RaftCore
	ctx  context.Context
	heartbeat *time.Ticker
}

func BecomeLeader(player node.RolePlayer) *Leader {
	core := player.ReleaseNode()
	return &Leader{
		core:      core,
		heartbeat: time.NewTicker(core.Config.HeartbeatTimeout),
		ctx:       context.Background(),
	}
}

func (l *Leader) ReleaseNode() *node.RaftCore {
	core := l.core
	l.core = nil
	return core
}

func (l *Leader) PlayRole() node.RolePlayer {
	for {
		select {
		case <-l.heartbeat.C:
			///////////////////////////////
			for _,neighbor := range l.core.Neighbors {
				msg := message.NewAppendEntries(
					&message.BaseRaftMessage{
					Owner: l.core.Addr,
					Dest: neighbor,
					CurrTerm: l.core.Term,
					},
					l.core.Entries[len(l.core.Entries) - 2].Term,
					len(l.core.Entries),
					make([]*message.Entry, 0),
				)
				go l.core.SendRaftMsg(message.RaftMessage(msg))
			}
			//////////////////////////////
			// отправляем heartbeat
		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {

			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {

			}
		}
	}
}
