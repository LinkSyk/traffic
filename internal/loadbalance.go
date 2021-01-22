package internal

// 负载均衡器的接口
type LoadBlanceAlg interface {
	GetBestNode() (Node, error)
	AddNode(node Node)
	RemoveNode(node Node)
}

func ChooseLoadBlance(alg LBAlg, nodes []Node) LoadBlanceAlg {
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
