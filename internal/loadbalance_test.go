package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LBRoundRobin(t *testing.T) {
	nodes := []MachineNode{
		NewTcpNode("node1", "127.0.0.1:7890", 1.0),
		NewTcpNode("node2", "127.0.0.1:7890", 1.0),
		NewTcpNode("node3", "127.0.0.1:7890", 1.0),
	}

	roundRobinAlg := NewRoundRoBinAlg(nodes)
	node1 := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[0].Name(), node1.Name(), fmt.Sprintf("want: %s, got: %s", nodes[0].Name(), node1.Name()))

	node2 := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[1].Name(), node2.Name(), fmt.Sprintf("want: %s, got: %s", nodes[1].Name(), node2.Name()))

	node3 := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[2].Name(), node3.Name(), fmt.Sprintf("want: %s, got: %s", nodes[2].Name(), node3.Name()))

	node4 := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[0].Name(), node4.Name(), fmt.Sprintf("want: %s, got: %s", nodes[0].Name(), node4.Name()))
}

func Benchmark_LBRoundRobin(b *testing.B) {
	nodes := []MachineNode{
		NewTcpNode("node1", "127.0.0.1:7890", 1.0),
		NewTcpNode("node2", "127.0.0.1:7890", 1.0),
		NewTcpNode("node3", "127.0.0.1:7890", 1.0),
		NewTcpNode("node4", "127.0.0.1:7890", 1.0),
		NewTcpNode("node5", "127.0.0.1:7890", 1.0),
		NewTcpNode("node6", "127.0.0.1:7890", 1.0),
		NewTcpNode("node7", "127.0.0.1:7890", 1.0),
		NewTcpNode("node8", "127.0.0.1:7890", 1.0),
		NewTcpNode("node9", "127.0.0.1:7890", 1.0),
		NewTcpNode("node10", "127.0.0.1:7890", 1.0),
		NewTcpNode("node11", "127.0.0.1:7890", 1.0),
		NewTcpNode("node12", "127.0.0.1:7890", 1.0),
		NewTcpNode("node13", "127.0.0.1:7890", 1.0),
		NewTcpNode("node14", "127.0.0.1:7890", 1.0),
		NewTcpNode("node15", "127.0.0.1:7890", 1.0),
		NewTcpNode("node16", "127.0.0.1:7890", 1.0),
	}
	roundRobinAlg := NewRoundRoBinAlg(nodes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = roundRobinAlg.GetBestNode()
	}
}
