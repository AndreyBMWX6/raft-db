package voter

import (
	"../../node"
	"../follower"
	"net"
	"time"
)

type Voter struct {
	core          *node.RaftCore
	candidateAddr net.Addr
}

func BecomeVoter(player node.RolePlayer, candidate net.Addr) *Voter {
	return &Voter{
		core:          player.ReleaseNode(),
		candidateAddr: candidate,
	}
}

func (v *Voter) ReleaseNode() *node.RaftCore {
	core := v.core
	v.core = nil
	return core
}

func (v *Voter) PlayRole() node.RolePlayer {
	v.MakeVote()
	time.Sleep(v.core.Config.VotingTimeout)

	return follower.BecomeFollower(v, v.candidateAddr)
}
