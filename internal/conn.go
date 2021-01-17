package internal

import (
	"net"

	log "github.com/LinkSyk/traffic/pkg/log"
)

// 入网连接
type InBoundConn struct {
	conn net.Conn
}

func NewInBoundConn(in net.Conn) *InBoundConn {
	return &InBoundConn{
		conn: in,
	}
}

func (c *InBoundConn) serve(backEnd BackEnd) {
	node := backEnd.GetBestNode()
	if err := node.Forward(c.conn); err != nil {
		log.Errorf("forward tcp data failed: %v", err)
		return
	}
}

// 出网连接
type OutBoundConn struct {
	conn net.Conn
}

func NewOutBoundConn(out net.Conn) *OutBoundConn {
	return &OutBoundConn{
		conn: out,
	}
}
