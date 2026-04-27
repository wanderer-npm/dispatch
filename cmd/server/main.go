package main

import (
	"log"
	"os"

	"github.com/wanderer-npm/dispatch/internal/config"
	"github.com/wanderer-npm/dispatch/internal/server"
)

func main() {
	cfgPath := "config.yml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("[dispatch] failed to load config: %v", err)
	}

	if err := server.New(cfg).Start(); err != nil {
		log.Fatalf("[dispatch] server error: %v", err)
	}
}
