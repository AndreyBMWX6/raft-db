package node

import (
	"log"
	"time"
)

// Follower follows current Leader
type Follower struct {
	core *RaftCore

	leaderURL string
	timer *time.Timer
}

func BecomeFollower(player RolePlayer) *Follower {
	core := player.ReleaseNode()
	return &Follower{
		core:       core,
		timer:      time.NewTimer(core.Config.FollowerTimeout),
	}
}

func RefreshFollower(f *Follower) *Follower{
	f.timer = time.NewTimer(f.core.Config.FollowerTimeout)
	return f
}

func (f *Follower) ReleaseNode() *RaftCore {
	f.timer.Stop()

	core := f.core
	f.core = nil

	return core
}

func (f *Follower) PlayRole() RolePlayer {
	for {
		select {
		case <-f.timer.C:
			log.Println("follower time is out")
			log.Println("[follower:", f.core.Term, "  -> candidate:", f.core.Term + 1, "]")
			return BecomeCandidate(f)
		default:
			if msg := f.core.TryRecvClientMsg(); msg != nil {
				f.ApplyClientMessage(msg)
			}
			if  msg := f.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := f.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
			if msg := f.core.TryRecvDBMsg(); msg != nil {
				f.ApplyDBMessage(msg)
			}
		}
	}
}
