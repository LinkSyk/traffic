package internal

type IPHashAlg struct {
}

func NewIPHashAlg(nodes []MachineNode) LoadBlanceAlg {
	return &IPHashAlg{}
}

func (lg *IPHashAlg) GetBestNode() (MachineNode, error) {
	panic("unimpl")
}

func (lg *IPHashAlg) AddNode(node MachineNode) {
	panic("unimpl")
}

func (lg *IPHashAlg) RemoveNode(node MachineNode) {
	panic("unimpl")
}
