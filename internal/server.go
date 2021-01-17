package internal

import (
	"net"

	log "github.com/LinkSyk/traffic/pkg/log"
)

// 代理服务
type TrafficServer struct {
	tcpAddr    string
	udpAddr    string
	tcpBackEnd map[string]BackEnd
	// udpBackEnd map[string]BackEnd
}

type Option func(svr *TrafficServer)

func WithTcpAddr(addr string) Option {
	return func(svc *TrafficServer) {
		svc.tcpAddr = addr
	}
}

func WithUdpAddr(addr string) Option {
	return func(svc *TrafficServer) {
		svc.udpAddr = addr
	}
}

func NewTrafficServer(opts ...Option) TrafficServer {
	ts := TrafficServer{}
	for _, opt := range opts {
		opt(&ts)
	}

	return ts
}

func (t *TrafficServer) RunTcpListener() error {
	l, err := net.Listen("tcp", t.tcpAddr)
	if err != nil {
		return err
	}

	go func() {
		tcpListener := l.(*net.TCPListener)
		for {
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Errorf("tcpListener accept connection failed: %d", err)
			}

			tc := NewInBoundConn(conn)
			go tc.serve(t.tcpBackEnd["default"])
		}
	}()
	return nil
}
