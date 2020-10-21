package node

import (
	"fmt"
	"net"
	"time"
)

// Follower follows current Leader
type Follower struct {
	core *RaftCore

	timer *time.Timer
	leaderAddr net.Addr
}

func BecomeFollower(player RolePlayer, leader net.Addr) *Follower {
	core := player.ReleaseNode()
	// logging
	fmt.Println(core.Addr, " became follower")
	return &Follower{
		core:       core,
		timer:      time.NewTimer(core.Config.FollowerTimeout),
		leaderAddr: leader,
	}
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
			// logging
			fmt.Print("follower: ")
			return BecomeCandidate(f)
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
