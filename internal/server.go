package internal

import (
	"context"
	"net"
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
type TrafficServer struct {
	tcpAddr    string
	udpAddr    string
	tcpBackEnd map[string]BackEnd
	udpBackEnd map[string]BackEnd
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

func WithBackEnd(backEndType BType, name string, alg LoadBlanceAlg, nodes []MachineNode) Option {
	return func(svc *TrafficServer) {
		switch backEndType {
		case TrafficTcpBackEnd:
			svc.tcpBackEnd[name] = NewTcpBackEnd(alg, nodes)
		case TrafficUdpBackEnd:
		default:
		}
	}
}

func NewTrafficServer(opts ...Option) TrafficServer {
	ts := TrafficServer{
		tcpBackEnd: make(map[string]BackEnd),
		udpBackEnd: make(map[string]BackEnd),
	}

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
	go func() {
		tcpListener := l.(*net.TCPListener)
		for {
			conn, err := tcpListener.AcceptTCP()
			log.Info("accept tcp connection")
			if err != nil {
				log.Errorf("tcpListener accept connection failed: %d", err)
			}

			tc := NewInBoundConn(conn)
			go tc.serve(t.tcpBackEnd["default"])
		}
	}()
	return nil
}
