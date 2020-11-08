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
			// change Term type from int to avoid int32 conversion()
			// initializing baseraftmessage
			baseraftmessage := &net_messages.BaseRaftMessage{}
			baseraftmessage.Owner = msg.OwnerAddr().String()
			baseraftmessage.Dest = msg.DestAddr().String()
			baseraftmessage.CurrTerm = int32(msg.Term())

			switch raftmessage := msg.(type) {
			case *message.AppendEntries:
				data := &net_messages.AppendEntries{}
				// initializing data
				data.Msg = baseraftmessage
				data.PrevTerm = int32(raftmessage.PrevTerm)
				data.NewIndex = int32(raftmessage.NewIndex)
				Entries := make([]*net_messages.Entry, 0)
				for _,entrie := range raftmessage.Entries {
					Entrie := &net_messages.Entry{}
					Entrie.Term = int32(entrie.Term)
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

				// sending UDP
				if _, err := conn.WriteToUDP(protodata, msg.DestAddr()); err != nil {
					panic(err)
					return
				}
			case *message.RequestVote:
				data := &net_messages.RequestVote{}
				// initializing data
				data.Msg = baseraftmessage
				data.TopIndex = int32(raftmessage.TopIndex)
				data.TopTerm = int32(raftmessage.TopTerm)

				// encrypting data
				protodata, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

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

					// sending UDP
					if _, err := conn.WriteToUDP(protodata, msg.DestAddr()); err != nil {
						panic(err)
						return
					}
			default:
				log.Print("unexpected type of message")
				}

		default:


			/*
				AppendEntries := net_messages.AppendEntries{}
				AppendAck := net_messages.AppendAck{}
				RequestVote := net_messages.RequestVote{}
				RequestAck := net_messages.RequestAck{}

				recvBuff := make([]byte, 1024)
				if length, ownerAddr, err := conn.ReadFromUDP(recvBuff); err == nil {
					data := recvBuff[:length]

					err := proto.Unmarshal(data, &AppendEntries)
					if err == nil {
						message.NewAppendEntries(
							&message.BaseRaftMessage{
								Owner: ,
								Dest: neighbor,
								CurrTerm: l.core.Term,
							},
							prevTerm,
							len(l.core.Entries),
							make([]*message.Entry, 0),
						)
					}


					err = proto.Unmarshal(data, &AppendAck)
					if err == nil {

					}

					err = proto.Unmarshal(data, &RequestVote)
					if err == nil {

					}

					err = proto.Unmarshal(data, &RequestAck)
					if err == nil {

					}
			*/
			if err != nil {
				panic(err)
			}
		}
	}
}
