package internal

// 负载均衡器的接口
type LoadBlanceAlg interface {
	GetBestNode() (MachineNode, error)
	AddNode(node MachineNode)
	RemoveNode(node MachineNode)
}

func ChooseLoadBlance(alg LBAlg, nodes []MachineNode) LoadBlanceAlg {
	switch alg {
	case LBRoundRoBin:
		return NewRoundRoBinAlg(nodes)
	case LBIPHash:
		return NewIPHashAlg(nodes)
	case LBAtLeast:
		return NewAtLeastAlg(nodes)
	default:
		return NewRoundRoBinAlg(nodes)
	}
}
