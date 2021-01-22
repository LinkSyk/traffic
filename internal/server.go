package internal

import (
	"context"
	"os"
	"os/signal"

	log "github.com/LinkSyk/traffic/pkg/log"
	"golang.org/x/sync/errgroup"
)

type BType uint

const (
	TrafficTcpBackEnd BType = iota
	TrafficUdpBackEnd
)

// 代理服务
type Traffic struct {
	tcpAddr      string
	udpAddr      string
	lb           LoadBlanceAlg
	tcpListeners map[string]*TcpListener
	cfg          *TrafficConfig
}

type Option func(svr *Traffic)

func WithListenTcpAddr(addr string) Option {
	return func(svc *Traffic) {
		svc.tcpAddr = addr
	}
}

func WithListenUdpAddr(addr string) Option {
	return func(svc *Traffic) {
		svc.udpAddr = addr
	}
}

func NewTrafficServer(opts ...Option) Traffic {
	ts := Traffic{}

	for _, opt := range opts {
		opt(&ts)
	}

	return ts
}

func (t *Traffic) Start() error {
	g, ctx := errgroup.WithContext(context.Background())

	// 注册信号
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	g.Go(func() error {
		if err := t.RunTcpListener(ctx); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		select {
		case <-ch:
			return ErrStopServer
		case <-ctx.Done():
			return ErrRunning
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (t *Traffic) RunTcpListener(ctx context.Context) error {
	for name, l := range t.tcpListeners {
		if err := l.Listen(); err != nil {
			log.Fatalf("run %s tcp listener failed: %v", name, err)
			return err
		}
	}
	return nil
}
