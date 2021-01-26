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
	tcpPreDialTimeOut = 3 * time.Second
	tcpDialTimeOut    = 3 * time.Second
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
	closeConn := func() {
		reader.Close()
		writer.Close()
	}

	// tcp的发送窗口大小一般在2^16，大约是64K buffer，tcp首部中的option可以扩充这个值，可以扩大到2^31
	ioBuffer := make([]byte, 0xffff)
	oiBuffer := make([]byte, 0xffff)
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
					log.Errorf("%s: in -> out read data failed: %v", info, err)
					return ErrSocketRead

				}
				log.Debugf("%s: in -> out read %d bytes data", info, cnt)

				cnt, err = writer.Write(ioBuffer[:cnt])
				if err != nil {
					log.Errorf("%s: in -> out write data failed: %v", info, err)
					return ErrSocketWrite
				}
				log.Infof("%s: in -> out write %d bytes data", info, cnt)
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
					log.Errorf("%s: out -> in read data failed: %v", info, err)
					return ErrSocketRead
				}
				log.Infof("%s: in -> out read %d bytes data", info, cnt)

				cnt, err = writer.Write(oiBuffer[:cnt])
				if err != nil {
					log.Errorf("%s: out -> in write data failed: %v", info, err)
					return ErrSocketWrite
				}
				log.Infof("%s: in -> out write %d bytes data", info, cnt)
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

func (n *SimpleNode) GetBestWriter() (TrafficWriteCloser, error) {
	return net.DialTimeout("tcp", n.addr, tcpDialTimeOut)
}
