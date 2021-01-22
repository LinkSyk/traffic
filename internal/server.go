package internal

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"

	log "github.com/LinkSyk/traffic/pkg/log"
	"golang.org/x/sync/errgroup"
)

type BType uint

const (
	TrafficTcpBackEnd BType = iota
	TrafficUdpBackEnd
)

// 代理服务
type TrafficServer struct {
	tcpAddr string
	udpAddr string
	backEnd BackEnd
}

type Option func(svr *TrafficServer)

func WithListenTcpAddr(addr string) Option {
	return func(svc *TrafficServer) {
		svc.tcpAddr = addr
	}
}

func WithListenUdpAddr(addr string) Option {
	return func(svc *TrafficServer) {
		svc.udpAddr = addr
	}
}

func WithBackEnd(backEnd BackEnd) Option {
	return func(svc *TrafficServer) {
		svc.backEnd = backEnd
	}
}

func NewTrafficServer(opts ...Option) TrafficServer {
	ts := TrafficServer{}

	for _, opt := range opts {
		opt(&ts)
	}

	return ts
}

func (t *TrafficServer) Run() error {
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

func (t *TrafficServer) RunTcpListener(ctx context.Context) error {

	l, err := net.Listen("tcp", t.tcpAddr)
	if err != nil {
		return err
	}

	// 用来关闭监听服务的
	go func() {
		<-ctx.Done()
		log.Info("stop tcp server")
		l.Close()
	}()

	log.Infof("start run tcp traffic in %s", t.tcpAddr)
	tcpListener := l.(*net.TCPListener)
	var wg sync.WaitGroup
	for {
		conn, err := tcpListener.AcceptTCP()
		log.Info("accept tcp connection")
		if err != nil {
			log.Errorf("tcpListener accept connection failed: %v", err)
			wg.Wait()
			return err
		}

		tc := NewInBoundConn(conn)
		wg.Add(1)
		go func() {
			tc.serve(ctx, t.backEnd)
			wg.Done()
		}()
	}
}
