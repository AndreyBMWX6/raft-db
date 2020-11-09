package manager

// each separate node have own raft manager, that communicates by RaftMessages
// RaftManageres communicate betwen each other by UDP protocol

// NODE1 <----RaftMessage----> RAFTMANAGER1 <------UDP------> RAFTMANAGER2 <----RaftMessage----> NODE2

import (
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"../message"
	"../net_message"
)

type RaftManager struct {
	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage
}

func (rm *RaftManager) RaftManagerProcessMessage() {
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
			// change Term type from int to avoid int32 conversion()
			// initializing baseraftmessage
			baseraftmessage := &net_messages.BaseRaftMessage{}
			Ownerip := msg.OwnerAddr().IP
			var Ownerport = uint32(msg.OwnerAddr().Port)
			Destip := msg.DestAddr().IP
			var Destport = uint32(msg.DestAddr().Port)

			baseraftmessage.Ownerip = Ownerip[len(Ownerip)-4:]
			baseraftmessage.Ownerport = Ownerport
			baseraftmessage.Dest = Destip[len(Destip)-4:]
			baseraftmessage.Destport = Destport
			baseraftmessage.CurrTerm = msg.Term()

			switch raftmessage := msg.(type) {
			case *message.AppendEntries:
				data := &net_messages.AppendEntries{}
				// initializing data
				data.Msg = baseraftmessage
				data.PrevTerm = raftmessage.PrevTerm
				data.NewIndex = raftmessage.NewIndex
				Entries := make([]*net_messages.Entry, 0)
				for _, entrie := range raftmessage.Entries {
					Entrie := &net_messages.Entry{}
					Entrie.Term = entrie.Term
					Entrie.Query = entrie.Query
					Entries = append(Entries, Entrie)
				}
				data.Entries = Entries

				// encrypting data
				protodata, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgtype uint8 = message.AppendEntriesType
				udpmessage := make([]byte, 1)
				udpmessage[0] = msgtype
				udpmessage = append(udpmessage, protodata...)

				// sending UDP
				if _, err := conn.WriteToUDP(udpmessage, msg.DestAddr()); err != nil {
					panic(err)
					return
				}
			case *message.RequestVote:
				data := &net_messages.RequestVote{}
				// initializing data
				data.Msg = baseraftmessage
				data.TopIndex = raftmessage.TopIndex
				data.TopTerm = raftmessage.TopTerm

				// encrypting data
				protodata, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgtype uint8 = message.RequestVoteType
				udpmessage := make([]byte, 1)
				udpmessage[0] = msgtype
				udpmessage = append(udpmessage, protodata...)

				// sending UDP
				if _, err := conn.WriteToUDP(protodata, msg.DestAddr()); err != nil {
					panic(err)
					return
				}
			case *message.AppendAck:
				data := &net_messages.AppendAck{}
				// initializing data
				data.Msg = baseraftmessage
				data.Appended = raftmessage.Appended

				// encrypting data
				protodata, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgtype uint8 = message.AppendAckType
				udpmessage := make([]byte, 1)
				udpmessage[0] = msgtype
				udpmessage = append(udpmessage, protodata...)

				// sending UDP
				if _, err := conn.WriteToUDP(protodata, msg.DestAddr()); err != nil {
					panic(err)
					return
				}
			case *message.RequestAck:
				data := &net_messages.RequestAck{}
				// initializing data
				data.Msg = baseraftmessage
				data.Voted = raftmessage.Voted

				// encrypting data
				protodata, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgtype uint8 = message.RequestAckType
				udpmessage := make([]byte, 1)
				udpmessage[0] = msgtype
				udpmessage = append(udpmessage, protodata...)

				// sending UDP
				if _, err := conn.WriteToUDP(protodata, msg.DestAddr()); err != nil {
					panic(err)
					return
				}
			default:
				log.Print("unexpected type of message")
			}

		default:
			recvBuff := make([]byte, 1024)
			if length, ownerAddr, err := conn.ReadFromUDP(recvBuff); err == nil {
				data := recvBuff[1:length]
				switch recvBuff[0] {
				case uint8(message.AppendEntriesType):
					AppendEntries := net_messages.AppendEntries{}
					err := proto.Unmarshal(data, &AppendEntries)
					if err == nil {
						// converting values
						Destip := net.IPv4(
							AppendEntries.Msg.Ownerip[0],
							AppendEntries.Msg.Ownerip[1],
							AppendEntries.Msg.Ownerip[2],
							AppendEntries.Msg.Ownerip[3])
						Destudp := net.UDPAddr{
							IP:   Destip,
							Port: int(AppendEntries.Msg.Destport),
						}
						Entries := make([]*message.Entry, 0)
						for _, protoentrie := range AppendEntries.Entries {
							entrie := &message.Entry{
								protoentrie.Term,
								protoentrie.Query,
							}
							Entries = append(Entries, entrie)
						}

						message.NewAppendEntries(
							&message.BaseRaftMessage{
								Owner:    *ownerAddr,
								Dest:     Destudp,
								CurrTerm: AppendEntries.Msg.CurrTerm,
							},
							AppendEntries.PrevTerm,
							AppendEntries.NewIndex,
							Entries,
						)
					}
					if err != nil {
						log.Fatal("unmarshaling error: ", err)
					}
				case uint8(message.AppendAckType):
					AppendAck := net_messages.AppendAck{}
					err := proto.Unmarshal(data, &AppendAck)
					if err == nil {
						// converting values
						Destip := net.IPv4(
							AppendAck.Msg.Ownerip[0],
							AppendAck.Msg.Ownerip[1],
							AppendAck.Msg.Ownerip[2],
							AppendAck.Msg.Ownerip[3])
						Destudp := net.UDPAddr{
							IP:   Destip,
							Port: int(AppendAck.Msg.Destport),
						}

						message.NewEntriesAck(
							&message.BaseRaftMessage{
								Owner:    *ownerAddr,
								Dest:     Destudp,
								CurrTerm: AppendAck.Msg.CurrTerm,
							},
							AppendAck.Appended,
						)
					}
					if err != nil {
						log.Fatal("unmarshaling error: ", err)
					}
				case uint8(message.RequestVoteType):
					RequestVote := net_messages.RequestVote{}
					err := proto.Unmarshal(data, &RequestVote)
					if err == nil {
						// converting values
						Destip := net.IPv4(
							RequestVote.Msg.Ownerip[0],
							RequestVote.Msg.Ownerip[1],
							RequestVote.Msg.Ownerip[2],
							RequestVote.Msg.Ownerip[3])
						Destudp := net.UDPAddr{
							IP:   Destip,
							Port: int(RequestVote.Msg.Destport),
						}

						message.NewRequestVote(
							&message.BaseRaftMessage{
								Owner:    *ownerAddr,
								Dest:     Destudp,
								CurrTerm: RequestVote.Msg.CurrTerm,
							},
							RequestVote.TopIndex,
							RequestVote.TopTerm,
						)
					}
					if err != nil {
						log.Fatal("unmarshaling error: ", err)
					}
				case uint8(message.RequestAckType):
					RequestAck := net_messages.RequestAck{}
					err := proto.Unmarshal(data, &RequestAck)
					if err == nil {
						// converting values
						Destip := net.IPv4(
							RequestAck.Msg.Ownerip[0],
							RequestAck.Msg.Ownerip[1],
							RequestAck.Msg.Ownerip[2],
							RequestAck.Msg.Ownerip[3])
						Destudp := net.UDPAddr{
							IP:   Destip,
							Port: int(RequestAck.Msg.Destport),
						}

						message.NewRequestAck(
							&message.BaseRaftMessage{
								Owner:    *ownerAddr,
								Dest:     Destudp,
								CurrTerm: RequestAck.Msg.CurrTerm,
							},
							RequestAck.Voted,
						)
					}
					if err != nil {
						log.Fatal("unmarshaling error: ", err)
					}
				}

				if err != nil {
					panic(err)
				}
			}
		}
	}
}