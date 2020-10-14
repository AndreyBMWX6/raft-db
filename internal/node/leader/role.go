package leader

import (
	"../../node"
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
			// отправляем heartbeat
		default:
			if msg := l.core.TryRecvClientMsg(); msg != nil {

			}
			if msg := l.core.TryRecvRaftMsg(); msg != nil {

			}
		}
	}
}
