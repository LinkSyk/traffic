package internal

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/LinkSyk/traffic/pkg/log"
	"golang.org/x/sync/errgroup"
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
	Forward(ctx context.Context, in *InBoundConn) error
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

func (n *TcpMachineNode) Forward(ctx context.Context, in *InBoundConn) error {
	// 新建连接，后续换成从连接池拿数据
	// fixme: 修复连接泄漏的地方
	out, err := net.DialTimeout("tcp", n.addr, tcpDialTimeOut)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		in.conn.Close()
		out.Close()
		return nil
	})

	// in -> out
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				log.Info(fmt.Sprintf("%s: in -> out stop serve", n.Info().String()))
				return ErrStopServer
			default:
				_, err := io.Copy(out, in.conn)
				if err != nil {
					log.Error(fmt.Sprintf("%s: in -> out forward data failed: %v", n.Info().String(), err))
					return err
				}

			}
		}
	})

	// out -> in
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				log.Info(fmt.Sprintf("%s: out -> in stop serve", n.Info().String()))
				return ErrStopServer
			default:
				_, err := io.Copy(in.conn, out)
				if err != nil {
					log.Error(fmt.Sprintf("%s: out -> in forward data failed: %v", n.Info().String(), err))
					return err
				}
			}
		}
	})

	return g.Wait()
}

func (n *TcpMachineNode) Info() *NodeInfo {
	return &NodeInfo{
		Name:     n.name,
		IP:       n.addr,
		ConnSums: n.connSums,
	}
}
