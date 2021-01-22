package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LBRoundRobin_GetBestNode(t *testing.T) {
	nodes := []Node{
		NewSimpleNode("node1", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node2", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node3", "127.0.0.1:7890", 1.0),
	}

	roundRobinAlg := NewRoundRoBinAlg(nodes)
	node1, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[0].Name(), node1.Name(), fmt.Sprintf("want: %s, got: %s", nodes[0].Name(), node1.Name()))

	node2, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[1].Name(), node2.Name(), fmt.Sprintf("want: %s, got: %s", nodes[1].Name(), node2.Name()))

	node3, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[2].Name(), node3.Name(), fmt.Sprintf("want: %s, got: %s", nodes[2].Name(), node3.Name()))

	node4, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, nodes[0].Name(), node4.Name(), fmt.Sprintf("want: %s, got: %s", nodes[0].Name(), node4.Name()))
}

func Test_LBRoundRobin_AddNode(t *testing.T) {
	nodes := []Node{
		NewSimpleNode("node1", "127.0.0.1:7890", 1.0),
	}
	roundRobinAlg := NewRoundRoBinAlg(nodes)
	roundRobinAlg.AddNode(NewSimpleNode("node2", "127.0.0.1:7890", 1.0))
	node1, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, "node1", node1.Name(), fmt.Sprintf("want: %s, got: %s", "node1", node1.Name()))
	node2, _ := roundRobinAlg.GetBestNode()
	assert.Equal(t, "node2", node2.Name(), fmt.Sprintf("want: %s, got: %s", "node2", node2.Name()))
}

func Benchmark_LBRoundRobin(b *testing.B) {
	nodes := []Node{
		NewSimpleNode("node1", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node2", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node3", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node4", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node5", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node6", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node7", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node8", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node9", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node10", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node11", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node12", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node13", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node14", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node15", "127.0.0.1:7890", 1.0),
		NewSimpleNode("node16", "127.0.0.1:7890", 1.0),
	}
	roundRobinAlg := NewRoundRoBinAlg(nodes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roundRobinAlg.GetBestNode()
	}
}
