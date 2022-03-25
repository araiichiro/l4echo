package main

import (
	"fmt"
	"time"

	"github.com/araiichiro/l4echo/internal/client"
	"github.com/araiichiro/l4echo/internal/server"
)

func main() {
	switch env.Mode {
	case "client":
		cs := []client.Config{
			{
				Concurrency: env.Concurrency,
				Network:     "udp",
				Address:     env.ClientEnv.UDPAddr,
				Workload: client.Workload{
					Loop:     env.Loop,
					Interval: env.Interval,
				},
			},
			{
				Concurrency: env.Concurrency,
				Network:     "tcp",
				Address:     env.ClientEnv.TCPAddr,
				RecvTimeout: 10 * time.Second,
				SendTimeout: 1 * time.Second,
				Workload: client.Workload{
					Loop:     env.Loop,
					Interval: env.Interval,
				},
			},
		}
		client.Run(cs)
	case "server":
		s := server.Server{
			TCPAddr: env.ServerEnv.TCPAddr,
			UDPAddr: env.ServerEnv.UDPAddr,
		}
		s.Serve()
	default:
		panic(fmt.Errorf("invalid mode: %s", env.Mode))
	}
}
