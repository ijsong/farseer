package main

import (
	"log"

	"go.uber.org/zap"
	"github.com/ijsong/farseer/cmd/datagather/app"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not initialize zap logger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	command := app.NewDataGatherCommand()
	if err := command.Execute(); err != nil {
		logger.Fatal("could not execute command", zap.Error(err))
	}
}
