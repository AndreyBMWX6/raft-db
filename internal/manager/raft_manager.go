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

	for {
		select {
		case msg := <-rm.RaftIn:
			// change Term type from int to avoid uint32 conversion()
			Term := make([]byte, 4)
			Type := make([]byte, 4)
			binary.LittleEndian.PutUint32(Term, uint32(msg.Term()))
			binary.LittleEndian.PutUint32(Type, uint32(msg.Type()))
			buff := make([]byte, 0)
			switch msg.Type() {
			case message.AppendEntriesType:
				switch entries := msg.(type) {
				case *message.AppendEntries:
					// term, type, prevterm, newidx, entries
					buff = make([]byte, 16 + len(entries.Entries))
					copy(buff, Term)
					copy(buff, Type)
					Prevterm := make([]byte, 4)
					Newidx := make([]byte, 4)
					Entries := make([]byte, len(entries.Entries))
					binary.LittleEndian.PutUint32(Prevterm, uint32(entries.PrevTerm))
					binary.LittleEndian.PutUint32(Newidx, uint32(entries.NewIndex))
					copy(buff, Prevterm)
					copy(buff, Newidx)
					copy(buff, Entries)
				default:
					log.Print("`AppendEntriesMessage` expected, got another type")
				}
			case message.RequestVoteType:
				switch requestvote := msg.(type) {
				case *message.RequestVote:
					// term, type, topindex, topterm
					buff = make([]byte, 16)
					copy(buff, Term)
					copy(buff, Type)
					Topindex := make([]byte, 4)
					Topterm := make([]byte, 4)
					binary.LittleEndian.PutUint32(Topindex, uint32(requestvote.TopIndex))
					binary.LittleEndian.PutUint32(Topterm, uint32(requestvote.TopTerm))
					copy(buff, Topindex)
					copy(buff, Topterm)

				default:
					log.Print("`RequestVoteMessage` expected, got another type")
				}
			case message.AppendAckType:
				switch appendack := msg.(type) {
				case *message.AppendAck:
					// term, type, appended
					buff = make([]byte, 12)
					copy(buff, Term)
					copy(buff, Type)
					appended_num := 0
					if appendack.Appended {
						appended_num = 1
					}
					appended := make([]byte, 4)
					binary.LittleEndian.PutUint32(appended, uint32(appended_num))
					copy(buff, appended)
				default:
					log.Print("`AppendAckMessage` expected, got another type")
				}
				case message.RequestAckType:

			}

			if _, err := conn.WriteToUDP(buff, msg.DestAddr()); err != nil {
				panic(err)
				return
			}

			default:
				recvBuff := make([]byte, 64)
				if _, ownerAddr, err := conn.ReadFromUDP(recvBuff); err == nil {
					udp_msg_term := int(binary.LittleEndian.Uint32(recvBuff[0:4]))
					udp_msg_type := int(binary.LittleEndian.Uint32(recvBuff[4:8]))

					switch udp_msg_type {
					case message.AppendEntriesType:
						Prevterm := int(binary.LittleEndian.Uint32(recvBuff[8:12]))
						Newidx := int(binary.LittleEndian.Uint32(recvBuff[12:16]))
						// implement Entries Initialization
						//Entries := recvBuff[16:]
						message.NewAppendEntries(
							&message.BaseRaftMessage{
							Owner: *ownerAddr,
							Dest: *myaddr,
							CurrTerm: udp_msg_term,
							},
							Prevterm,
							Newidx,
							//Entries,
							make([] *message.Entry, 0),
						)
					case message.RequestVoteType:
						//appended_num := int(binary.LittleEndian.Uint32(recvBuff[8:12]))

					case message.AppendAckType:

					case message.RequestAckType:

					}
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
