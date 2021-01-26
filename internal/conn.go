package internal

import (
	"context"
	"io"
	"net"
)

type TrafficRW interface {
	io.Reader
	io.Writer
	io.Closer
}

// 入网连接
type InTcpConn struct {
	conn *net.TCPConn
	ctx  context.Context
}

func NewInTCPConn(in *net.TCPConn) *InTcpConn {
	return &InTcpConn{
		conn: in,
	}
}

func (i *InTcpConn) Read(p []byte) (n int, err error) {
	return i.conn.Read(p)
}

func (i *InTcpConn) Close() error {
	return i.conn.Close()
}

func (i *InTcpConn) Write(p []byte) (n int, err error) {
	return i.conn.Write(p)
}

// 出网连接
type OutTcpConn struct {
	conn *net.TCPConn
}

func NewOutTCPConn(out *net.TCPConn) *OutTcpConn {
	return &OutTcpConn{
		conn: out,
	}
}

func (o *OutTcpConn) Read(p []byte) (n int, err error) {
	return o.conn.Read(p)
}

func (o *OutTcpConn) Write(p []byte) (n int, err error) {
	return o.conn.Write(p)
}

func (o *OutTcpConn) Close() error {
	return o.conn.Close()
}
