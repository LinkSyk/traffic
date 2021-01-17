package internal

import (
	"math/rand"

	log "github.com/LinkSyk/traffic/pkg/log"
)

type BackEnd interface {
	RegisterNode(node MachineNode) error
	RegisterNodes(nodes []MachineNode) error
	GetBestNode() MachineNode
}

type TcpBackEnd struct {
	nodes []MachineNode
}

func NewTcpBackEnd(upstreams []MachineNode) BackEnd {
	tb := &TcpBackEnd{
		nodes: make([]MachineNode, 0, len(upstreams)),
	}
	for _, node := range upstreams {
		// todo: 检查每个node的健康状态, 又不健康的直接退出
		if !node.IsAlive() {
			log.Fatalf("node is deaded: %s", node.Info())
			continue
		}
		tb.nodes = append(tb.nodes, node)
	}
	return tb
}

func (t *TcpBackEnd) GetBestNode() MachineNode {
	// 伪随机算法
	// todo: 策略可随时替换, ip hash、最少连接等等
	return t.nodes[rand.Int()%len(t.nodes)]
}

func (t *TcpBackEnd) RegisterNode(node MachineNode) error {
	panic("unimpl")
}

func (t *TcpBackEnd) RegisterNodes(nodes []MachineNode) error {
	panic("unimpl")
}
