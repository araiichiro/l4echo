package client

import (
	"context"
	"net"
	"time"

	"github.com/araiichiro/l4echo/internal/log"
	"github.com/araiichiro/l4echo/internal/payload"
)

type Client struct {
	conn net.Conn
	*Stats
}

type Workload struct {
	Loop     int
	Interval time.Duration
}

func (c *Client) Start(w *Workload) {
	errCh := make(chan error, 1)
	go func() {
		if err := c.receiveLoop(w); err != nil {
			errCh <- err
		}
		close(errCh)
	}()
	if err := c.sendLoop(w); err != nil {
		log.Error("failed to send:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil && err != context.DeadlineExceeded {
			log.Error("error in Loop:", err)
		}
	case err := <-errCh:
		if err != nil {
			log.Error("failed to receive:", err)
		}
	}

}

func (c *Client) sendLoop(w *Workload) error {
	p := payload.New()
	for i := 0; i < w.Loop; i++ {
		p.SetSeq(uint64(i))
		p.SetTime(time.Now())
		if err := payload.Send(c.conn, p); err != nil {
			return err
		}
		c.Stats.Sending()
		time.Sleep(w.Interval)
	}
	return nil

}

func (c *Client) receiveLoop(w *Workload) error {
	next := uint64(0)
	buf := make([]byte, payload.Size)
	for i := 0; i < w.Loop; i++ {
		p, err := payload.Receive(c.conn, buf)
		if err != nil {
			return err
		}
		seq := p.Seq()
		if seq == next {
			c.Stats.Update(time.Now().Sub(p.Time()))
			c.Stats.Received()
			next = seq + 1
		} else if seq > next {
			c.Stats.Drop(seq - next)
			next = seq + 1
		} else {
			c.Stats.Delayed()
		}
	}
	return nil
}
