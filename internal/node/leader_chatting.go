package node

import (
	"log"

	"github.com/AndreyBMWX6/raft-db/internal/message"
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
						l.replicated[ack.TopIndex]++

						// a sign of committed changes
						success := make([]*message.Entry, 1)
						success[0] = nil
						l.updates[ack.OwnerAddr().String()]<-success

						if l.replicated[ack.TopIndex] > ((len(l.core.Neighbors) + 1) / 2) {
							response := message.NewResponseClientMessage(
								&message.BaseClientMessage{
									Owner: nil,
									Dest:  nil,
								},
								l.responses[ack.TopIndex],
								false,
							)

							if l.replicated[ack.TopIndex] > ((len(l.core.Neighbors) + 1) / 2) + 1 {
								return nil
							} else {
								go l.core.SendClientMsg(response)
							}
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

func (l *Leader) ApplyClientMessage(msg message.ClientMessage) {
	switch rawClient := msg.(type) {
	case *message.RawClientMessage:
		l.core.SendDBMsg(rawClient.DBReq)
	default:
		log.Print("`RawClientMessage` expected, got another type")
	}
}

func (l *Leader) ApplyDBMessage(msg message.DBMessage) {
	switch dbResp := msg.(type) {
	case *message.DBResponse:
		if dbResp.Entry != nil {
			dbResp.Entry.Term = l.core.Term
			l.core.Entries = append(l.core.Entries, dbResp.Entry)
			l.replicated = append(l.replicated, 1)
			l.responses = append(l.responses, dbResp)

			log.Println("Leader added new entry")
			log.Println("Leader log:     ", l.core.Entries)
			var entriesTerms []uint32
			for _,entry := range l.core.Entries {
				entriesTerms = append(entriesTerms, entry.Term)
			}
			log.Println("Log terms:      ", entriesTerms)

			var entries []*message.Entry
			entries = append(entries, dbResp.Entry)
			log.Println("Append entries: ", entries)
			entriesTerms = nil
			for _,entry := range entries {
				entriesTerms = append(entriesTerms, entry.Term)
			}
			log.Println("Entries terms:  ", entriesTerms)

			for _, update := range l.updates {
				upd := make([]*message.Entry, 1)
				upd[0] = dbResp.Entry
				update <-upd
			}
		} else {
			response := message.NewResponseClientMessage(
				&message.BaseClientMessage{
					Owner: nil,
					Dest:  nil,
				},
				dbResp,
				false,
			)
			go l.core.SendClientMsg(response)
		}

	default:
		log.Print("`RawClientMessage` expected, got another type")
	}
}
