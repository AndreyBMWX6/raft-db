package manager

// each separate node have own raft manager, that communicates by RaftMessages
// RaftManageres communicate betwen each other by UDP protocol

// NODE1 <----RaftMessage----> RAFTMANAGER1 <------UDP------> RAFTMANAGER2 <----RaftMessage----> NODE2

import (
	"../message"
	"encoding/binary"
	"log"
	"net"
)

type RaftManager struct {
	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage
}

func (rm *RaftManager) RaftManagerProcessMessage(msg message.RaftMessage) {
	// Resolving address
	myaddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:800")
	if err != nil {
		log.Fatal(err)
	}

	// Build listening connections
	conn, err := net.ListenUDP("udp", myaddr)
	defer conn.Close()
	if err != nil {
		log.Println("Error: ", err)
	}

	recvBuff := make([]byte, 64)

	for {
		select {
		case msg := <-rm.RaftIn:
			// change Term type from int to avoid uint32 conversion()
			Term := make([]byte, 4)
			Type := make([]byte, 4)
			binary.LittleEndian.PutUint32(Term, uint32(msg.Term()))
			binary.LittleEndian.PutUint32(Type, uint32(msg.Type()))
			switch msg.Type() {
			case message.AppendEntriesType:

			case message.RequestVoteType:

			case message.AppendAckType:

			case message.VoteType:

			}
			buff := make([]byte, 8)
			copy(buff, Term)
			copy(buff, Type)

			if _, err := conn.WriteToUDP(buff, msg.DestAddr()); err != nil {
				panic(err)
				return
			}

			default:
				if _, ownerAddr, err := conn.ReadFromUDP(recvBuff); err == nil {
					udp_msg_term := int(binary.LittleEndian.Uint32(recvBuff[:4]))
					udp_msg_type := int(binary.LittleEndian.Uint32(recvBuff[4:]))

					switch udp_msg_type {
					case message.AppendEntriesType:
						message.NewAppendEntries(
							&message.BaseRaftMessage{
							Owner: *ownerAddr,
							Dest: *myaddr,
							CurrTerm: udp_msg_term,
							},
							3, // get value from udp data
							3, // get value from udp data
							make([]*message.Entry, 0), // get value from udp data
						)
					case message.RequestVoteType:

					case message.AppendAckType:

					case message.VoteType:

					}
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
