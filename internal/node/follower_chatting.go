package node

import (
	"log"

	"../message"
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
				f.ApplyAppendEntries(entries)
				//f.timer.Reset(f.core.Config.FollowerTimeout)
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
				f.core.ProcessRequestVote(request)
				
				if msg.Term() > oldFollowerTerm {
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
	ack := message.NewEntriesAck(
		&message.BaseRaftMessage{
			Owner:    f.core.Addr,
			Dest:     entries.Owner,
			CurrTerm: f.core.Term,
		},
		false,
	)

	if entries.Entries == nil {
		ack.Appended = true
	} else {
		log.Println("follower got updates from leader")
		if entries.NewIndex < uint32(len(f.core.Entries)) {
			ack.Appended = false
		} else {
			var prevTerm uint32 = 0
			if len(f.core.Entries) != 0 && entries.NewIndex > 0 {
				prevTerm = f.core.Entries[entries.NewIndex-1].Term
			}
			if entries.PrevTerm != prevTerm {
				ack.Appended = false
			} else {
					f.core.Entries = append(f.core.Entries[:entries.NewIndex],
						entries.Entries...)
				ack.Appended = true
				log.Println("follower committed updates")
			}
		}
	}

	log.Println("Node:", ack.Owner.String(), " send EntriesAck:", ack.CurrTerm,
		" to Node:", ack.Dest.String())
	f.core.SendRaftMsg(
		message.RaftMessage(ack),
	)
}

func (f *Follower) ApplyClientMessage(msg message.ClientMessage) {
	f.core.SendClientMsg(msg)
}
