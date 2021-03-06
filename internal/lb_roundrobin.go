package internal

import (
	"log"
	"sync/atomic"
)

type RoundRoBinAlg struct {
	idx       int32
	upstreams atomic.Value
}

func NewRoundRoBinAlg(upstreams []Node) LoadBlanceAlg {
	rr := &RoundRoBinAlg{}
	nodes := make([]Node, 0, len(upstreams))
	for _, node := range upstreams {
		// todo: 检查每个node的健康状态, 又不健康的直接退出
		if !node.IsAlive() {
			log.Fatalf("node is deaded: %s", node.Info().String())
			continue
		}
		nodes = append(nodes, node)
	}

	rr.upstreams.Store(nodes)
	return rr
}

func (lg *RoundRoBinAlg) GetBestNode() (Node, error) {
	nodes := lg.upstreams.Load().([]Node)
	defer func() {
		// atomic.StoreInt32(&lg.idx, (lg.idx+1)%int32(len(nodes)))
		lg.idx = (lg.idx + 1) % int32(len(nodes))
	}()

	l := int32(len(nodes))
	if l == 0 {
		return nil, ErrNoAvailableNode
	}
	if l < lg.idx {
		lg.idx = l
	}

	return nodes[int(lg.idx)], nil
}

func (lg *RoundRoBinAlg) AddNode(node Node) {
	nodes := lg.upstreams.Load().([]Node)
	nodes = append(nodes, node)
	lg.upstreams.Store(nodes)
}

func (lg *RoundRoBinAlg) RemoveNode(node Node) {
	nodes := lg.upstreams.Load().([]Node)
	newNodes := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if node.Name() == n.Name() {
			continue
		}
		newNodes = append(newNodes, n)
	}
	lg.upstreams.Store(newNodes)
}
