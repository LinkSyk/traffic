package main

import (
	_ "net/http/pprof"
	"os"

	"github.com/LinkSyk/traffic/internal"
	log "github.com/LinkSyk/traffic/pkg/log"
)

func main() {
	f, err := os.Open("./traffic.yaml")
	if err != nil {
		log.Fatalf("open config file failed: %v", err)
	}

	cfg, err := internal.ParserConfig(f)
	if err != nil {
		log.Fatalf("parser config file failed: %v", err)
	}

	traffic, err := internal.BuildTraffic(cfg)
	if err != nil {
		log.Fatalf("build traffic failed: %v", err)
	}

	if err := traffic.Start(); err != nil {
		log.Fatalf("start traffic failed: %v", err)
	}
}
