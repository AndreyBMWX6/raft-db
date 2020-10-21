package node

import (
	"../message"
	"fmt"
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
				request := message.NewRequestVote(
					&message.BaseRaftMessage{
						Owner:	  msg.OwnerAddr(),
						Dest: 	  msg.DestAddr(),
						CurrTerm: msg.Term(),
					},
				)
				l.core.ProcessRequestVote(request)
				return BecomeFollower(l, msg.OwnerAddr())
			}
		default:
			return nil
		}
	}

	return nil
}

