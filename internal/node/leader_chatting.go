package node

import (
	"../message"
	"fmt"
	"log"
)

func (l *Leader) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < l.core.Term {
		return nil
	}

	if msg.Term() >= l.core.Term {
		switch msg.Type() {
		case message.AppendEntriesType: // may be should add processing of query
			if msg.Term() > l.core.Term {
				l.core.Term = msg.Term()
				// logging
				fmt.Print("leader: ")
				return BecomeFollower(l, msg.OwnerAddr())
			}
		case message.RequestVoteType:

			if msg.Term() > l.core.Term {
				l.core.Term = msg.Term()
				switch requestvote := msg.(type) {
				case *message.RequestVote:
					request := message.NewRequestVote(
					&message.BaseRaftMessage{
						Owner:    *msg.OwnerAddr(),
						Dest:     *msg.DestAddr(),
						CurrTerm: msg.Term(),
					},
					requestvote.TopIndex,
					requestvote.TopTerm,
				)
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

