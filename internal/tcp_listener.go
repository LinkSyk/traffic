package internal

import (
	log "github.com/LinkSyk/traffic/pkg/log"
	"net"
	"sync"
)

type TcpListener struct {
	cfg  *TcpListenerConfig
	lb   LoadBlanceAlg
	stop chan struct{}
}

func NewTcpListener(cfg *TcpListenerConfig) *TcpListener {
	tl := &TcpListener{
		cfg:  cfg,
		stop: make(chan struct{}),
	}

	return tl
}

func (t *TcpListener) Listen() error {
	l, err := net.Listen("tcp", t.cfg.Listen)
	if err != nil {
		return err
	}

	// 用来关闭监听服务的
	go func() {
		log.Info("stop tcp server")
		l.Close()
	}()

	log.Infof("start run tcp traffic in %s", t.cfg.Listen)
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
		go t.serve(conn)
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
	node, err := t.lb.GetBestNode()
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
