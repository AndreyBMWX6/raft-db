package node

import (
	"fmt"
	"time"

	"../message"
)

type Candidate struct {
	core     *RaftCore
	timer    *time.Timer
	maxVotes int
	voters   map[string]struct{}
}

func BecomeCandidate(player RolePlayer) *Candidate {
	core := player.ReleaseNode()
	// logging
	fmt.Println(core.Addr, " became candidate")
	return NewCandidate(core)
}

func NewCandidate(core *RaftCore) *Candidate {
	maxVotes := (len(core.Neighbors)+1)/2
	return &Candidate{
		core:     core,
		timer:    time.NewTimer(core.Config.VotingTimeout),
		maxVotes: maxVotes,
		voters:   make(map[string]struct{}, maxVotes),
	}
}

func (c *Candidate) ReleaseNode() *RaftCore {
	c.timer.Stop()

	core := c.core
	c.core = nil
	return core
}

func (c *Candidate) PlayRole() RolePlayer {
	// Votes for itself
	c.voters[c.core.Addr.String()] = struct{}{}

	// implementation of parallel RequestVote


	for _, neighbor := range c.core.Neighbors {
		Topindex := len(c.core.Entries)
		Topterm := 0
		if Topindex != 0 {
			Topterm = c.core.Entries[len(c.core.Entries) - 1].Term
		}

		msg := message.NewRequestVote(
			&message.BaseRaftMessage{
				Owner: c.core.Addr,
				Dest: neighbor,
				CurrTerm: c.core.Term,
			},
			Topindex,
			Topterm,
	)
		msg.TopIndex = len(c.core.Entries)
		// if no Entries, Topterm = 0
		if msg.TopIndex != 0 {
			msg.TopTerm = c.core.Entries[len(c.core.Entries)-1].Term
		} else {
			msg.TopTerm = 0
		}
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
