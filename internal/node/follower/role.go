package follower

import (
	"net"
	"time"

	"../../node"
	"../candidate"
)

// Follower follows current Leader
type Follower struct {
	core *node.RaftCore

	timer *time.Timer
	leaderAddr net.Addr
}

func BecomeFollower(player node.RolePlayer, leader net.Addr) *Follower {
	core := player.ReleaseNode()
	return &Follower{
		core:       core,
		timer:      time.NewTimer(core.Config.FollowerTimeout),
		leaderAddr: leader,
	}
}

func (f *Follower) ReleaseNode() *node.RaftCore {
	f.timer.Stop()

	core := f.core
	f.core = nil

	return core
}

func (f *Follower) PlayRole() node.RolePlayer {
	for {
		select {
		case <-f.timer.C:
			return candidate.BecomeCandidate(f)
		default:
			if msg := f.core.TryRecvClientMsg(); msg != nil {
				f.ApplyClientMessage(msg)
			}
			if msg := f.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := f.ApplyRaftMessage(msg); msg != nil {
					return nextRole
				}
			}
		}
	}
}
