package node

import (
	"../message"
	"log"
)

func (c *Candidate) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < c.core.Term {
		return nil
	}

	if msg.Term() >= c.core.Term {
		switch msg.Type() {
		case message.AppendEntriesType:
			// may be will add processing of query
			log.Println("[candidate -> follower ]")
			return BecomeFollower(c, msg.OwnerAddr())
		case message.RequestVoteType:
			if msg.Term() > c.core.Term {
				c.core.Term = msg.Term()
				switch requestVote := msg.(type) {
				case *message.RequestVote:
					request := message.NewRequestVote(
						&message.BaseRaftMessage{
							Owner:    *msg.OwnerAddr(),
							Dest:     *msg.DestAddr(),
							CurrTerm: msg.Term(),
						},
						requestVote.TopIndex,
						requestVote.TopTerm,
					)
					c.core.ProcessRequestVote(request)
					log.Println("[candidate -> follower ]")
					return BecomeFollower(c, msg.OwnerAddr())
				default:
					log.Print("`RequestVoteMessage` expected, got another type")
				}
			}
		case message.RequestAckType:
			if _, found := c.voters[msg.OwnerAddr().String()]; found {
				return nil
			} else {
				c.voters[msg.OwnerAddr().String()] = struct{}{}
			}

			if len(c.voters) >= c.maxVotes {
				// logging
				log.Println("[candidate -> leader   ]")
				return BecomeLeader(c)
			}

		default:
			return nil
		}
	}

	return nil
}
