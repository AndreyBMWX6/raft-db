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
		switch msg.Type() {
		case message.AppendEntriesType: // may be should add processing of query
			if msg.Term() > l.core.Term {
				log.Println("[leader:", l.core.Term, "   -> follower:", msg.Term(), " ]")
				l.core.Term = msg.Term()
				return BecomeFollower(l, msg.OwnerAddr())
			} else {
				switch entries := msg.(type) {
				case *message.AppendEntries:
					entrTopTerm := entries.Entries[entries.NewIndex - 1].Term
					coreTopTerm := l.core.Entries[len(l.core.Entries)].Term
					if entrTopTerm > coreTopTerm {
						return BecomeFollower(l, msg.OwnerAddr())
					} else if entrTopTerm == coreTopTerm && len(entries.Entries) > len(l.core.Entries) {
						return BecomeFollower(l, msg.OwnerAddr())
					}
					return nil

				default:
					log.Print("`AppendEntriesMessage` expected, got another type")
				}
			}
		case message.RequestVoteType:
			if msg.Term() > l.core.Term {
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

				log.Println("[leader:", l.core.Term, "   -> follower:", msg.Term(), " ]")
				l.core.Term = msg.Term()
				l.core.ProcessRequestVote(request)
				return BecomeFollower(l, msg.OwnerAddr())
				default:
					log.Print("`RequestVoteMessage` expected, got another type")
				}
			}
		default:
			return nil
		}
	}

	return nil
}
