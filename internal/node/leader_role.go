package node

import (
	"../message"
	"context"
	"fmt"
	"time"
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
	fmt.Println(core.Addr, " became leader")
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
			// отправляем heartbeat
			var prevTerm uint32 = 0 // no entries in leader case no need to assign newIdx as it's assigned by len(Entries)
			if len(l.core.Entries) > 1 {
				prevTerm = l.core.Entries[len(l.core.Entries) - 2].Term
			}
			for _,neighbor := range l.core.Neighbors {
				msg := message.NewAppendEntries(
					&message.BaseRaftMessage{
					Owner: l.core.Addr,
					Dest: neighbor,
					CurrTerm: l.core.Term,
					},
					prevTerm,
					uint32(len(l.core.Entries)),
					make([]*message.Entry, 0),
				)
				go l.core.SendRaftMsg(message.RaftMessage(msg))
			}

		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {
				// as heartbeat, but with data
			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {

			}
		}
	}
}
