package main

import (
	"github.com/LinkSyk/traffic/internal"
	"github.com/LinkSyk/traffic/pkg/log"
)

func main() {
	svr := internal.NewTrafficServer(
		internal.WithListenTcpAddr("0.0.0.0:9999"),
		internal.WithBackEnd(internal.TrafficTcpBackEnd, "default", internal.ChooseLoadBlance(internal.LBRoundRoBin), []internal.MachineNode{
			internal.NewTcpNode("127.0.0.1:7890", 1),
			// internal.NewTcpNode("127.0.0.1:8089"),
		}),
	)

	if err := svr.Run(); err != nil {
		log.Fatalf("run traffic failed: %v", err)
	}
}
