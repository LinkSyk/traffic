package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/LinkSyk/traffic/internal"
	"github.com/LinkSyk/traffic/pkg/log"
)

func main() {
	nodes := []internal.Node{
		internal.NewSimpleNode("node1", "127.0.0.1:7890", 1),
		// internal.NewSimpleNode("node2", "127.0.0.1:7891", 1),
	}
	lg := internal.NewRoundRoBinAlg(nodes)
	backEnd := internal.NewBackEnd(lg)

	svr := internal.NewTrafficServer(
		internal.WithListenTcpAddr("0.0.0.0:9999"),
		internal.WithBackEnd(backEnd),
	)

	go http.ListenAndServe("0.0.0.0:9991", nil)

	if err := svr.Start(); err != nil {
		log.Fatalf("run traffic failed: %v", err)
	}
}
