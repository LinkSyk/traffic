package internal

type LBAlg int

const (
	LBRoundRoBin LBAlg = iota
	LBAtLeast
	LBIPHash
)

// backend 没有必要用interface实现，下层已经抽象，上层无需抽象
// backend需要处理一些业务上的东西，比如接受http请求发送过来的配置更新需求之类的
type BackEnd struct {
	lb LoadBlanceAlg
}

func NewBackEnd(lbAlg LoadBlanceAlg) BackEnd {
	tb := BackEnd{
		lb: lbAlg,
	}
	return tb
}

func (t *BackEnd) GetBestNode() MachineNode {
	// 伪随机算法
	// done: 策略可随时替换, ip hash、最少连接等等
	node, _ := t.lb.GetBestNode()
	return node
}
