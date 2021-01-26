package internal

import (
	"context"
	"net"
	"sync"

	log "github.com/LinkSyk/traffic/pkg/log"
)

type TcpListener struct {
	cfg         *TcpListenerConfig
	loadbalance LoadBlanceAlg
	stop        chan struct{}
}

func NewTcpListener(cfg *TcpListenerConfig) *TcpListener {
	nodes := make([]Node, 0, len(cfg.Nodes))
	for i, n := range cfg.Nodes {
		nodes[i] = NewSimpleNode(n.Name, n.Name, float32(n.Weight))
	}

	lb := NewRoundRoBinAlg(nodes)

	tl := &TcpListener{
		cfg:         cfg,
		stop:        make(chan struct{}),
		loadbalance: lb,
	}

	return tl
}

func (t *TcpListener) Listen(ctx context.Context) error {
	l, err := net.Listen("tcp", t.cfg.Listen)
	if err != nil {
		return err
	}

	// 用来关闭监听服务的
	go func() {
		log.Info("stop tcp server")
		l.Close()
	}()

	log.Debugf("start run tcp traffic in %s", t.cfg.Listen)
	tcpListener := l.(*net.TCPListener)
	var wg sync.WaitGroup
	for {
		conn, err := tcpListener.AcceptTCP()
		log.Debug("accept tcp connection")
		if err != nil {
			log.Errorf("tcpListener accept connection failed: %v", err)
			wg.Wait()
			return err
		}

		wg.Add(1)
		go func() {
			t.serve(conn)
			wg.Done()
		}()
	}
}

func (t *TcpListener) UpdateConfig(cfg *TcpListenerConfig) error {
	t.cfg = cfg
	return nil
}

func (t *TcpListener) Stop() {

}

func (t *TcpListener) serve(conn *net.TCPConn) {
	reader := NewInTCPConn(conn)
	node, err := t.loadbalance.GetBestNode()
	if err != nil {
		log.Infof("tcpListener get best node failed: %v", err)
		return
	}

	writer, err := node.GetBestWriter()
	if err != nil {
		log.Infof("tcpListener get best writer failed: %v", err)
	}

	if err := node.Forward(reader, writer); err != nil {
		log.Infof("tcpListener node forward failed: %v", err)
	}

}
