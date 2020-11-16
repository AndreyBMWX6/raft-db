package node

import (
	"log"
	"net"
	"time"
)

// Follower follows current Leader
type Follower struct {
	core *RaftCore

	timer *time.Timer
	leaderAddr net.UDPAddr
}

func BecomeFollower(player RolePlayer, leader *net.UDPAddr) *Follower {
	core := player.ReleaseNode()
	return &Follower{
		core:       core,
		timer:      time.NewTimer(core.Config.FollowerTimeout),
		leaderAddr: *leader,
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
			log.Println("[follower  -> candidate]")
			return BecomeCandidate(f)
		default:
			if msg := f.core.TryRecvClientMsg(); msg != nil {
				f.ApplyClientMessage(msg)
			}
			if  msg := f.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := f.ApplyRaftMessage(msg); msg != nil {
					return nextRole
				}
			}
		}
	}
}
