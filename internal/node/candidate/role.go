package candidate

import (
	"net"
	"time"

	"../../node"
	"../../message"
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

	// implement parallel RequestVote
	//ctx, cancel := context.WithTimeout(context.Background(), c.core.Config.VotingTimeout)
	//  можно убрать таймер из кандидата и использовать контекст с таймером, а мб и нет ????

	for _, neighbor := range c.core.Neighbors {
		msg := message.NewRequestVote(
			&message.BaseRaftMessage{
				Owner: c.core.Addr,
				Dest: neighbor,
				CurrTerm: c.core.Term,
		},
	)
		msg.TopIndex = len(c.core.Entries)
		msg.TopTerm = c.core.Entries[len(c.core.Entries) - 1].Term
		go c.core.SendRaftMsg(msg)
	}

	for {
		select {
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
