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
	// lb           LoadBlanceAlg
	tcpListeners map[string]*TcpListener
	readyStatus  bool
}

type Option func(svr *Traffic)

func withListenTcpAddr(name string, l *TcpListener) Option {
	return func(svc *Traffic) {
		svc.tcpListeners[name] = l
	}
}

func newTrafficServer(opts ...Option) Traffic {
	ts := Traffic{}

	for _, opt := range opts {
		opt(&ts)
	}

	return ts
}

func (t *Traffic) Start() error {
	if !t.readyStatus {
		return ErrServerNotReady
	}

	g, ctx := errgroup.WithContext(context.Background())

	// 注册信号
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	g.Go(func() error {
		if err := t.runTcpListener(ctx); err != nil {
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

func (t *Traffic) runTcpListener(ctx context.Context) error {
	for name, l := range t.tcpListeners {
		if err := l.Listen(ctx); err != nil {
			log.Fatalf("run %s tcp listener failed: %v", name, err)
			return err
		}
	}
	return nil
}

func BuildTraffic(cfg *TrafficConfig) (Traffic, error) {
	t := Traffic{
		tcpListeners: map[string]*TcpListener{},
	}
	for _, tcpUp := range cfg.TcpUpstreams {
		t.tcpListeners[tcpUp.Name] = NewTcpListener(tcpUp)
	}

	t.readyStatus = true

	return t, nil
}
