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

type NodeInfo struct {
	Name     string
	IP       string
	Weight   float32
	ConnSums uint32
}

func (n *NodeInfo) String() string {
	return fmt.Sprintf("ip: %s, weight: %f, connection num: %d", n.IP, n.Weight, n.ConnSums)
}

// 节点的抽象
type MachineNode interface {
	Name() string
	Forward(in net.Conn) error
	IsAlive() bool
	Info() *NodeInfo
}

type TcpMachineNode struct {
	name     string
	addr     string
	weight   float32
	connSums uint32
}

func NewTcpNode(name string, addr string, weight float32) MachineNode {
	return &TcpMachineNode{
		name:   name,
		addr:   addr,
		weight: weight,
	}
}

func (n *TcpMachineNode) Name() string {
	return n.name
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

func (n *TcpMachineNode) Info() *NodeInfo {
	return &NodeInfo{
		Name:     n.name,
		IP:       n.addr,
		ConnSums: n.connSums,
	}
}
