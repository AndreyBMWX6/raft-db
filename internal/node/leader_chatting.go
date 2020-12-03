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
			switch ack := msg.(type) {
			case *message.AppendAck:
				if ack.Appended == true {
					if ack.Heartbeat == true {
						return nil
					} else {
						l.replicated++

						// a sign of committed changes
						success := make([]*message.Entry, 1)
						success[0] = nil
						l.updates[ack.OwnerAddr().String()]<-success

						if l.replicated > ((len(l.core.Neighbors) + 1) / 2) {
							response := message.NewResponseClientMessage(
								&message.BaseClientMessage{
									Owner: nil,
									Dest:  nil,
								},
							)

							go l.core.SendClientMsg(response)
							l.replicated = 0
						}
						return nil
					}
				} else {
					retry := make([]*message.Entry, 2)
					retry[0] = nil
					retry[1] = nil
					l.updates[ack.OwnerAddr().String()]<-retry
					return nil
				}
			default:
				log.Print("`AppendAckMessage` expected, got another type")
			}
		default:
			return nil
		}
	}

	return nil
}
