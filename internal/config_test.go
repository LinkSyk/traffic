package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParserConfig(t *testing.T) {
	f, err := os.Open("../test/traffic.yaml")
	if err != nil {
		t.Fatalf("open file error")
	}

	config, err := ParserConfig(f)
	if err != nil {
		t.Fatalf("parser config file failed")
	}

	tcpStream := config.TcpUpstreams[0]
	assert.Equal(t, "round-robin", tcpStream.LbAlg, "want: %s, got: %s", "round-robin", tcpStream.LbAlg)
	node := tcpStream.Nodes[0]
	assert.Equal(t, "127.0.0.1:7890", node.Addr, "want: %s, got: %s", "127.0.0.1:7890", node.Addr)
}
