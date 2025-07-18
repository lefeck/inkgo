package main

import (
	"flag"
	"go.uber.org/zap"
	"inkgo/config"
	"inkgo/server"
	"log"
)

var (
	appConfig = flag.String("config", "config/app.yaml", "application config path")
)

func main() {
	flag.Parse()
	conf, err := config.LoadConfig(*appConfig)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any zap.ReplaceGlobals(logger)

	s, err := server.New(conf, logger)
	if err != nil {
		log.Fatalf("Init server failed: %v", err)
	}
	if err := s.Run(); err != nil {
		log.Fatalf("Run server failed: %v", err)
	}
}
