package internal

import (
	"context"
	"fmt"
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
type Node interface {
	Name() string
	Forward(reader TrafficReadCloser, writer TrafficWriteCloser) error
	GetBestWriter() (TrafficWriteCloser, error)
	IsAlive() bool
	Info() *NodeInfo
}

type SimpleNode struct {
	name     string
	addr     string
	weight   float32
	connSums uint32
}

func NewSimpleNode(name string, addr string, weight float32) Node {
	return &SimpleNode{
		name:   name,
		addr:   addr,
		weight: weight,
	}
}

func (n *SimpleNode) Name() string {
	return n.name
}

func (n *SimpleNode) IsAlive() bool {
	c, err := net.DialTimeout("tcp", n.addr, tcpPreDialTimeOut)
	if err != nil {
		return false
	}
	defer c.Close()
	return true
}

func (n *SimpleNode) Forward(reader TrafficReadCloser, writer TrafficWriteCloser) error {
	// 新建连接，后续换成从连接池拿数据
	// fixme: 修复连接泄漏的地方
	out, err := net.DialTimeout("tcp", n.addr, tcpDialTimeOut)
	if err != nil {
		return err
	}

	closeConn := func() {
		reader.Close()
		writer.Close()
	}

	ioBuffer := make([]byte, 10240)
	oiBuffer := make([]byte, 10240)
	info := n.Info().String()
	g, ctx := errgroup.WithContext(context.Background())

	// close connection
	g.Go(func() error {
		<-ctx.Done()
		closeConn()
		return nil
	})

	// in -> out
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ErrStopServer
			default:
				cnt, err := reader.Read(ioBuffer)
				if err != nil {
					log.Error(fmt.Sprintf("%s: in -> out read data failed: %v", info, err))
					return ErrSocketRead

				}
				log.Info(fmt.Sprintf("%s: in -> out read %d bytes data", info, cnt))

				cnt, err = out.Write(ioBuffer[:cnt])
				if err != nil {
					log.Error(fmt.Sprintf("%s: in -> out write data failed: %v", info, err))
					return ErrSocketWrite
				}
				log.Info(fmt.Sprintf("%s: in -> out write %d bytes data", info, cnt))
			}
		}
	})

	// out -> in
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ErrStopServer
			default:
				cnt, err := reader.Read(oiBuffer)
				if err != nil {
					log.Error(fmt.Sprintf("%s: out -> in read data failed: %v", n.Info().String(), err))
					return ErrSocketRead
				}
				log.Info(fmt.Sprintf("%s: in -> out read %d bytes data", info, cnt))

				cnt, err = writer.Write(oiBuffer[:cnt])
				if err != nil {
					log.Error(fmt.Sprintf("%s: out -> in write data failed: %v", n.Info().String(), err))
					return ErrSocketWrite
				}
				log.Info(fmt.Sprintf("%s: in -> out write %d bytes data", info, cnt))
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (n *SimpleNode) Info() *NodeInfo {
	return &NodeInfo{
		Name:     n.name,
		IP:       n.addr,
		ConnSums: n.connSums,
	}
}

func (n *SimpleNode) GetOutConn() *OutTcpConn {
	return nil
}

func (n *SimpleNode) GetBestWriter() (TrafficWriteCloser, error) {
	return nil, nil
}
