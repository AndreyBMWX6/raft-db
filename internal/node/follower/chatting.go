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

	if msg.Term() >= f.core.Term {
		f.core.Term = msg.Term()
		switch msg.Type() {
		case message.AppendEntriesType:
			return BecomeFollower(f, msg.OwnerAddr())
		case message.RequestVoteType:
			request := message.NewRequestVote(
				&message.BaseRaftMessage{
					Owner:	  msg.OwnerAddr(),
					Dest: 	  msg.DestAddr(),
					CurrTerm: msg.Term(),
				},
			)
			// msg can be passed to function to remove if
			response := f.ProcessRequestVote(request)
			if response { // to return in ProcessQuery
				return voter.BecomeVoter(f, msg.OwnerAddr())
			}
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

func (f *Follower) ProcessRequestVote(request *message.RequestVote) bool {
	ack := message.NewRequestAck(
		&message.BaseRaftMessage{
			Owner:    f.core.Addr,
			Dest:     request.Owner,
			CurrTerm: f.core.Term,
		},
	)
	topterm := f.core.Entries[len(f.core.Entries) - 1].Term
	if request.TopTerm < topterm {
		return false
	} else {
		var topindex = len(f.core.Entries)
		if (request.TopTerm == topterm) && (request.TopIndex < topindex) {
			return false
			// not sending request_ack may be simplified
		} else {
			// sending request_ack signalize that node voted for candidate
			f.core.SendRaftMsg(
				message.RaftMessage(ack),
			)
			return true
		}
	}
}

func (f *Follower) ApplyClientMessage(msg message.ClientMessage) {
	f.core.SendClientMsg(msg)
}
