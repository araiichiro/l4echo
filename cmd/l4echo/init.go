package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	env = &Env{
		Mode: "",
		ClientEnv: ClientEnv{
			Concurrency: 1024,
			Loop:        500,
			Interval:    80 * time.Millisecond,
			TCPAddr:     "localhost:7000",
			UDPAddr:     "localhost:7001",
		},
		ServerEnv: ServerEnv{
			TCPAddr: ":7000",
			UDPAddr: ":7001",
		},
	}
)

func init() {
	if err := envconfig.Process("", env); err != nil {
		panic(err)
	}
}

type Env struct {
	Mode string
	ClientEnv
	ServerEnv
}

type ClientEnv struct {
	Concurrency int
	Loop        int
	Interval    time.Duration
	TCPAddr     string
	UDPAddr     string
}

type ServerEnv struct {
	TCPAddr string
	UDPAddr string
}
