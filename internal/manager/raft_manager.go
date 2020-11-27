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
			// initializing baseraftmessage
			baseRaftMsg := &net_messages.BaseRaftMessage{}
			ownerIp := msg.OwnerAddr().IP
			var ownerPort = uint32(msg.OwnerAddr().Port)
			destIp := msg.DestAddr().IP
			var destPort = uint32(msg.DestAddr().Port)

			baseRaftMsg.Ownerip = ownerIp[len(ownerIp)-4:]
			baseRaftMsg.Ownerport = ownerPort
			baseRaftMsg.Dest = destIp[len(destIp)-4:]
			baseRaftMsg.Destport = destPort
			baseRaftMsg.CurrTerm = msg.Term()

			switch raftMsg := msg.(type) {
			case *message.AppendEntries:
				data := &net_messages.AppendEntries{}
				// initializing data
				data.Msg = baseRaftMsg
				data.PrevTerm = raftMsg.PrevTerm
				data.NewIndex = raftMsg.NewIndex
				entries := make([]*net_messages.Entry, 0)
				for _, entry := range raftMsg.Entries {
					Entry := &net_messages.Entry{}
					Entry.Term = entry.Term
					Entry.Query = entry.Query
					entries = append(entries, Entry)
				}
				data.Entries = entries

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgType uint8 = message.AppendEntriesType
				udpMsg := make([]byte, 1)
				udpMsg[0] = msgType
				udpMsg = append(udpMsg, protoData...)

				// sending UDP
				if _, err := conn.WriteToUDP(udpMsg, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.RequestVote:
				data := &net_messages.RequestVote{}
				// initializing data
				data.Msg = baseRaftMsg
				data.TopIndex = raftMsg.TopIndex
				data.TopTerm = raftMsg.TopTerm

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgType uint8 = message.RequestVoteType
				udpMsg := make([]byte, 1)
				udpMsg[0] = msgType
				udpMsg = append(udpMsg, protoData...)

				// sending UDP
				if _, err := conn.WriteToUDP(udpMsg, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.AppendAck:
				data := &net_messages.AppendAck{}
				// initializing data
				data.Msg = baseRaftMsg
				data.Appended = raftMsg.Appended

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgType uint8 = message.AppendAckType
				udpMsg := make([]byte, 1)
				udpMsg[0] = msgType
				udpMsg = append(udpMsg, protoData...)

				// sending UDP
				if _, err := conn.WriteToUDP(udpMsg, msg.DestAddr()); err != nil {
					log.Fatal(err)
					return
				}
			case *message.RequestAck:
				data := &net_messages.RequestAck{}
				// initializing data
				data.Msg = baseRaftMsg
				data.Voted = raftMsg.Voted

				// encrypting data
				protoData, err := proto.Marshal(data)
				if err != nil {
					log.Fatal("marshaling error: ", err)
					return
				}

				// adding message type to message
				var msgType uint8 = message.RequestAckType
				udpMsg := make([]byte, 1)
				udpMsg[0] = msgType
				udpMsg = append(udpMsg, protoData...)

				// sending UDP
				if _, err := conn.WriteToUDP(udpMsg, msg.DestAddr()); err != nil {
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
			data := recvBuff[1:length]
			switch recvBuff[0] {
			case uint8(message.AppendEntriesType):
				var appendEntries *message.AppendEntries
				rm.RaftOut <- appendEntries.Unmarshal(data)
			case uint8(message.AppendAckType):
				var appendAck *message.AppendAck
				rm.RaftOut <- appendAck.Unmarshal(data)
			case uint8(message.RequestVoteType):
				var requestVote *message.RequestVote
				rm.RaftOut <- requestVote.Unmarshal(data)
			case uint8(message.RequestAckType):
				var requestAck *message.RequestAck
				rm.RaftOut <- requestAck.Unmarshal(data)
			default:

			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
