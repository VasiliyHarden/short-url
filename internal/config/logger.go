package config

import (
	"go.uber.org/zap"
	"log"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error initializing logger", err)
	}

	return logger
}
