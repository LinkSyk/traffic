package main

import (
	"github.com/LinkSyk/traffic/internal"
	"github.com/LinkSyk/traffic/pkg/log"
)

func main() {
	nodes := []internal.MachineNode{
		internal.NewTcpNode("node1", "127.0.0.1:7890", 1),
		internal.NewTcpNode("node2", "127.0.0.1:7891", 1),
	}
	lg := internal.NewRoundRoBinAlg(nodes)
	backEnd := internal.NewBackEnd(lg)

	svr := internal.NewTrafficServer(
		internal.WithListenTcpAddr("0.0.0.0:9999"),
		internal.WithBackEnd(backEnd),
	)

	if err := svr.Run(); err != nil {
		log.Fatalf("run traffic failed: %v", err)
	}
}
