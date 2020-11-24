package node

import (
	"../message"
	"log"
	"time"
)

func (c *Candidate) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < c.core.Term {
		return nil
	}

	if msg.Term() >= c.core.Term {
		switch msg.Type() {
		case message.AppendEntriesType:
			// may be will add processing of query
			log.Println("[candidate:", c.core.Term, " -> follower:,", msg.Term(), " ]")
			c.core.Term = msg.Term()
			return BecomeFollower(c, msg.OwnerAddr())
		case message.RequestVoteType:
			if msg.Term() > c.core.Term {
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
					log.Println("[candidate:", c.core.Term, " -> follower:", msg.Term(), " ]")
					c.core.Term = msg.Term()
					c.core.ProcessRequestVote(request)
					c.timer = time.NewTimer(c.core.Config.HeartbeatTimeout)
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
				log.Println("[candidate:", c.core.Term, " -> leader:", c.core.Term, "   ]")
				return BecomeLeader(c)
			}

		default:
			return nil
		}
	}

	return nil
}
