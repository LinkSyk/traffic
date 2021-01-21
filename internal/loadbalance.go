package internal

import (
	"log"
	"sync/atomic"
)

type LoadBlanceAlg interface {
	GetBestNode() MachineNode
	AddNode(node MachineNode)
	RemoveNode(node MachineNode)
}

type RoundRoBinAlg struct {
	idx       int32
	upstreams atomic.Value
}

func NewRoundRoBinAlg(upstreams []MachineNode) LoadBlanceAlg {
	rr := &RoundRoBinAlg{}
	nodes := make([]MachineNode, 0, len(upstreams))
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

func (lg *RoundRoBinAlg) GetBestNode() MachineNode {
	nodes := lg.upstreams.Load().([]MachineNode)
	defer func() {
		// atomic.StoreInt32(&lg.idx, (lg.idx+1)%int32(len(nodes)))
		lg.idx = (lg.idx + 1) % int32(len(nodes))
	}()

	l := int32(len(nodes))
	if l < lg.idx {
		lg.idx = l
	}

	return nodes[int(lg.idx)]
}

func (lg *RoundRoBinAlg) AddNode(node MachineNode) {
	panic("unimpl")
}

func (lg *RoundRoBinAlg) RemoveNode(node MachineNode) {
	panic("unimpl")
}

type IPHashAlg struct {
}

func NewIPHashAlg() LoadBlanceAlg {
	return &IPHashAlg{}
}

func (lg *IPHashAlg) GetBestNode() MachineNode {
	panic("unimpl")
}

func (lg *IPHashAlg) AddNode(node MachineNode) {
	panic("unimpl")
}

func (lg *IPHashAlg) RemoveNode(node MachineNode) {
	panic("unimpl")
}

type AtLeastAlg struct {
}

func NewAtLeastAlg() LoadBlanceAlg {
	return &AtLeastAlg{}
}

func (lg *AtLeastAlg) GetBestNode() MachineNode {
	panic("unimpl")
}

func (lg *AtLeastAlg) AddNode(node MachineNode) {
	panic("unimpl")
}

func (lg *AtLeastAlg) RemoveNode(node MachineNode) {
	panic("unimpl")
}

func ChooseLoadBlance(alg LBAlg, nodes []MachineNode) LoadBlanceAlg {
	switch alg {
	case LBRoundRoBin:
		return NewRoundRoBinAlg(nodes)
	case LBIPHash:
		return NewIPHashAlg()
	case LBAtLeast:
		return NewAtLeastAlg()
	default:
		return NewRoundRoBinAlg(nodes)
	}
}
