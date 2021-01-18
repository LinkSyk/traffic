package internal

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	tcpPreDialTimeOut = 2 * time.Second
	tcpDialTimeOut    = 1 * time.Second
)

// 节点的抽象
type MachineNode interface {
	Forward(in net.Conn) error
	IsAlive() bool
	Info() string
}

type TcpMachineNode struct {
	addr   string
	weight float32
}

func NewTcpNode(addr string, weight float32) MachineNode {
	return &TcpMachineNode{
		addr:   addr,
		weight: weight,
	}
}

func (n *TcpMachineNode) IsAlive() bool {
	c, err := net.DialTimeout("tcp", n.addr, tcpPreDialTimeOut)
	if err != nil {
		return false
	}
	defer c.Close()
	return true
}

func (n *TcpMachineNode) Forward(in net.Conn) error {
	// 新建连接，后续换成从连接池拿数据
	// fixme: 修复连接泄漏的地方
	out, err := net.DialTimeout("tcp", n.addr, tcpDialTimeOut)
	if err != nil {
		return err
	}

	// in -> out
	go func() {
		for {
			_, err := io.Copy(out, in)
			if err != nil {
				return
			}
		}
	}()

	// out -> in
	go func() {
		for {
			_, err := io.Copy(in, out)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (n *TcpMachineNode) Info() string {
	return fmt.Sprintf("addr: %s", n.addr)
}
