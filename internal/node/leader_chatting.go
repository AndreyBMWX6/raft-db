package node

import (
	"log"

	"../message"
)

func (l *Leader) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < l.core.Term {
		return nil
	}

	if msg.Term() >= l.core.Term {
		oldLeaderTerm := l.core.Term
		l.core.Term = msg.Term()
		switch msg.Type() {
		case message.AppendEntriesType: // may be should add processing of query
			if msg.Term() > oldLeaderTerm {
				log.Println("[leader:", l.core.Term, "   -> follower:", msg.Term(), " ]")
				return BecomeFollower(l)
			}
		case message.RequestVoteType:
			if msg.Term() >= l.core.Term {
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
				l.core.ProcessRequestVote(request)

				if msg.Term() > oldLeaderTerm {
					log.Println("[leader:", oldLeaderTerm, "   -> follower:", msg.Term(), " ]")
					return BecomeFollower(l)
				} else {
					return nil
				}

				default:
					log.Print("`RequestVoteMessage` expected, got another type")
				}
			}
		case message.AppendAckType:
			// make logic to send AppendEntries when needed and when not - only heartbeat
			return nil
		default:
			return nil
		}
	}

	return nil
}
