package node

import (
	"log"

	"github.com/AndreyBMWX6/raft-db/internal/message"
)

func (c *Candidate) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < c.core.Term {
		return nil
	}

	if msg.Term() >= c.core.Term {
		oldCandidateTerm := c.core.Term
		switch msg.Type() {
		case message.AppendEntriesType:
			// may be will add processing of query
			log.Println("[candidate:", oldCandidateTerm, " -> follower:", msg.Term(), " ]")
			c.core.Term = msg.Term()
			return BecomeFollower(c)
		case message.RequestVoteType:
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

					if msg.Term() > c.core.Term {
						c.core.Term = msg.Term()
						c.core.Voted = false
					}
					c.core.ProcessRequestVote(request)

					if msg.Term() > oldCandidateTerm {
						log.Println("[candidate:", oldCandidateTerm, " -> follower:", msg.Term(), " ]")
						return BecomeFollower(c)
					} else {
						return nil
					}
				default:
					log.Print("`RequestVoteMessage` expected, got another type")
				}
		case message.RequestAckType:
			if msg.Term() > oldCandidateTerm {
				log.Println("[candidate:", oldCandidateTerm, " -> follower:", msg.Term(), " ]")
				return BecomeFollower(c)
			}
			switch requestAck := msg.(type) {
			case *message.RequestAck:
				if requestAck.Voted {
					if _, found := c.voters[msg.OwnerAddr().String()]; found {
						return nil
					} else {
						c.voters[msg.OwnerAddr().String()] = struct{}{}
					}
				} else {
					return nil
				}
			default:
				log.Print("`RequestAckMessage` expected, got another type")
			}

			if len(c.voters) > c.maxVotes {
				log.Println("[candidate:", c.core.Term, " -> leader:", c.core.Term, "   ]")
				return BecomeLeader(c)
			}

		default:
			return nil
		}
	}

	return nil
}
