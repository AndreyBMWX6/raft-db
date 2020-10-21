package node

import (
	"../message"
	"fmt" // temporary implementation of logging cluster nodes condition and communication
)

func (c *Candidate) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < c.core.Term {
		return nil
	}

	if msg.Term() >= c.core.Term {
		switch msg.Type() {
		case message.AppendEntriesType:
			// may be will add processing of query
			// logging
			fmt.Print("candidate:")
			return BecomeFollower(c, msg.OwnerAddr())
		case message.RequestVoteType:
			if msg.Term() > c.core.Term {
				request := message.NewRequestVote(
					&message.BaseRaftMessage{
						Owner:	  msg.OwnerAddr(),
						Dest: 	  msg.DestAddr(),
						CurrTerm: msg.Term(),
					},
				)
				// 2 lower lines can be changed to BecomeVoter
				c.core.ProcessRequestVote(request)
				// logging
				fmt.Print("candidate:")
				return BecomeFollower(c, msg.OwnerAddr())
			}
		default:
			return nil
		}
	}

	switch msg.Type() {
	case message.VoteType:
		if _, found := c.voters[msg.OwnerAddr()]; found {
			return nil
		} else {
			c.voters[msg.OwnerAddr()] = struct{}{}
		}

		if len(c.voters) >= c.maxVotes {
			// logging
			fmt.Print("candidate:")
			return BecomeLeader(c)
		}
	}

	return nil
}
