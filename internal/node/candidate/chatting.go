package candidate

import (
	"../../message"
	"../../node"
	"../follower"
	"../leader"
	"../voter"
)

func (c *Candidate) ApplyRaftMessage(msg message.RaftMessage) node.RolePlayer {
	if msg.Term() < c.core.Term {
		return nil
	}

	if msg.Term() > c.core.Term {
		switch msg.Type() {
		case message.AppendEntriesType:
			return follower.BecomeFollower(c, msg.OwnerAddr())
		case message.RequestVoteType:
			return voter.BecomeVoter(c, msg.OwnerAddr())
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
			return leader.BecomeLeader(c)
		}
	}

	return nil
}
