package internal

import (
	"io"
	"net"
)

// 节点的抽象
type MachineNode interface {
	Forward(in net.Conn) error
	IsAlive() bool
}

type TcpMachineNode struct {
	addr string
}

func NewTcpNode(addr string) MachineNode {
	return &TcpMachineNode{
		addr: addr,
	}
}

func (n *TcpMachineNode) IsAlive() bool {
	return true
}

func (n *TcpMachineNode) Forward(in net.Conn) error {
	// 新建连接，后续换成从连接池拿数据
	out, err := net.Dial("tcp", n.addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err := io.Copy(out, in)
			if err != nil {
				return
			}
		}
	}()

	return nil
}
