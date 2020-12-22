package node

import (
	"log"

	"github.com/AndreyBMWX6/raft-db/internal/message"
)

func (f *Follower) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < f.core.Term {
		return nil
	}

	if msg.Term() >= f.core.Term {
		oldFollowerTerm := f.core.Term
		f.core.Term = msg.Term()
		switch msg.Type() {
		case message.AppendEntriesType:
			switch entries := msg.(type) {
			case *message.AppendEntries:
				if entries.LeaderURL != "" {
					f.leaderURL = entries.LeaderURL
				}
				f.ApplyAppendEntries(entries)
			default:
				log.Print("`AppendEntriesMessage` expected, got another type")
			}
			return RefreshFollower(f)
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

				if msg.Term() > oldFollowerTerm {
					f.core.Voted = false
				} else {
					f.core.Voted = true
				}
				fVoted := f.core.ProcessRequestVote(request)
				
				if fVoted == true {
					log.Println("[follower:", oldFollowerTerm, "  -> follower:", msg.Term(), " ]")
					return BecomeFollower(f)
				} else {
					return nil
				}
			default:
				log.Print("`RequestVoteMessage` expected, got another type")
			}
		default:
			return nil
		}
	}


	return nil
}

func (f *Follower) ApplyAppendEntries(entries *message.AppendEntries) {
	ack := message.NewAppendAck(
		&message.BaseRaftMessage{
			Owner:    f.core.Addr,
			Dest:     entries.Owner,
			CurrTerm: f.core.Term,
		},
		false,
		false,
		0,
	)

	if entries.Entries == nil {
		ack.Appended = true
		ack.Heartbeat = true
	} else {
		log.Println("Append entries:", entries.Entries)
		var entriesTerms []uint32
		for _,entry := range entries.Entries {
			entriesTerms = append(entriesTerms, entry.Term)
		}
		log.Println("Entries terms: ", entriesTerms)

		// metadata check
		if entries.NewIndex > uint32(len(f.core.Entries)) {
			ack.Appended = false
			log.Println("Failed to add new entry - Metadata checks error: " +
				"Follower log length: < NewIndex in AppendEntries    ",
				uint32(len(f.core.Entries)), " < ", entries.NewIndex)
		} else {
			var prevTerm uint32 = 0
			if len(f.core.Entries) != 0 && entries.NewIndex > 0 {
				prevTerm = f.core.Entries[entries.NewIndex-1].Term
			}

			if len(f.core.Entries) == 0 {
				if uint32(cap(f.core.Entries)) >= entries.NewIndex {
					f.core.Entries = append(f.core.Entries[:entries.NewIndex],
						entries.Entries...)
				} else {
					f.core.Entries = append(f.core.Entries, entries.Entries...)
				}
				ack.Appended = true
				ack.TopIndex = uint32(len(f.core.Entries) - 1)
				log.Println("New entry added successfully")
			} else {
				if entries.PrevTerm != prevTerm {
					ack.Appended = false
					log.Println("Failed to add new entry - Metadata checks error: " +
						"Follower entry's at index", entries.NewIndex - 1, "term", " != PrevTerm in AppendEntries",
						prevTerm, "!=", entries.PrevTerm)
				} else {
					if uint32(cap(f.core.Entries)) >= entries.NewIndex {
						f.core.Entries = append(f.core.Entries[:entries.NewIndex],
							entries.Entries...)
					} else {
						f.core.Entries = append(f.core.Entries, entries.Entries...)
					}
					ack.Appended = true
					ack.TopIndex = uint32(len(f.core.Entries) - 1)
					log.Println("New entry added successfully")
				}
			}
		}

		log.Println("Follower log:  ", f.core.Entries)
		entriesTerms = nil
		for _,entry := range f.core.Entries {
			entriesTerms = append(entriesTerms, entry.Term)
		}
		log.Println("Log terms:     ", entriesTerms)
	}

	log.Println("Node:", ack.Owner.String(), " send ", ack.Appended, "AppendAck:", ack.CurrTerm,
		" to Node:", ack.Dest.String())

	f.core.SendRaftMsg(
		message.RaftMessage(ack),
	)
}

func (f *Follower) ApplyClientMessage(msg message.ClientMessage) {
	switch clientMsg := msg.(type) {
	case *message.RawClientMessage:
		if clientMsg.ReqType == message.GetRequestType {
			go f.core.SendDBMsg(clientMsg.DBReq)
		} else {
			response := message.NewResponseClientMessage(
				&message.BaseClientMessage{
					Owner: nil,
					Dest:  nil,
				},
				nil,
				true,
			)

			response.LeaderURL = f.leaderURL

			log.Println("Node:", f.core.Addr.String(), " redirected client message to leader" +
				"(URL:", f.leaderURL, ")")
			go f.core.SendClientMsg(response)
		}
	default:
		log.Print("`RawClientMessage` expected, got another type")
	}
}

func (f *Follower) ApplyDBMessage(msg message.DBMessage) {
	switch dbResp := msg.(type) {
	case *message.DBResponse:
		response := message.NewResponseClientMessage(
			&message.BaseClientMessage{
				Owner: nil,
				Dest:  nil,
			},
			dbResp,
			false,
		)
		go f.core.SendClientMsg(response)
	default:
		log.Print("`RawClientMessage` expected, got another type")
	}
}
