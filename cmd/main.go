package main

import (
	"../internal/manager"
	"../internal/node"
	"../internal/router"
	"flag"
	"log"
)

func main() {
	runAll := flag.Bool("all", false, "run multiple servers using common config")
	flag.Parse()

	if *runAll == true {
		/*
		execPath := flag.Arg(0)
		if execPath == "" {
			log.Println("Error: no exec path")
		}
		configPath := flag.Arg(1)
		if configPath == "" {
			log.Println("Error: no config path")
		}
		 */
		ip := flag.Arg(0)
		if ip == "" {
			log.Println("Error: no ip")
		}
		ipPort := flag.Arg(1)
		if ipPort == "" {
			log.Println("Error: no port")
		}
		urlPort := flag.Arg(2)
		if urlPort == "" {
			log.Println("Error: no URL port")
		}
		username := flag.Arg(3)
		if urlPort == "" {
			log.Println("Error: no username")
		}
		// all run
		var raftNode = node.NewAllRunRaftCore(ip, ipPort, urlPort)

		var candidate = node.NewCandidate(raftNode)

		rm := &manager.RaftManager{
			RaftIn:  raftNode.RaftOut,
			RaftOut: raftNode.RaftIn,
			Addr:    raftNode.Addr,
		}

		cm := &manager.ClientManager{
			ClientIn:  raftNode.ClientOut,
			ClientOut: raftNode.ClientIn,
			Port: urlPort,
		}

		dbm := &manager.MongoDBManager{
			DBIn:       raftNode.DBOut,
			DBOut:      raftNode.DBIn,
			Username:   username,
		}

		r := router.NewRouterRunAll()

		go rm.ProcessMessage()
		go cm.ProcessEntries()
		go dbm.ProcessMessage()
		go r.RunRouter()

		node.RunRolePlayer(candidate)
	} else {
		// simple run
		var raftNode = node.NewRaftCore()

		var candidate = node.NewCandidate(raftNode)

		rm := &manager.RaftManager{
			RaftIn:  raftNode.RaftOut,
			RaftOut: raftNode.RaftIn,
		}

		cm := &manager.ClientManager{
			ClientIn:  raftNode.ClientOut,
			ClientOut: raftNode.ClientIn,
		}

		dbm := &manager.MongoDBManager{
			DBIn:       raftNode.DBOut,
			DBOut:      raftNode.DBIn,
		}

		r := router.NewRouter()

		go rm.ProcessMessage()
		go cm.ProcessEntries()
		go dbm.ProcessMessage()
		go r.RunRouter()

		node.RunRolePlayer(candidate)
	}
}
