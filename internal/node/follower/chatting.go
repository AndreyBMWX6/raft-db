package follower

import (
	"log"

	"../../message"
	"../../node"
	"../voter"
)

func (f *Follower) ApplyRaftMessage(msg message.RaftMessage) node.RolePlayer {
	if msg.Term() < f.core.Term {
		return nil
	}

	if msg.Term() > f.core.Term {
		f.core.Term = msg.Term()
		switch msg.Type() {
		case message.AppendEntriesType:
			return BecomeFollower(f, msg.OwnerAddr())
		case message.RequestVoteType:
			return voter.BecomeVoter(f, msg.OwnerAddr())
		default:
			return nil
		}
	}

	switch msg.Type() {
	case message.AppendEntriesType:
		if msg.OwnerAddr() != f.leaderAddr {
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

	if entries.NewIndex > len(f.core.Entries) {
		ack.Appended = false
	} else {
		var prevTerm = f.core.Entries[entries.NewIndex-1].Term
		if entries.PrevTerm != prevTerm {
			ack.Appended = false
		} else {
			f.core.Entries = append(f.core.Entries[:entries.NewIndex],
				entries.Entries...)
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
