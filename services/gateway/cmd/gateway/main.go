package main

import (
	"log"

	"github.com/jmandel1027/perspex/services/gateway/pkg/config"
	"github.com/jmandel1027/perspex/services/gateway/pkg/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Print(err)
	}

	go server.HTTP(&cfg)

	select {}
}
