package config

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/AndreyBMWX6/raft-db/internal/message"
)

type Config struct {
	FollowerTimeout  time.Duration
	VotingTimeout    time.Duration
	HeartbeatTimeout time.Duration
	DelayTime        int // in ms

	Servers   []net.UDPAddr
	URLs      []string
	Usernames []string

	// default [0, 0, 0,...]
	Terms []uint32
	// default [[],[],[],...]
	Entries [][]*message.Entry
}

type RouterConfig struct {
	URLs  []string
}

func NewConfig() *Config {
	// *SET PORTS RANGES HERE*

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

	var usernames []string
	userCnt := lIPPort - fIPPort
	for i := 1; i <= userCnt; i++ {
		usernames = append(usernames, "user" + strconv.Itoa(i))
	}

	return &Config{
		FollowerTimeout  : 4000*time.Millisecond,
		VotingTimeout    : time.Duration(rand.Intn(1000) + 1000)*time.Millisecond,
		HeartbeatTimeout : 1000*time.Millisecond,
		DelayTime        : 1,
		Servers          : servers,
		URLs			 : urls,
		Usernames        : usernames,
		Terms            : make([]uint32, len(servers)),
		Entries          : make([][]*message.Entry, len(servers)),
	}
}

func (cfg *Config) Delay(time int) {
	for i := 0; i < time; i++ {
		// code for delay
		m := make(map[int]int)
		for i := 0; i < 10000; i++ {
			m[i] = i
		}
	}
}

func NewRouterConfig() *RouterConfig {
	cfg := NewConfig()

	return &RouterConfig{
		URLs:  cfg.URLs,
	}
}
