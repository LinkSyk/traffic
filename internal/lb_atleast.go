package internal

type AtLeastAlg struct {
}

func NewAtLeastAlg(nodes []MachineNode) LoadBlanceAlg {
	return &AtLeastAlg{}
}

func (lg *AtLeastAlg) GetBestNode() (MachineNode, error) {
	panic("unimpl")
}

func (lg *AtLeastAlg) AddNode(node MachineNode) {
	panic("unimpl")
}

func (lg *AtLeastAlg) RemoveNode(node MachineNode) {
	panic("unimpl")
}
