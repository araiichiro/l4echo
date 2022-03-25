package client

import (
	"net"
	"sync"
	"time"

	"github.com/araiichiro/l4echo/internal/log"
	"github.com/araiichiro/l4echo/internal/network"
)

func Run(configs []Config) {
	wg := &sync.WaitGroup{}
	statsMap := map[string]*Stats{}
	for _, config := range configs {
		statsMap[config.Network] = &Stats{}

		for i := 0; i < config.Concurrency; i++ {
			wg.Add(1)

			conn, err := net.Dial(config.Network, config.Address)
			if err != nil {
				log.Error("failed to dial:", err)
				return
			}
			if config.Network == "tcp" {
				conn = &network.ConnWithTimeout{
					Conn:         conn,
					ReadTimeout:  config.RecvTimeout,
					WriteTimeout: config.SendTimeout,
				}
			}
			client := &Client{
				conn:  conn,
				Stats: statsMap[config.Network],
			}

			go func() {
				defer wg.Done()
				defer func() {
					if err := conn.Close(); err != nil {
						log.Error("failed to close:", err)
					}
				}()
				client.Start(&config.Workload)
			}()
		}
	}
	wg.Wait()

	for name, stats := range statsMap {
		log.Infof("Stats of %s: %+v", name, stats)
	}
}

type Config struct {
	Concurrency int
	Network     string
	Address     string
	RecvTimeout time.Duration
	SendTimeout time.Duration
	Workload
}
