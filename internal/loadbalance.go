package internal

type LoadBlanceAlg interface {
	GetBestNode(nodes []MachineNode) MachineNode
}

type RoundRoBinAlg struct {
	idx int32
}

func NewRoundRoBinAlg() LoadBlanceAlg {
	return &RoundRoBinAlg{}
}

func (lg *RoundRoBinAlg) GetBestNode(nodes []MachineNode) MachineNode {
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

type IPHashAlg struct {
}

func NewIPHashAlg() LoadBlanceAlg {
	return &IPHashAlg{}
}

func (lg *IPHashAlg) GetBestNode(nodes []MachineNode) MachineNode {
	panic("unimpl")
}

type AtLeastAlg struct {
}

func NewAtLeastAlg() LoadBlanceAlg {
	return &AtLeastAlg{}
}

func (lg *AtLeastAlg) GetBestNode(nodes []MachineNode) MachineNode {
	panic("unimpl")
}

func ChooseLoadBlance(alg LBAlg) LoadBlanceAlg {
	switch alg {
	case LBRoundRoBin:
		return NewRoundRoBinAlg()
	case LBIPHash:
		return NewIPHashAlg()
	case LBAtLeast:
		return NewAtLeastAlg()
	default:
		return NewRoundRoBinAlg()
	}
}
