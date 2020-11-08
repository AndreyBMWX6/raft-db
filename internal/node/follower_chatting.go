package node

import (
	"fmt"
	"log"

	"../message"
)

func (f *Follower) ApplyRaftMessage(msg message.RaftMessage) RolePlayer {
	if msg.Term() < f.core.Term {
		return nil
	}

	if msg.Term() >= f.core.Term {
		f.core.Term = msg.Term()
		switch msg.Type() {
		case message.AppendEntriesType:
			// logging
			fmt.Print("follower: ")
			return BecomeFollower(f, msg.OwnerAddr())
		case message.RequestVoteType:
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
				f.core.ProcessRequestVote(request)
				return BecomeFollower(f, msg.OwnerAddr())
			default:
				log.Print("`RequestVoteMessage` expected, got another type")
			}
		default:
			return nil
		}
	}

	switch msg.Type() {
	case message.AppendEntriesType:

		if msg.OwnerAddr().String() != f.leaderAddr.String() {
			return nil
		}

		switch entries := msg.(type) {
		case *message.AppendEntries:
			f.ApplyAppendEntries(entries)
			f.timer.Reset(f.core.Config.FollowerTimeout)
		default:
			log.Print("`AppendEntriesMessage` expected, got another type")
		}
	}

	return nil
}

func (f *Follower) ApplyAppendEntries(entries *message.AppendEntries) {
	ack := message.NewEntriesAck(
		&message.BaseRaftMessage{
			Owner:    f.core.Addr,
			Dest:     f.leaderAddr,
			CurrTerm: f.core.Term,
		},
	)

	if entries.NewIndex > uint32(len(f.core.Entries)) {
		ack.Appended = false
	} else {
		var prevTerm = f.core.Entries[entries.NewIndex-1].Term
		if entries.PrevTerm != prevTerm {
			ack.Appended = false
		} else {
			// if nil, got heartbeat, so don't need to append
			if entries.Entries != nil {
				f.core.Entries = append(f.core.Entries[:entries.NewIndex],
					entries.Entries...)
			}
			ack.Appended = true
		}
	}

	f.core.SendRaftMsg(
		message.RaftMessage(ack),
	)
}

func (f *Follower) ApplyClientMessage(msg message.ClientMessage) {
	f.core.SendClientMsg(msg)
}
