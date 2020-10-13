package config

import "time"

type Config struct {
	FollowerTimeout  time.Duration
	VotingTimeout    time.Duration
	HeartbeatTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		FollowerTimeout:  100*time.Millisecond,
		VotingTimeout:    100*time.Millisecond,
		HeartbeatTimeout: 100*time.Millisecond,
	}
}
