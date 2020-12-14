package config

import (
	"../message"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type AllConfig struct {
	FollowerTimeout  time.Duration
	VotingTimeout    time.Duration
	HeartbeatTimeout time.Duration

	Servers []net.UDPAddr
	URLs    []string

	// default [0, 0, 0,...]
	Terms []uint32
	// default [[],[],[],...]
	Entries [][]*message.Entry
}

func NewAllConfig() *AllConfig {
	// ip ports range
	firstIPPort := "8001"
	lastIPPort  := "8006"

	// URL ports range
	firstURLPort := "8081"
	lastURLPort := "8086"

	fIPPort, err := strconv.Atoi(firstIPPort)
	if err != nil {
		log.Fatal(err)
	}
	lIPPort, err := strconv.Atoi(lastIPPort)
	if err != nil {
		log.Fatal(err)
	}

	fURLPort, err := strconv.Atoi(firstURLPort)
	if err != nil {
		log.Fatal(err)
	}
	lURLPort, err := strconv.Atoi(lastURLPort)
	if err != nil {
		log.Fatal(err)
	}

	var serversStr []string
	serversStr = append(serversStr, ("127.0.0.1:" + firstIPPort))

	var urls []string
	urls = append(urls, ("http://localhost:" + firstURLPort))

	for fIPPort != lIPPort && fURLPort != lURLPort {
		fIPPort++
		fURLPort++
		serversStr = append(serversStr, "127.0.0.1:" + strconv.Itoa(fIPPort))
		urls = append(urls, "http://localhost:" + strconv.Itoa(fURLPort))
	}

	var servers []net.UDPAddr
	for _, server := range serversStr {
		addr, err := net.ResolveUDPAddr("udp4", server)
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, *addr)
	}

	return &AllConfig{
		FollowerTimeout  : 4000*time.Millisecond,
		VotingTimeout    : time.Duration(rand.Intn(1000) + 1000)*time.Millisecond,
		HeartbeatTimeout : 1000*time.Millisecond,
		Servers          : servers,
		URLs			 : urls,
		Terms            : make([]uint32, len(servers)),
		Entries          : make([][]*message.Entry, len(servers)),
	}
}

func NewRouterRunAllConfig() *RouterConfig {
allCfg := NewAllConfig()

return &RouterConfig{
		URLs:  allCfg.URLs,
	}
}

