package internal

type AtLeastAlg struct {
}

func NewAtLeastAlg(nodes []Node) LoadBlanceAlg {
	return &AtLeastAlg{}
}

func (lg *AtLeastAlg) GetBestNode() (Node, error) {
	panic("unimpl")
}

func (lg *AtLeastAlg) AddNode(node Node) {
	panic("unimpl")
}

func (lg *AtLeastAlg) RemoveNode(node Node) {
	panic("unimpl")
}
