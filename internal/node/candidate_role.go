package node

import (
	"log"
	"math/rand"
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
	core.Term++
	return NewCandidate(core)
}

func NewCandidate(core *RaftCore) *Candidate {
	maxVotes := (len(core.Neighbors)+1)/2
	return &Candidate{
		core:     core,
		timer:    time.NewTimer(time.Duration(rand.Intn(1000) + 100)*time.Millisecond),
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
	c.core.Voted = true

	// implementation of parallel RequestVote
	for _, neighbour := range c.core.Neighbors {
		var topIndex = uint32(len(c.core.Entries))
		var topTerm uint32 = 0
		if topIndex != 0 {
			topTerm = c.core.Entries[len(c.core.Entries) - 1].Term
		}

		msg := message.NewRequestVote(
			&message.BaseRaftMessage{
				Owner:    c.core.Addr,
				Dest:     neighbour,
				CurrTerm: c.core.Term,
			},
			topIndex,
			topTerm,
	)
		//msg.TopIndex = uint32(len(c.core.Entries))
		// if no Entries, topTerm = 0
		if msg.TopIndex != 0 {
			msg.TopTerm = c.core.Entries[len(c.core.Entries)-1].Term
		} else {
			msg.TopTerm = 0
		}

		log.Println("Node:", msg.Owner.String(), " send RequestVote:", msg.CurrTerm,
			" to Node:", msg.Dest.String())
		go c.core.SendRaftMsg(msg)
	}

	for {
		select {
		case <-c.timer.C:
			log.Println("voting time is out")
			log.Println("[candidate:", c.core.Term, " -> candidate:", c.core.Term + 1, "]")
			return BecomeCandidate(c)

		default:
			//if msg := c.core.TryRecvClientMsg(); msg != nil {
				//с.ApplyClientMessage(msg)
			//}
			if  msg := c.core.TryRecvRaftMsg(); msg != nil {
				if nextRole := c.ApplyRaftMessage(msg); nextRole != nil {
					return nextRole
				}
			}
		}
	}
}
