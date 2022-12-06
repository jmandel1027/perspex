package main

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/jmandel1027/perspex/services/gateway/pkg/config"
	"github.com/jmandel1027/perspex/services/gateway/pkg/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		otelzap.L().Ctx(context.TODO()).Error(err.Error())
	}

	go server.HTTP(&cfg)

	select {}
}
