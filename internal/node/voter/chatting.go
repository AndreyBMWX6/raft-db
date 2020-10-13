package voter

import (
	"../../message"
)

func (v *Voter) MakeVote() {
	v.core.SendRaftMsg(
		message.NewVote(
			&message.BaseRaftMessage{
				Owner:    v.core.Addr,
				Dest:     v.candidateAddr,
				CurrTerm: v.core.Term,
			},
		),
	)
}
