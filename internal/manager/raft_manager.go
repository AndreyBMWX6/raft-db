package manager

// each separate node have own raft manager, that communicates by RaftMessages
// RaftManagers communicate between each other by UDP protocol

// NODE1 <----RaftMessage----> RAFTMANAGER1 <------UDP------> RAFTMANAGER2 <----RaftMessage----> NODE2

import (
	"google.golang.org/protobuf/proto"
	"log"
	"net"

	"../config"
	"../message"
	"../net_message"
)

type RaftManager struct {
	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage
}

func (rm *RaftManager) ProcessMessage() {
	// Resolving address
	myAddr := config.NewConfig().Addr

	// Build listening connections
	conn, err := net.ListenUDP("udp", &myAddr)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	go rm.ListenToUDP(conn)

	for {
		select {
		case msg := <-rm.RaftIn:
			// change Term type from int to avoid int32 conversion()
			// initializing BaseRaftMessage
			baseRaftMsg := &net_message.BaseRaftMessage{}
			ownerIp := msg.OwnerAddr().IP
			var ownerPort = uint32(msg.OwnerAddr().Port)
			destIp := msg.DestAddr().IP
			var destPort = uint32(msg.DestAddr().Port)

			baseRaftMsg.OwnerIp = ownerIp[len(ownerIp)-4:]
			baseRaftMsg.OwnerPort = ownerPort
			baseRaftMsg.DestIp = destIp[len(destIp)-4:]
			baseRaftMsg.DestPort = destPort
			baseRaftMsg.CurrTerm = msg.Term()

			switch raftMsg := msg.(type) {
			case *message.AppendEntries:
				// initializing entries
				entries := make([]*net_message.Entry, 0)
				if raftMsg.Entries == nil {
					entries = nil
				} else {
					for _, entry := range raftMsg.Entries {
						Entry := &net_message.Entry{}
						Entry.Term = entry.Term
						Entry.Query = entry.Query
						entries = append(entries, Entry)
					}
				}

				// initializing data
				data := &net_message.Message{
				RaftMessage: &net_message.Message_AppendEntries{
					AppendEntries: &net_message.AppendEntries{
						Msg:      baseRaftMsg,
						PrevTerm: raftMsg.PrevTerm,
						NewIndex: raftMsg.NewIndex,
						Entries:  entries,
						},
					},
				}

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// sending UDP
				if _, err := conn.WriteToUDP(protoData, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.RequestVote:
				// initializing data
				data := &net_message.Message{
					RaftMessage: &net_message.Message_RequestVote{
						RequestVote: &net_message.RequestVote{
							Msg:      baseRaftMsg,
							TopIndex: raftMsg.TopIndex,
							TopTerm:  raftMsg.TopTerm,
						},
					},
				}

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// sending UDP
				if _, err := conn.WriteToUDP(protoData, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.AppendAck:
				// initializing data
				data := &net_message.Message{
					RaftMessage: &net_message.Message_AppendAck{
						AppendAck: &net_message.AppendAck{
							Msg:       baseRaftMsg,
							Appended:  raftMsg.Appended,
							Heartbeat: raftMsg.Heartbeat,
						},
					},
				}

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// sending UDP
				if _, err := conn.WriteToUDP(protoData, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.RequestAck:
				// initializing data
				data := &net_message.Message{
					RaftMessage: &net_message.Message_RequestAck{
						RequestAck: &net_message.RequestAck{
							Msg:   baseRaftMsg,
							Voted: raftMsg.Voted,
						},
					},
				}

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// sending UDP
				if _, err := conn.WriteToUDP(protoData, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			default:
				log.Print("unexpected type of message")
			}
		default:
		}
	}
}

func (rm *RaftManager) ListenToUDP(conn *net.UDPConn) {
	recvBuff := make([]byte, 1024)
	for {
		if length, _, err := conn.ReadFromUDP(recvBuff); err == nil {
			msg := net_message.Message{}
			pErr := proto.Unmarshal(recvBuff[:length], &msg)
			if pErr == nil {
				switch msg.RaftMessage.(type) {
				case *net_message.Message_AppendEntries:
					var appendEntries *message.AppendEntries
					rm.RaftOut <- appendEntries.Unmarshal(&msg)
				case *net_message.Message_AppendAck:
					var appendAck *message.AppendAck
					rm.RaftOut <- appendAck.Unmarshal(&msg)
				case *net_message.Message_RequestVote:
					var requestVote *message.RequestVote
					rm.RaftOut <- requestVote.Unmarshal(&msg)
				case *net_message.Message_RequestAck:
					var requestAck *message.RequestAck
					rm.RaftOut <- requestAck.Unmarshal(&msg)
				}
			}
			if pErr != nil {
				log.Println(recvBuff)
				log.Fatal("unmarshalling error: ", err)
			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
