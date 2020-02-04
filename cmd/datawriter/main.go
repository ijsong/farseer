package main

import (
	"log"

	"github.com/ijsong/farseer/cmd/datawriter/app"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not initialize zap logger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	command := app.NewDataWriterCommand()
	if err := command.Execute(); err != nil {
		logger.Fatal("could not execute command", zap.Error(err))
	}
}
