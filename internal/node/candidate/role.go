package candidate

import (
	"net"
	"time"

	"../../node"
)

type Candidate struct {
	core     *node.RaftCore
	timer    *time.Timer
	maxVotes int
	voters   map[net.Addr]struct{}
}

func BecomeCandidate(player node.RolePlayer) *Candidate {
	core := player.ReleaseNode()
	return NewCandidate(core)
}

func NewCandidate(core *node.RaftCore) *Candidate {
	maxVotes := (len(core.Neighbors)+1)/2
	return &Candidate{
		core:     core,
		timer:    time.NewTimer(core.Config.VotingTimeout),
		maxVotes: maxVotes,
		voters:   make(map[net.Addr]struct{}, maxVotes),
	}
}

func (c *Candidate) ReleaseNode() *node.RaftCore {
	c.timer.Stop()

	core := c.core
	c.core = nil
	return core
}

func (c *Candidate) PlayRole() node.RolePlayer {
	// Votes for itself
	c.voters[c.core.Addr] = struct{}{}

	for {
		select {
		// повторное голосование начинается не по истечении времени таймера кандидата,
		// а в случае, когда никто не набрал наибольшее кол-во голосов, поэтому кандидат
		//  должен ожидать остальных. Реализовать это.
		case <-c.timer.C:
			return BecomeCandidate(c)
		default:
			if msg := c.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := c.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
		}
	}
}
