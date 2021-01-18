package internal

import (
	"sync/atomic"

	log "github.com/LinkSyk/traffic/pkg/log"
)

type LBAlg int

const (
	LBRoundRoBin LBAlg = iota
	LBAtLeast
	LBIPHash
)

type BackEnd interface {
	RegisterNode(node MachineNode) error
	RegisterNodes(nodes []MachineNode) error
	RemoveNode(node MachineNode)
	RemoveNodes(node []MachineNode)
	GetBestNode() MachineNode
}

type TcpBackEnd struct {
	nodes atomic.Value
	lb    LoadBlanceAlg
}

func NewTcpBackEnd(lbAlg LoadBlanceAlg, upstreams []MachineNode) BackEnd {
	tb := &TcpBackEnd{
		lb: lbAlg,
	}

	nodes := make([]MachineNode, 0, len(upstreams))
	for _, node := range upstreams {
		// todo: 检查每个node的健康状态, 又不健康的直接退出
		if !node.IsAlive() {
			log.Fatalf("node is deaded: %s", node.Info().String())
			continue
		}
		nodes = append(nodes, node)
	}

	tb.nodes.Store(nodes)
	return tb
}

func (t *TcpBackEnd) GetBestNode() MachineNode {
	// 伪随机算法
	// done: 策略可随时替换, ip hash、最少连接等等
	return t.lb.GetBestNode(t.nodes.Load().([]MachineNode))
}

func (t *TcpBackEnd) RegisterNode(node MachineNode) error {
	nodes := t.nodes.Load().([]MachineNode)
	t.nodes.Store(append(nodes, node))
	return nil
}

func (t *TcpBackEnd) RegisterNodes(nodes []MachineNode) error {
	ns := t.nodes.Load().([]MachineNode)
	t.nodes.Store(append(ns, nodes...))
	return nil
}

func (t *TcpBackEnd) RemoveNode(nodes MachineNode) {
	ns := t.nodes.Load().([]MachineNode)
	newNodes := make([]MachineNode, 0, len(ns))
	for _, n := range ns {
		if nodes.Name() == n.Name() {
			continue
		}
		newNodes = append(newNodes, n)
	}
	t.nodes.Store(newNodes)
}

func (t *TcpBackEnd) RemoveNodes(nodes []MachineNode) {
	ns := t.nodes.Load().([]MachineNode)
	newNodes := make([]MachineNode, 0, len(ns))
	m := make(map[string]struct{})
	for _, v := range nodes {
		m[v.Name()] = struct{}{}
	}

	for _, v := range ns {
		if _, ok := m[v.Name()]; ok {
			continue
		}
		newNodes = append(newNodes, v)
	}
	t.nodes.Store(newNodes)
}
