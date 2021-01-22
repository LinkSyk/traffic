package internal

type IPHashAlg struct {
}

func NewIPHashAlg(nodes []Node) LoadBlanceAlg {
	return &IPHashAlg{}
}

func (lg *IPHashAlg) GetBestNode() (Node, error) {
	panic("unimpl")
}

func (lg *IPHashAlg) AddNode(node Node) {
	panic("unimpl")
}

func (lg *IPHashAlg) RemoveNode(node Node) {
	panic("unimpl")
}
